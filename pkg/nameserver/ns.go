package nameserver

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
		servers      map[string]struct{}
		replicatecnt uint
		root         root
	}

	NameServerError struct {
		fmt  string
		name string
	}
)

var (
	grpcopts = []grpc.DialOption{grpc.WithInsecure()}

	ErrNoFileSevers = fmt.Errorf("No fileservers available")

	ErrFileExists    = fmt.Errorf("File already exists")
	ErrFileNotExists = fmt.Errorf("File does not exists")
	ErrDirNotExists  = fmt.Errorf("Directory does not exists")
	ErrDirIsFile     = fmt.Errorf("Directory is actually a file")
	ErrFileIsDir     = fmt.Errorf("File is actually a directory")

	ErrNoPeerInfo = fmt.Errorf("Cannot obtain name server address from context")
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
	serversset := make(map[string]struct{})
	return &GRPCServer{servers: serversset, replicatecnt: 2}
}

func (g *GRPCServer) pickServers(n int) ([]string, error) {
	// Avoid IndexOutOfBounds
	if len(g.servers) < n {
		n = len(g.servers)
	}

	if n <= 0 {
		return nil, ErrNoFileSevers
	}

	seed := sync.Once{}
	seed.Do(func() {
		rand.Seed(time.Now().Unix())
	})

	keys := make([]string, 0, len(g.servers))
	for k := range g.servers {
		keys = append(keys, k)
	}

	rand.Shuffle(n, func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	return keys[:n], nil
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
		return nil, errors.Wrap(ErrDirIsFile, "Lookup")
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
		return nil, errors.Wrap(ErrDirIsFile, "Create")
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

	fileservers, err := g.pickServers(int(g.replicatecnt))
	if err != nil {
		return nil, err
	}
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
			log.Print(err)
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

// MapFS return url of fileserver that stores needed file together with
// inode of the file on that server
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
		return nil, errors.Wrap(ErrDirIsFile, "MapFS")
	}

	e, ok := dir.entries.Load(fname)
	if !ok {
		return nil, errors.Wrap(ErrFileNotExists, "MapFS")
	}

	file, ok := e.(*file)
	if !ok {
		return nil, errors.Wrap(ErrFileIsDir, "MapFS")
	}

	return &nsapi.MapFSResponse{Fsurl: file.fileservers[0], Inode: file.inode}, nil
}

// ConnectFileServer is a request that fileservers use to connect to
// nameserver. After connection name server uses them to store files
func (g *GRPCServer) ConnectFileServer(
	ctx context.Context,
	_ *nsapi.ConnectRequest,
) (*nsapi.ConnectResponse, error) {
	// Obtain sender address from context
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, ErrNoPeerInfo
	}

	addr := peer.Addr.String()

	_, ok = g.servers[addr]
	if ok {
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

	g.servers[addr] = struct{}{}
	return &nsapi.ConnectResponse{}, nil
}
