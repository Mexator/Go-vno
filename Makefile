

generate:
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/nameserver --go_opt=paths=source_relative api/proto/nameserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/fileserver --go_opt=paths=source_relative api/proto/fileserver.proto

fileserver:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/fileserver ./cmd/fileserver
	docker build -f docker/fileserver/Dockerfile .
nameserver:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/nameserver ./cmd/nameserver
fuse:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/client ./cmd/fuse

.PHONY: generate fileserver
