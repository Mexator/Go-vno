package fuse

import (
	"context"
	"hash/fnv"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"google.golang.org/grpc"

	fsapi "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	"github.com/Mexator/Go-vno/pkg/cache"
)

type (
	FS struct {
		Nsurl string
	}

	Dir struct {
		path string
		conn nsapi.NameServerClient
	}

	File struct {
		path string
		conn nsapi.NameServerClient
	}
)

var (
	grpcopts = []grpc.DialOption{grpc.WithInsecure()}
)

func (f *FS) Root() (fs.Node, error) {
	conn, err := grpc.Dial(f.Nsurl, grpcopts...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed connect to name server")
	}
	return &Dir{path: "/", conn: nsapi.NewNameServerClient(conn)}, nil
}

func getInode(fname string) uint64 {
	h := fnv.New64()
	_, err := h.Write([]byte(fname))
	if err != nil {
		log.Fatal(err)
	}
	return h.Sum64()
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = getInode(d.path)
	a.Mode = os.ModeDir | 0o555
	a.Size = 4096 // Why not lol
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	resp, err := d.conn.ReadDirAll(ctx, &nsapi.ReadDirAllRequest{Path: d.path})
	if err != nil {
		return nil, errors.Wrap(err, "Failed lookup in name server")
	}

	for _, n := range resp.Nodes {
		if n.Name != name {
			continue
		}

		path := filepath.Join(d.path, n.Name)
		if n.IsDir {
			return &Dir{path: path, conn: d.conn}, nil
		} else {
			return &File{path: path, conn: d.conn}, nil
		}
	}

	return nil, syscall.ENOENT
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	resp, err := d.conn.ReadDirAll(ctx, &nsapi.ReadDirAllRequest{Path: d.path})
	if err != nil {
		return nil, errors.Wrap(err, "Failed lookup in name server")
	}

	ret := []fuse.Dirent{}

	for _, n := range resp.Nodes {
		dirent := fuse.Dirent{
			Inode: getInode(filepath.Join(d.path, n.Name)),
			Name:  n.Name,
		}

		if n.IsDir {
			dirent.Type = fuse.DT_Dir
		} else {
			dirent.Type = fuse.DT_File
		}

		ret = append(ret, dirent)
	}

	return ret, nil
}

func (f *File) getSize(ctx context.Context) (uint64, error) {
	resp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.path})
	if err != nil {
		return 0, errors.Wrap(err, "Failed MapFS")
	}

	conn, err := cache.GrpcDial(resp.Fsurl, grpcopts...)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to connect to fileserver")
	}
	cl := fsapi.NewFileServerClient(conn)

	szresp, err := cl.Size(ctx, &fsapi.SizeRequest{Inode: resp.Inode})
	if err != nil {
		return 0, errors.Wrap(err, "Failed to retrieve size")
	}

	return szresp.Size, nil
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	sz, err := f.getSize(ctx)
	if err != nil {
		return errors.Wrap(err, "Attr:")
	}

	a.Inode = getInode(f.path)
	a.Mode = 0o777 // No mode management lol
	a.Size = sz
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	resp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.path})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve fileserver from nameserver")
	}
	conn, err := cache.GrpcDial(resp.Fsurl, grpcopts...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to fileserver")
	}
	cl := fsapi.NewFileServerClient(conn)

	szresp, err := cl.Size(ctx, &fsapi.SizeRequest{Inode: resp.Inode})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get size")
	}

	rresp, err := cl.Read(ctx, &fsapi.ReadRequest{
		Inode:  resp.Inode,
		Offset: 0,
		Size:   szresp.Size,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read from fileserver")
	}

	return rresp.Content, nil
}

func (f *File) Write(
	ctx context.Context,
	req *fuse.WriteRequest,
	resp *fuse.WriteResponse,
) error {
	mresp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.path})
	if err != nil {
		resp.Size = 0
		return errors.Wrap(err, "Failed to retrieve fileserver from nameserver")
	}
	conn, err := cache.GrpcDial(mresp.Fsurl, grpcopts...)
	if err != nil {
		resp.Size = 0
		return errors.Wrap(err, "Failed to connect fileserver")
	}
	cl := fsapi.NewFileServerClient(conn)

	_, err = cl.Write(ctx, &fsapi.WriteRequest{
		Inode:   mresp.Inode,
		Offset:  uint64(req.Offset),
		Content: req.Data,
	})
	if err != nil {
		resp.Size = 0
		return errors.Wrap(err, "Failed to write to fileserver")
	}

	resp.Size = len(req.Data)
	return nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	rreq := &nsapi.RemoveRequest{Path: path.Join(d.path, req.Name)}
	_, err := d.conn.Remove(ctx, rreq)
	if err != nil {
		return errors.Wrap(err, "Failed to remove entry")
	}
	return nil
}

func (d *Dir) Create(
	ctx context.Context,
	req *fuse.CreateRequest,
	resp *fuse.CreateResponse,
) (fs.Node, fs.Handle, error) {
	creq := &nsapi.CreateRequest{
		Path:  path.Join(d.path, req.Name),
		IsDir: req.Mode.IsDir(),
	}
	_, err := d.conn.Create(ctx, creq)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create node")
	}

	node, err := d.Lookup(ctx, req.Name)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to retrieve created node")
	}

	return node, node, err
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	creq := &nsapi.CreateRequest{
		Path:  path.Join(d.path, req.Name),
		IsDir: true,
	}
	_, err := d.conn.Create(ctx, creq)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create node")
	}

	node, err := d.Lookup(ctx, req.Name)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve created node")
	}

	return node, err
}
