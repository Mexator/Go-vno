package nameserver

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	fsapi "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	"github.com/Mexator/Go-vno/pkg/cache"
)

type (
	GRPCServer struct {
		servers      []string
		replicatecnt uint
		root         root
	}
	root struct{ dir directory }

	dirEntry interface {
		AsNode() *nsapi.Node
		Remove(context.Context) error
	}

	directory struct {
		name    string
		entries sync.Map // map[string]dirEntry
	}
	file struct {
		name string

		inode       string
		fileservers []string
	}

	NameServerError struct {
		fmt  string
		name string
	}
)

var (
	grpcopts = []grpc.DialOption{grpc.WithInsecure()}

	ErrFileExists    = fmt.Errorf("File already exists")
	ErrFileNotExists = fmt.Errorf("File does not exists")
	ErrDirNotExists  = fmt.Errorf("Directory does not exists")
	ErrDirFile       = fmt.Errorf("Directory is actually a file")
	ErrFileDir       = fmt.Errorf("File is actually a directory")
)

func getInode() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return id.String()
}

func (d *directory) AsNode() *nsapi.Node {
	return &nsapi.Node{IsDir: true, Name: d.name, Size: 4096}
}

func (d *directory) Remove(ctx context.Context) error {
	var err error
	d.entries.Range(func(k, nodeint interface{}) bool {
		if e := nodeint.(dirEntry).Remove(ctx); e != nil {
			err = e
		}
		d.entries.Delete(k)
		return true
	})
	return err
}

func (f *file) AsNode() *nsapi.Node {
	return &nsapi.Node{IsDir: false, Name: f.name}
}

func (f *file) Remove(ctx context.Context) error {
	var err error
	for _, fsurl := range f.fileservers {
		conn, e := cache.GrpcDial(fsurl, grpcopts...)
		if err != nil {
			err = e
		}
		fs := fsapi.NewFileServerClient(conn)

		_, e = fs.Remove(ctx, &fsapi.RemoveRequest{Inode: f.inode})
		if err != nil {
			err = e
		}
	}
	return err
}

func NewServer(servers []string) nsapi.NameServerServer {
	return &GRPCServer{servers: servers, replicatecnt: 2}
}

func (r *root) lookup(path string) (dirEntry, error) {
	sep := string(os.PathSeparator)
	parts := strings.Split(path, sep)[1:] // starts with /

	var d *directory = &r.dir
	cur := sep

	for _, p := range parts {
		if p == "" {
			return d, nil
		}

		ent, ok := d.entries.Load(p)
		cur = cur + sep + p
		if !ok {
			return nil, ErrDirNotExists
		}

		d, ok = ent.(*directory)
		if !ok {
			return nil, ErrDirFile
		}
	}

	return d, nil
}

func (g *GRPCServer) pickServers(n int) []string {
	seed := sync.Once{}
	seed.Do(func() {
		rand.Seed(time.Now().Unix())
	})

	a := g.servers
	rand.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a[:n]
}

func (g *GRPCServer) ReadDirAll(
	ctx context.Context,
	lreq *nsapi.ReadDirAllRequest,
) (*nsapi.ReadDirAllResponse, error) {
	ent, err := g.root.lookup(lreq.Path)
	if err != nil {
		return nil, errors.Wrap(err, "Lookup")
	}
	d, ok := ent.(*directory)
	if !ok {
		return nil, errors.Wrap(ErrDirFile, "Lookup")
	}
	nodes := []*nsapi.Node{}
	d.entries.Range(func(_, value interface{}) bool {
		nodes = append(nodes, value.(dirEntry).AsNode())
		return true
	})
	return &nsapi.ReadDirAllResponse{Nodes: nodes}, nil
}

func (g *GRPCServer) Create(
	ctx context.Context,
	creq *nsapi.CreateRequest,
) (*nsapi.CreateResponse, error) {
	var d, node string

	d, node = path.Split(creq.Path)

	ent, err := g.root.lookup(d)
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}

	dir, ok := ent.(*directory)
	if !ok {
		return nil, errors.Wrap(ErrDirFile, "Create")
	}

	_, ok = dir.entries.Load(node)
	if ok {
		return nil, errors.Wrap(ErrFileExists, "Create")
	}

	if creq.IsDir {
		var e dirEntry = &directory{name: node, entries: sync.Map{}}
		dir.entries.Store(node, e)
		return &nsapi.CreateResponse{}, nil
	}

	fileservers := g.pickServers(int(g.replicatecnt))
	inode := getInode()

	// Falsestarting servers in order to get some errors beforehand
	for _, fs := range fileservers {
		_, err := cache.GrpcDial(fs, grpcopts...)
		if err != nil {
			return nil, errors.Wrap(err, "Create: Failed to dial with fileserver")
		}
	}

	var i int
	var e dirEntry

	for i = 0; i < len(fileservers); i++ {
		conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		_, err := f.Create(ctx, &fsapi.CreateRequest{Inode: inode})
		if err != nil {
			goto cleanup
		}
	}

	e = &file{name: node, inode: inode, fileservers: fileservers}
	dir.entries.Store(node, e)
	return &nsapi.CreateResponse{}, nil

cleanup:
	for ; i >= 0; i-- {
		conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		f.Remove(ctx, &fsapi.RemoveRequest{Inode: inode})
	}
	return nil, err
}

func (g *GRPCServer) Remove(
	ctx context.Context,
	rreq *nsapi.RemoveRequest,
) (*nsapi.RemoveResponse, error) {
	d, fname := path.Split(rreq.Path)
	dirint, err := g.root.lookup(d)
	if err != nil {
		return nil, errors.Wrap(err, "Remove")
	}
	dir := dirint.(*directory)
	node, ok := dir.entries.Load(fname)
	if !ok {
		return nil, ErrFileNotExists
	}

	err = node.(dirEntry).Remove(ctx)
	if err != nil {
		return nil, err
	}
	dir.entries.Delete(fname)
	return &nsapi.RemoveResponse{}, nil
}

// MapFS maps
func (g *GRPCServer) MapFS(
	ctx context.Context,
	mreq *nsapi.MapFSRequest,
) (*nsapi.MapFSResponse, error) {
	d, fname := path.Split(mreq.Path)
	ent, err := g.root.lookup(d)
	if err != nil {
		return nil, errors.Wrap(err, "MapFS")
	}
	dir, ok := ent.(*directory)
	if !ok {
		return nil, errors.Wrap(ErrDirFile, "MapFS")
	}

	e, ok := dir.entries.Load(fname)
	if !ok {
		return nil, errors.Wrap(ErrFileNotExists, "MapFS")
	}

	file, ok := e.(*file)
	if !ok {
		return nil, errors.Wrap(ErrFileDir, "MapFS")
	}

	return &nsapi.MapFSResponse{Fsurl: file.fileservers[0], Inode: file.inode}, nil
}
