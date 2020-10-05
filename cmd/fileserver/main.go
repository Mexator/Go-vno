package main

import (
	"flag"
	"log"
	"net"

	"github.com/Mexator/Go-vno/pkg/fileserver"

	api "github.com/Mexator/Go-vno/pkg/api/fileserver"
	"google.golang.org/grpc"
)

func main() {
	config := flag.String("config", "./config.json", "Path to JSON config file")
	flag.Parse()
	if err := startFileServer(*config); err != nil {
		log.Fatal(err)
	}
	log.Print("Server started")
}

func startFileServer(configPath string) error {
	s := grpc.NewServer()
	srv, err := fileserver.MakeFileServer(configPath)
	if err != nil {
		return err
	}
	api.RegisterFileServerServer(s, &srv)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	if err := s.Serve(listener); err != nil {
		return err
	}
	return nil
}
