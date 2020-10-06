package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	ns "github.com/Mexator/Go-vno/pkg/nameserver"
	"google.golang.org/grpc"
)

var (
	port = flag.Uint64("p", 3000, "Port for grpc name server")
	host = flag.String("h", "", "Hostname for grpc name server")
)

func main() {
	flag.Parse()

	var fileservers []string

	for i := 0; i < flag.NArg(); i++ {
		fileservers = append(fileservers, flag.Arg(i))
	}

	s := grpc.NewServer()
	srv := ns.NewServer(fileservers)
	nsapi.RegisterNameServerServer(s, srv)

	log.Printf("Server is listening on %s:%d", *host, *port)
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
