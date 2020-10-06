package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Mexator/Go-vno/pkg/fileserver"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	"google.golang.org/grpc"
)

var (
	port       = flag.Uint64("p", 8080, "Port for grpc file server")
	host       = flag.String("h", "", "Hostname for grpc file server")
	storageDir = flag.String(
		"f",
		"/var/data",
		"Folder to stored files on the server",
	)
)

func main() {
	flag.Parse()

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
