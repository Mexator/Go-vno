package nameserver

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"path"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	fsapi "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	"github.com/Mexator/Go-vno/pkg/cache"
)

type (
	GRPCServer struct {
		servers      sync.Map // map[string]chan struct{}
		replicatecnt uint
		root         root
	}
)

var (
	grpcopts = []grpc.DialOption{grpc.WithInsecure()}

	ErrDirIsFile      = fmt.Errorf("Directory is actually a file")
	ErrDirNotExists   = fmt.Errorf("Directory does not exists")
	ErrEntryExists    = fmt.Errorf("Entry already exists")
	ErrEntryNotExists = fmt.Errorf("Entry does not exists")
	ErrFileExists     = fmt.Errorf("File already exists")
	ErrFileIsDir      = fmt.Errorf("File is actually a directory")
	ErrFileNotExists  = fmt.Errorf("File does not exists")
	ErrNoFileSevers   = fmt.Errorf("Needed number of fileservers is not available")
	ErrNoPeerInfo     = fmt.Errorf("Cannot obtain name server address from context")
)

// Generates unique IDs for file
func getInode() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return id.String()
}

// NewServer creates new NameServer
func NewServer() nsapi.NameServerServer {
	return &GRPCServer{replicatecnt: 2}
}

func (g *GRPCServer) nservers() int {
	len := 0
	g.servers.Range(func(key, value interface{}) bool {
		len++
		return true
	})
	return len
}

