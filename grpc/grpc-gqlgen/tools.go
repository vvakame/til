// +build tools

package main

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/uber/prototool/cmd/prototool"
	_ "github.com/rakyll/statik"
)
