package fuse

import (
	"context"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"google.golang.org/grpc"

	fsapi "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
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
		size uint64
		conn nsapi.NameServerClient
	}
)

func (f FS) Root() (fs.Node, error) {
	conn, err := grpc.Dial(f.Nsurl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return Dir{path: "/", conn: nsapi.NewNameServerClient(conn)}, nil
}

func getInode(fname string) uint64 {
	h := fnv.New64()
	_, err := h.Write([]byte(fname))
	if err != nil {
		log.Fatal(err)
	}
	return h.Sum64()
}

func (d Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = getInode(d.path)
	a.Mode = os.ModeDir | 0o555
	a.Size = 4096 // Why not lol
	return nil
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	resp, err := d.conn.Lookup(ctx, &nsapi.LookupRequest{Path: d.path})
	if err != nil {
		return nil, err
	}

	for _, n := range resp.Nodes {
		if n.Name != name {
			continue
		}

		path := filepath.Join(d.path, n.Name)
		if n.IsDir {
			return Dir{path: path, conn: d.conn}, nil
		} else {
			return File{path: path, conn: d.conn, size: n.Size}, nil
		}
	}

	return nil, syscall.ENOENT
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	resp, err := d.conn.Lookup(ctx, &nsapi.LookupRequest{Path: d.path})
	if err != nil {
		return nil, err
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

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = getInode(f.path)
	a.Mode = 0o777 // No mode management lol
	a.Size = f.size
	return nil
}

// caching fileserver connections
var fileservers sync.Map //map[string]fsapi.FileServerClient

func getFileServer(fsurl string) (fsapi.FileServerClient, error) {
	clint, ok := fileservers.Load(fsurl)
	if ok {
		return clint.(fsapi.FileServerClient), nil
	}
	conn, err := grpc.Dial(fsurl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := fsapi.NewFileServerClient(conn)
	fileservers.Store(fsurl, c)
	return c, nil
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	resp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.path})
	if err != nil {
		return nil, err
	}
	cl, err := getFileServer(resp.Fsurl)
	if err != nil {
		return nil, err
	}

	rresp, err := cl.Read(ctx, &fsapi.ReadRequest{
		Inode:  resp.Inode,
		Offset: 0,
		Size:   f.size,
	})
	if err != nil {
		return nil, err
	}

	return rresp.Content, nil
}

func (f File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	mresp, err := f.conn.MapFS(ctx, &nsapi.MapFSRequest{Path: f.path})
	if err != nil {
		resp.Size = 0
		return err
	}
	cl, err := getFileServer(mresp.Fsurl)
	if err != nil {
		resp.Size = 0
		return err
	}

	_, err = cl.Write(ctx, &fsapi.WriteRequest{
		Inode:   mresp.Inode,
		Offset:  uint64(req.Offset),
		Content: req.Data,
	})
	if err != nil {
		resp.Size = 0
		return err
	}
	resp.Size = len(req.Data)

	return nil
}
