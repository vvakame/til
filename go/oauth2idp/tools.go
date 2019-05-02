// +build tools

package main

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "github.com/favclip/jwg"
	_ "github.com/favclip/qbg"
	_ "github.com/google/wire/cmd/wire"
	_ "golang.org/x/tools/cmd/goimports"
)
