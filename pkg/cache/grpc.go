package cache

import (
	"sync"

	"google.golang.org/grpc"
)

var connections sync.Map //map[string]*grpc.ClientConn

// GrpcDial caches grpc connections if they are not already present
func GrpcDial(fsurl string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	clint, ok := connections.Load(fsurl)
	if ok {
		return clint.(*grpc.ClientConn), nil
	}
	conn, err := grpc.Dial(fsurl, opts...)
	if err != nil {
		return nil, err
	}
	connections.Store(fsurl, conn)
	return conn, nil
}
