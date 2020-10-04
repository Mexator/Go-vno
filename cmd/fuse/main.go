package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	dfs "github.com/Mexator/Go-vno/pkg/fuse"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s NAMESERVERURL MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}
	nsurl := flag.Arg(0)
	mountpoint := flag.Arg(1)

	c, err := fuse.Mount(mountpoint, fuse.FSName("govno"), fuse.Subtype("dfs"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, dfs.FS{Nsurl: nsurl})
	if err != nil {
		log.Fatal(err)
	}
}
