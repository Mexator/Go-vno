package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	nsapi "github.com/Mexator/Go-vno/pkg/api/nameserver"
	ns "github.com/Mexator/Go-vno/pkg/nameserver"
	"google.golang.org/grpc"
)

var (
	port = flag.Uint64("p", 8080, "Port for grpc name server")
	host = flag.String("h", "", "Hostname for grpc name server")
)

func main() {
	flag.Parse()
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	s := grpc.NewServer()
	srv := ns.NewServer(lines)
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
