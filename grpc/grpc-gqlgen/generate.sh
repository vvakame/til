#!/bin/bash -eu

cd `dirname $0`

PACKAGE_NAME=$(go list -m)
GATEWAY_PACKAGE_PATH=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)

set -x

# Apply tools
export PATH=$(pwd)/bin:$(pwd):$PATH
command -v protoc-gen-go protoc-gen-grpc-gateway protoc-gen-gqlgen
protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --go_out=plugins=grpc,paths=source_relative:./echopb echo.proto
protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:./echopb ./echo.proto

protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --go_out=paths=source_relative:./ proto-extentions/gqlgen.proto
protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --gqlgen_out=:./echopb ./echo.proto
