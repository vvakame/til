#!/bin/sh -eux

cd "$(dirname "$0")"

# build tools
rm -rf bin/
mkdir bin/

go mod download
# go mod tidy

GOBIN="$(pwd -P)/bin"
export GOBIN
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go install github.com/uber/prototool/cmd/prototool
go install github.com/rakyll/statik
go install github.com/google/wire/cmd/wire
go install github.com/99designs/gqlgen
