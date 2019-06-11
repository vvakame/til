// +build tools

package main

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "github.com/golang/protobuf"
	_ "github.com/grpc-ecosystem/grpc-gateway"
	_ "github.com/rakyll/statik"
)
