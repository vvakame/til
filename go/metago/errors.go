package metago

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
)

type ErrorLevel int

const (
	ErrorLevelError ErrorLevel = iota
	ErrorLevelWarning
	ErrorLevelNotice
	ErrorLevelDebug
)

var _ error = (NodeErrors)(nil)

type NodeErrors []*NodeError

func (nErrs NodeErrors) Error() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	for idx, nErr := range nErrs {
		if idx != 0 {
			_, _ = buf.Write([]byte("\n"))
		}
		pos := nErr.Fset.Position(nErr.Node.Pos())
		errPath, err := filepath.Rel(cwd, pos.Filename)
		if err != nil {
			panic(err)
		}
		errText := fmt.Sprintf("%s:%d:%d: %s", errPath, pos.Line, pos.Column, nErr.Message)
		_, _ = buf.Write([]byte(errText))
	}

	return buf.String()
}

type NodeError struct {
	ErrorLevel ErrorLevel
	Fset       *token.FileSet
	Node       ast.Node
	Message    string
}
