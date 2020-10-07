package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Mexator/Go-vno/pkg/fileserver"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"

	"google.golang.org/grpc"
)

var (
	port       = flag.Uint64("p", 2000, "Port for grpc file server")
	host       = flag.String("h", "", "Hostname for grpc file server")
	storageDir = flag.String(
		"f",
		"/var/data",
		"Folder to stored files on the server",
	)
	serverURL = flag.String(
		"s",
		"nameserver:3000",
		"hostname of the nameserver")
)

func main() {
	flag.Parse()

	s := grpc.NewServer()
	srv, err := fileserver.MakeFileServer(*storageDir)
	if err != nil {
		log.Fatal(err)
	}
	api.RegisterFileServerServer(s, &srv)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	go attachToNS(*serverURL)
	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func attachToNS(serverAddress string) {
	var grpcopts = []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddress, grpcopts...)
	if err != nil {
		log.Fatal("Can not connect to name server")
	}
	client := nsapi.NewNameServerClient(conn)
	client.ConnectFileServer(context.Background(), &nsapi.ConnectRequest{})
}
