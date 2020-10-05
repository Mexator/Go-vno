generate:
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/nameserver --go_opt=paths=source_relative api/proto/nameserver.proto
	protoc -I api/proto --go_out=plugins=grpc:pkg/api/fileserver --go_opt=paths=source_relative api/proto/fileserver.proto

.PHONY: generate
