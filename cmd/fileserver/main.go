package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	"github.com/Mexator/Go-vno/pkg/fileserver"
	"github.com/Mexator/Go-vno/pkg/utils"

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
	nsurl, storage := flag.Arg(0), flag.Arg(1)

	s := utils.GrpcServer()
	srv, err := fileserver.MakeFileServer(storage, nsurl)
	if err != nil {
		log.Fatal(err)
	}
	api.RegisterFileServerServer(s, &srv)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	go attachToNS(nsurl)
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
	client.ConnectFileServer(context.Background(), &nsapi.ConnectRequest{Port: int32(*port)})
}
