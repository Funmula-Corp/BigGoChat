.PHONY: build

export PATH:=$(PATH):$(shell go env GOPATH)/bin

build:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./service.proto
