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
	port = flag.Uint64("p", 2000, "Port for grpc file server")
	host = flag.String("h", "", "Hostname for grpc file server")
)

func main() {
	flag.Usage = func() {
		fmt.Println("fileserver NSADDR STORAGEDIR")
		fmt.Println("NSADDR - hostname of the nameserver")
		fmt.Println("STORAGEDIR - Folder to stored files on file server")
		flag.PrintDefaults()
	}
	flag.Parse()

	s := grpc.NewServer()
	srv, err := fileserver.MakeFileServer(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	api.RegisterFileServerServer(s, &srv)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	go attachToNS(flag.Arg(0))
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
