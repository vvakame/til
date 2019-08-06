// +build tools

package metago

// from https://github.com/golang/go/issues/25922#issuecomment-412992431

import (
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
