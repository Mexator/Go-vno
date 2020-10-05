package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Mexator/Go-vno/pkg/fileserver"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	"google.golang.org/grpc"
)

var (
	port = flag.Uint64("p", 8080, "Port for grpc name server")
	host = flag.String("h", "", "Hostname for grpc name server")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "  %s FILEDIR\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	s := grpc.NewServer()
	srv, err := fileserver.MakeFileServer(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	api.RegisterFileServerServer(s, &srv)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
