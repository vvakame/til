#!/bin/bash -eu

cd "$(dirname "$0")"

GATEWAY_PACKAGE_PATH=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)
export GATEWAY_PACKAGE_PATH

set -x

# Apply tools
PATH=$(pwd)/bin:$(pwd):$PATH
export PATH
command -v protoc-gen-go protoc-gen-grpc-gateway protoc-gen-gqlgen prototool

# statik
go generate ./cmd/protoc-gen-gqlgen

rm -rf ./echopb ./todopb
# mkdir -p ./echopb ./todopb

envsubst <prototool.base.yaml >prototool.yaml

prototool format -w
prototool lint
prototool generate

protoc -I=./proto -I="$GATEWAY_PACKAGE_PATH/third_party/googleapis" \
    --gqlgen_out=:./echopb \
    ./proto/echopb/echo.proto
protoc -I=./proto -I="$GATEWAY_PACKAGE_PATH/third_party/googleapis" \
    --gqlgen_out=:./todopb \
    ./proto/todopb/todo.proto

goimports -w ./**/*.gql.go

# wire & gqlgen
go generate ./...
