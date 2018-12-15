#!/bin/sh -eux

cd `dirname $0`

# build tools
rm -rf bin/
mkdir bin/

go mod download
# go mod tidy
# go generate のため
go mod vendor

export GOBIN=`pwd -P`/bin
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
