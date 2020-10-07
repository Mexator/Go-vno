package fuse

import (
	"context"
	"os"
	"path"
	"sync"
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
		dir  *Dir
		name string
		conn nsapi.NameServerClient
		sync.Map
	}
	File struct {
		dir  *Dir
		name string
		conn nsapi.NameServerClient
	}
	Node interface {
		rename(name string)
		changedir(dir *Dir)
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
	return &Dir{conn: nsapi.NewNameServerClient(conn)}, nil
}

func (d *Dir) getPath() string {
	if d.dir == nil {
		return string(os.PathSeparator)
	} else {
		return path.Join(d.dir.getPath(), d.name)
	}
}

func (f *File) getPath() string    { return path.Join(f.dir.getPath(), f.name) }
func (f *File) rename(name string) { f.name = name }
func (d *Dir) rename(name string)  { d.name = name }
func (f *File) changedir(dir *Dir) { f.dir = dir }
func (d *Dir) changedir(dir *Dir)  { d.dir = dir }

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0o555
	a.Size = 4096 // Why not lol
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	resp, err := d.conn.ReadDirAll(ctx, &nsapi.ReadDirAllRequest{Path: d.getPath()})
	if err != nil {
		return nil, errors.Wrap(err, "Failed lookup in name server")
	}

	for _, n := range resp.Nodes {
		if n.Name != name {
			continue
		}

		var node fs.Node

		if n.IsDir {
			node = &Dir{dir: d, name: name, conn: d.conn}
		} else {
			node = &File{dir: d, name: name, conn: d.conn}
		}

		d.Store(name, node)
		return node, nil
	}

	return nil, syscall.ENOENT
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	resp, err := d.conn.ReadDirAll(ctx, &nsapi.ReadDirAllRequest{Path: d.getPath()})
	if err != nil {
		return nil, errors.Wrap(err, "Failed lookup in name server")
	}

	ret := []fuse.Dirent{}

	for _, n := range resp.Nodes {
		dirent := fuse.Dirent{Name: n.Name}

		if n.IsDir {
			dirent.Type = fuse.DT_Dir
		} else {
			dirent.Type = fuse.DT_File
		}

		ret = append(ret, dirent)
	}

	return ret, nil
}

func (f *File) fsGetSize(ctx context.Context, fsurl, inode string) (uint64, error) {
	conn, err := cache.GrpcDial(fsurl, grpcopts...)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to connect to fileserver")
	}
	cl := fsapi.NewFileServerClient(conn)

	szresp, err := cl.Size(ctx, &fsapi.SizeRequest{Inode: inode})
	if err != nil {
		return 0, errors.Wrap(err, "Failed to retrieve size")
	}

	return szresp.Size, nil
}

func (f *File) getSize(ctx context.Context) (uint64, error) {
	resp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.getPath()})
	if err != nil {
		return 0, errors.Wrap(err, "Failed MapFS")
	}

	var sz uint64

	for _, url := range resp.Fsurls {
		sz, err = f.fsGetSize(ctx, url, resp.Inode)
		if err == nil {
			return sz, nil
		}
	}

	return 0, err
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	sz, err := f.getSize(ctx)
	if err != nil {
		return errors.Wrap(err, "Attr:")
	}

	a.Mode = 0o777 // No mode management lol
	a.Size = sz
	return nil
}

func (f *File) Read(ctx context.Context, fsurl, inode string) ([]byte, error) {
	conn, err := cache.GrpcDial(fsurl, grpcopts...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to fileserver")
	}
	cl := fsapi.NewFileServerClient(conn)

	szresp, err := cl.Size(ctx, &fsapi.SizeRequest{Inode: inode})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get size")
	}

	rresp, err := cl.Read(ctx, &fsapi.ReadRequest{Inode: inode, Offset: 0, Size: szresp.Size})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read from fileserver")
	}

	return rresp.Content, nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	resp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.getPath()})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve fileserver from nameserver")
	}

	var cont []byte

	for _, url := range resp.Fsurls {
		cont, err = f.Read(ctx, url, resp.Inode)
		if err == nil {
			return cont, nil
		}
	}

	return nil, err
}

func (f *File) Write(
	ctx context.Context,
	req *fuse.WriteRequest,
	resp *fuse.WriteResponse,
) error {
	mresp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.getPath()})
	resp.Size = 0
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve fileserver from nameserver")
	}

	for _, fsurl := range mresp.Fsurls {
		_, err = cache.GrpcDial(fsurl, grpcopts...)
		if err != nil {
			return errors.Wrap(err, "Failed to connect fileserver")
		}
	}

	for i := 0; i < len(mresp.Fsurls); i++ {
		conn, _ := cache.GrpcDial(mresp.Fsurls[i], grpcopts...)
		cl := fsapi.NewFileServerClient(conn)

		_, err = cl.Write(ctx, &fsapi.WriteRequest{
			Inode:   mresp.Inode,
			Offset:  uint64(req.Offset),
			Content: req.Data,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to write to fileserver")
		}
	}

	resp.Size = len(req.Data)
	return nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	rreq := &nsapi.RemoveRequest{Path: path.Join(d.getPath(), req.Name)}
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
		Path:  path.Join(d.getPath(), req.Name),
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
		Path:  path.Join(d.getPath(), req.Name),
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

func (from *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, to_node fs.Node) error {
	to := to_node.(*Dir)
	rreq := &nsapi.RenameRequest{
		FromDir:  from.getPath(),
		FromName: req.OldName,
		ToDir:    to.getPath(),
		ToName:   req.NewName,
	}
	_, err := from.conn.Rename(ctx, rreq)
	if err != nil {
		return errors.Wrap(err, "Failed to rename:")
	}

	v, l := from.Map.LoadAndDelete(req.OldName)
	if l {
		v.(Node).rename(req.NewName)
		v.(Node).changedir(to)
		to.Map.Store(req.NewName, v)
	}

	return nil
}
