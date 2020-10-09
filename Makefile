MOUNTDIR=mnt
NFILESERVERS=3
DCC=docker-compose -f docker/docker-compose.yml

all: generate fileserver nameserver client

generate: pkg/api/nameserver/nameserver.pb.go pkg/api/fileserver/fileserver.pb.go

pkg/api/nameserver/nameserver.pb.go: api/proto/nameserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/nameserver --go_opt=paths=source_relative api/proto/nameserver.proto

pkg/api/fileserver/fileserver.pb.go: api/proto/fileserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/fileserver --go_opt=paths=source_relative api/proto/fileserver.proto

build/fileserver: cmd/fileserver/* pkg/fileserver/* pkg/api/*
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/fileserver ./cmd/fileserver
fileserver: build/fileserver
	docker build -f docker/fileserver/Dockerfile -t mexator/fileserver .

build/nameserver: cmd/nameserver/* pkg/nameserver/* pkg/api/*
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/nameserver ./cmd/nameserver
nameserver: build/nameserver
	docker build -f docker/nameserver/Dockerfile -t mexator/nameserver .

build/client: cmd/fuse/* pkg/fuse/*
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/client ./cmd/fuse
client: build/client
	docker build -f docker/client/Dockerfile -t mexator/client .

local_run: all
	$(DCC) down || true
	$(DCC) up -d --scale fileserver=$(NFILESERVERS)
	while [ $$($(DCC) logs | grep '"grpc.method": "ConnectFileServer"' | grep OK | wc -l) = $(NFILESERVERS) ]; do :; done
	$(DCC) exec client sh

local_logs:
	$(DCC) logs

.PHONY: generate fileserver nameserver all local_run
