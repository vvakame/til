// +build tools

package multi_schema

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "github.com/99designs/gqlgen"
	_ "golang.org/x/lint"
	_ "golang.org/x/tools"
)
