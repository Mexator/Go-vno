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
		servers   []string
		replicate uint
		root      root
	}
	root struct{ dir directory }

	ent interface {
		AsNode() *nsapi.Node
	}

	directory struct {
		name    string
		entries sync.Map // map[string]ent
	}
	file struct {
		name string
		size uint64

		inode       string
		nameservers []string
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
	return &nsapi.Node{IsDir: false, Name: d.name, Size: 4096}
}

func (f *file) AsNode() *nsapi.Node {
	return &nsapi.Node{IsDir: true, Name: f.name, Size: f.size}
}

func NewServer(servers []string) nsapi.NameServerServer {
	return &GRPCServer{servers: servers, replicate: 2}
}

func (r *root) lookup(path string) (*directory, error) {
	sep := string(os.PathSeparator)
	parts := strings.Split(path, sep)

	var d *directory = &r.dir
	cur := sep

	for _, p := range parts {
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

func (g *GRPCServer) Lookup(
	ctx context.Context,
	lreq *nsapi.LookupRequest,
) (*nsapi.LookupResponse, error) {
	d, err := g.root.lookup(lreq.Path)
	if err != nil {
		return nil, errors.Wrap(err, "Lookup failed")
	}
	nodes := []*nsapi.Node{}
	d.entries.Range(func(_, value interface{}) bool {
		nodes = append(nodes, value.(ent).AsNode())
		return true
	})
	return &nsapi.LookupResponse{Nodes: nodes}, nil
}

func (g *GRPCServer) Create(
	ctx context.Context,
	creq *nsapi.CreateRequest,
) (*nsapi.CreateResponse, error) {
	d, node := path.Split(creq.Path)
	dir, err := g.root.lookup(d)
	if err != nil {
		return nil, errors.Wrap(err, "Create: Lookup failed")
	}

	_, ok := dir.entries.Load(node)
	if ok {
		return nil, errors.Wrap(err, "Create: entry already exists")
	}

	if creq.IsDir {
		var e ent = &directory{name: node}
		dir.entries.Store(node, e)
		return &nsapi.CreateResponse{}, nil
	}

	fileservers := g.pickServers(int(g.replicate))
	inode := getInode()

	// Falsestarting servers in order to get some errors beforehand
	for _, fs := range fileservers {
		_, err := cache.GrpcDial(fs, grpcopts...)
		if err != nil {
			return nil, errors.Wrap(err, "Create: Failed to dial with fileserver")
		}
	}

	for i := 0; i < len(fileservers); i++ {
		conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		_, err := f.Create(ctx, &fsapi.CreateRequest{Inode: inode})
		if err == nil {
			continue
		}
		for ; i >= 0; i-- {
			conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
			f := fsapi.NewFileServerClient(conn)
			f.Remove(ctx, &fsapi.RemoveRequest{Inode: inode})
		}
		return nil, err
	}

	var e ent = &file{name: node, inode: inode}
	dir.entries.Store(node, e)
	return &nsapi.CreateResponse{}, nil
}

func (g *GRPCServer) Remove(
	ctx context.Context,
	rreq *nsapi.RemoveRequest,
) (*nsapi.RemoveResponse, error) {
	panic("Todo")
	return nil, nil
}

func (g *GRPCServer) MapFS(
	ctx context.Context,
	mreq *nsapi.MapFSRequest,
) (*nsapi.MapFSResponse, error) {
	d, fname := path.Split(mreq.Path)
	dir, err := g.root.lookup(d)
	if err != nil {
		return nil, errors.Wrap(err, "MapFS: Lookup failed")
	}

	ent, ok := dir.entries.Load(fname)
	if !ok {
		return nil, ErrFileNotExists
	}

	file, ok := ent.(*file)
	if !ok {
		return nil, ErrFileDir
	}

	return &nsapi.MapFSResponse{Fsurl: file.nameservers[0], Inode: file.inode}, nil
}
