all: generate fileserver nameserver fuse

generate: pkg/api/nameserver/nameserver.pb.go pkg/api/fileserver/fileserver.pb.go

pkg/api/nameserver/nameserver.pb.go: api/proto/nameserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/nameserver --go_opt=paths=source_relative api/proto/nameserver.proto

pkg/api/fileserver/fileserver.pb.go: api/proto/fileserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/fileserver --go_opt=paths=source_relative api/proto/fileserver.proto

fileserver:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/fileserver ./cmd/fileserver
	docker build -f docker/fileserver/Dockerfile -t mexator/fileserver .

nameserver:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/nameserver ./cmd/nameserver
	docker build -f docker/nameserver/Dockerfile -t mexator/nameserver .

fuse:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/client ./cmd/fuse


.PHONY: generate fileserver nameserver fuse all