func (g *GRPCServer) pickServers(n int) ([]string, error) {
	if g.nservers() < n {
		return nil, ErrNoFileSevers
	}

	seed := sync.Once{}
	seed.Do(func() {
		rand.Seed(time.Now().Unix())
	})

	keys := make([]string, 0, g.nservers())
	g.servers.Range(func(k, _ interface{}) bool {
		keys = append(keys, k.(string))
		return true
	})

	rand.Shuffle(n, func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	return keys[:n], nil
}

func (g *GRPCServer) ReadDirAll(
	ctx context.Context,
	lreq *nsapi.ReadDirAllRequest,
) (*nsapi.ReadDirAllResponse, error) {
	d, err := g.root.lookup_dir(lreq.Path)
	if err != nil {
		return nil, errors.Wrap(err, "Lookup")
	}
	nodes := []*nsapi.Node{}
	d.entries.Range(func(k, value interface{}) bool {
		node := value.(dirEntry).AsNode()
		node.Name = k.(string)
		nodes = append(nodes, node)
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

	dir, err := g.root.lookup_dir(d)
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}

	_, ok := dir.entries.Load(node)

	if ok {
		return nil, errors.Wrap(ErrFileExists, "Create")
	}

	if creq.IsDir {
		dir.entries.Store(node, &directory{})
		return &nsapi.CreateResponse{}, nil
	}

	fsarr, err := g.pickServers(int(g.replicatecnt))
	if err != nil {
		return nil, err
	}
	inode := getInode()

	// Falsestarting servers in order to get some errors beforehand
	for _, fs := range fsarr {
		_, err := cache.GrpcDial(fs, grpcopts...)
		if err != nil {
			return nil, errors.Wrap(err, "Create: Failed to dial with fileserver")
		}
	}

	fl := &file{inode: inode}

	for _, url := range fsarr {
		conn, _ := cache.GrpcDial(url, grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		_, err := f.Create(ctx, &fsapi.CreateRequest{Inode: inode})
		if err != nil {
			log.Print(err)
			goto cleanup
		}
		fl.fileservers.Store(url, struct{}{})
	}

	dir.entries.Store(node, fl)
	return &nsapi.CreateResponse{}, nil

cleanup:
	fl.fileservers.Range(func(fsurl, _ interface{}) bool {
		conn, _ := cache.GrpcDial(fsurl.(string), grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		f.Remove(ctx, &fsapi.RemoveRequest{Inode: inode})
		return true
	})
	return nil, err
}

func (g *GRPCServer) Rename(
	ctx context.Context,
	rreq *nsapi.RenameRequest,
) (*nsapi.RenameResponse, error) {
	log.Println("Rename", rreq)
	from, err := g.root.lookup_dir(rreq.FromDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Rename: Failed to find directory `%s'", rreq.FromDir)
	}

	to, err := g.root.lookup_dir(rreq.ToDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Rename: Failed to find directory `%s'")
	}

	e, ok := to.entries.Load(rreq.ToName)
	if ok {
		e.(dirEntry).Remove(ctx)
	}

	ent_int, ok := from.entries.LoadAndDelete(rreq.FromName)
	if !ok {
		fmt := "Rename: Failed to find entry `%s' in `%s'"
		return nil, errors.Wrapf(ErrEntryNotExists, fmt, rreq.FromName, rreq.FromDir)
	}

	to.entries.Store(rreq.ToName, ent_int.(dirEntry))

	return &nsapi.RenameResponse{}, nil
}

func (g *GRPCServer) Remove(
	ctx context.Context,
	rreq *nsapi.RemoveRequest,
) (*nsapi.RemoveResponse, error) {
	d, fname := path.Split(rreq.Path)
	dir, err := g.root.lookup_dir(d)
	if err != nil {
		return nil, errors.Wrap(err, "Remove")
	}
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

// MapFS return url of fileserver that stores needed file together with
// inode of the file on that server
func (g *GRPCServer) MapFS(
	ctx context.Context,
	mreq *nsapi.MapFSRequest,
) (*nsapi.MapFSResponse, error) {
	d, fname := path.Split(mreq.Path)
	dir, err := g.root.lookup_dir(d)
	if err != nil {
		return nil, errors.Wrapf(err, "MapFS: dir: `%s' name: `%s'", d, fname)
	}

	e, ok := dir.entries.Load(fname)
	if !ok {
		return nil, errors.Wrapf(ErrFileNotExists, "MapFS: dir: `%s' name: `%s'", d, fname)
	}

	file, ok := e.(*file)
	if !ok {
		return nil, errors.Wrapf(ErrFileIsDir, "MapFS: dir: `%s' name: `%s'", d, fname)
	}

	var Fsurls []string
	file.fileservers.Range(func(k, _ interface{}) bool {
		Fsurls = append(Fsurls, k.(string))
		return true
	})

	return &nsapi.MapFSResponse{Fsurls: Fsurls, Inode: file.inode}, nil
}

func (g *GRPCServer) healthcheck(addr string, ch chan struct{}) {
	for {
		select {
		case <-ch:
			continue
		case <-time.After(10 * time.Second):
			log.Println("Disconnected", addr)
			g.servers.Delete(addr)
			g.root.dir.replicate(addr, g)
			return
		}
	}
}

// ConnectFileServer is a request that fileservers use to connect to
// nameserver. After connection name server uses them to store files
func (g *GRPCServer) ConnectFileServer(
	ctx context.Context,
	req *nsapi.ConnectRequest,
) (*nsapi.ConnectResponse, error) {
	// Obtain sender address from context
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrNoPeerInfo
	}

	host, _, err := net.SplitHostPort(peer.Addr.String())
	log.Printf("Connect request from %s", host)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to split host and port")
	}

	addr := fmt.Sprintf("%s:%d", host, req.Port)
	c, ok := g.servers.Load(addr)
	if ok {
		log.Println("Disconnect timed out", addr)
		c.(chan struct{}) <- struct{}{}
		return &nsapi.ConnectResponse{}, nil
	}

	// Check if sender is a working file server
	conn, err := cache.GrpcDial(addr, grpcopts...)
	if err != nil {
		return nil, err
	}

	fileserver := fsapi.NewFileServerClient(conn)
	_, err = fileserver.ReportDiskSpace(ctx, &fsapi.Empty{})
	if err != nil {
		return nil, err
	}

	ch := make(chan struct{})
	g.servers.Store(addr, ch)
	go g.healthcheck(addr, ch)
	return &nsapi.ConnectResponse{}, nil
}
