#!/bin/bash -eu

cd `dirname $0`

PACKAGE_NAME=$(go list .)
GATEWAY_PACKAGE_PATH=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)

set -x

# Apply tools
export PATH=$(pwd)/bin:$PATH
protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --go_out=plugins=grpc:$GOPATH/src echo.proto
protoc -I. -I./vendor -I$GATEWAY_PACKAGE_PATH/third_party/googleapis --grpc-gateway_out=logtostderr=true:$GOPATH/src ./echo.proto
