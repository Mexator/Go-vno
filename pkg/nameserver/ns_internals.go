package nameserver

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"

	fsapi "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	"github.com/Mexator/Go-vno/pkg/cache"
)

// This file contains implementation of internal structures needed for
// work of nameserver

type (
	dirEntry interface {
		AsNode() *nsapi.Node
		Remove(context.Context) error
		replicate(addr string, g *GRPCServer)
	}

	directory struct {
		name    string
		entries sync.Map // map[string]dirEntry
	}
	file struct {
		name string

		inode       string
		fileservers sync.Map // map[string]struct{}
	}

	root struct{ dir directory }
)

func (d *directory) AsNode() *nsapi.Node { return &nsapi.Node{IsDir: true} }
func (f *file) AsNode() *nsapi.Node      { return &nsapi.Node{IsDir: false} }

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

func (f *file) Remove(ctx context.Context) error {
	var err error
	f.fileservers.Range(func(fsurl, value interface{}) bool {
		conn, e := cache.GrpcDial(fsurl.(string), grpcopts...)
		if err != nil {
			err = e
		}
		fs := fsapi.NewFileServerClient(conn)

		_, e = fs.Remove(ctx, &fsapi.RemoveRequest{Inode: f.inode})
		if err != nil {
			err = e
		}
		return true
	})
	return err
}

func (r *root) lookup(path string) (dirEntry, error) {
	sep := string(os.PathSeparator)
	parts := strings.Split(path, sep)

	if parts[0] != "" { // starts with /
		return nil, ErrDirNotExists
	}
	parts = parts[1:]

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
			return nil, ErrDirIsFile
		}
	}

	return d, nil
}

func (r *root) lookup_dir(path string) (*directory, error) {
	ent, err := r.lookup(path)
	if err != nil {
		return nil, err
	}
	dir, ok := ent.(*directory)
	if !ok {
		return nil, ErrDirIsFile
	}
	return dir, nil
}

func (f *file) replicate(addr string, g *GRPCServer) {
	var err error = errors.New("")
	var s []string

	for err != nil {
		s, err = g.pickServers(1)
		if err == nil {
			_, ok := f.fileservers.Load(s[0])
			if ok {
				err = errors.New("")
			}
		}
	}

	var old, new fsapi.FileServerClient

	f.fileservers.Delete(addr)

	f.fileservers.Range(func(k, _ interface{}) bool {
		conn, err := cache.GrpcDial(k.(string), grpcopts...)
		if err != nil {
			return false
		}
		old = fsapi.NewFileServerClient(conn)
		return false
	})
	if old == nil {
		// Hope for the best
		return
	}
	conn, err := cache.GrpcDial(s[0], grpcopts...)
	if err != nil {
		// Hope for the best
		return
	}
	new = fsapi.NewFileServerClient(conn)

	ctx := context.Background()
	szr, err := old.Size(ctx, &fsapi.SizeRequest{Inode: f.inode})
	if err != nil {
		return
	}
	sz := szr.Size

	var readsize uint64 = 2 * 1024 * 1024
	var i uint64

	for i = 0; i < sz; i += readsize {
		if i+readsize >= sz {
			readsize = sz - i - 1
		}

		rreq := &fsapi.ReadRequest{
			Inode:  f.inode,
			Offset: i,
			Size:   readsize,
		}
		rresp, err := old.Read(ctx, rreq)
		if err != nil {
			return
		}

		wreq := &fsapi.WriteRequest{
			Inode:   f.inode,
			Offset:  i,
			Content: rresp.Content,
		}
		_, err = new.Write(ctx, wreq)
		if err != nil {
			return
		}
	}

	f.fileservers.Store(s[0], struct{}{})
}

func (d *directory) replicate(addr string, g *GRPCServer) {
	d.entries.Range(func(_, eint interface{}) bool {
		eint.(dirEntry).consensus(addr, g)
		return true
	})
}
