package nameserver

import (
	"context"
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
