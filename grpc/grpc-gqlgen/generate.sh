#!/bin/bash -eu

cd "$(dirname "$0")"

GATEWAY_PACKAGE_PATH=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)

set -x

# Apply tools
PATH=$(pwd)/bin:$(pwd):$PATH
export PATH
command -v protoc-gen-go protoc-gen-grpc-gateway protoc-gen-gqlgen

rm -rf ./echopb ./todopb
mkdir -p ./echopb ./todopb

go generate ./cmd/protoc-gen-gqlgen

protoc -I=. -I="$GATEWAY_PACKAGE_PATH/third_party/googleapis" \
    --go_out=paths=source_relative:./ \
    gqlgen-proto/options.proto

protoc -I=. -I="$GATEWAY_PACKAGE_PATH/third_party/googleapis" \
    --go_out=plugins=grpc,paths=source_relative:./echopb \
    --grpc-gateway_out=logtostderr=true,paths=source_relative:./echopb \
    --gqlgen_out=:./echopb \
    echo.proto

protoc -I=. -I="$GATEWAY_PACKAGE_PATH/third_party/googleapis" \
    --go_out=plugins=grpc,paths=source_relative:./todopb \
    --grpc-gateway_out=logtostderr=true,paths=source_relative:./todopb \
    --gqlgen_out=:./todopb \
    todo.proto

go generate ./...
goimports -w ./**/*.gql.go
