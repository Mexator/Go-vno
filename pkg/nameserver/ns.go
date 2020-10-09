package nameserver

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"path"
	"strings"
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
	serversset := make(map[string]struct{})
	return &GRPCServer{servers: serversset, replicatecnt: 2}
}

func (g *GRPCServer) pickServers(n int) ([]string, error) {
	if len(g.servers) < n {
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

	for i = 0; i < len(fileservers); i++ {
		conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		_, err := f.Create(ctx, &fsapi.CreateRequest{Inode: inode})
		if err != nil {
			log.Print(err)
			goto cleanup
		}
	}

	dir.entries.Store(node, &file{inode: inode, fileservers: fileservers})
	return &nsapi.CreateResponse{}, nil

cleanup:
	for ; i >= 0; i-- {
		conn, _ := cache.GrpcDial(fileservers[i], grpcopts...)
		f := fsapi.NewFileServerClient(conn)
		f.Remove(ctx, &fsapi.RemoveRequest{Inode: inode})
	}
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

	_, ok := to.entries.Load(rreq.ToName)
	if ok {
		fmt := "Rename: Entry `%s' in `%s' already exists"
		return nil, errors.Wrapf(ErrEntryExists, fmt, rreq.ToName, rreq.ToDir)
	}

	ent_int, ok := from.entries.LoadAndDelete(rreq.FromName)
	if !ok {
		fmt := "Rename: Failed to find entry `%s' in `%s'"
		return nil, errors.Wrapf(ErrEntryNotExists, fmt, rreq.FromName, rreq.FromDir)
	}

	to.entries.Store(rreq.ToName, ent_int.(dirEntry))
	log.Println("Ok Rename", rreq)

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

	return &nsapi.MapFSResponse{Fsurls: file.fileservers, Inode: file.inode}, nil
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

	fullAddr := peer.Addr.String()
	addr := fullAddr[:strings.LastIndex(fullAddr, ":")]
	addr = fmt.Sprintf("%s:%d", addr, req.Port)

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
