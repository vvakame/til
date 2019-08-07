package metago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

type ErrorLevel int

const (
	ErrorLevelError ErrorLevel = iota
	ErrorLevelWarning
	ErrorLevelNotice
	ErrorLevelDebug
)

var _ error = (NodeErrors)(nil)
var _ json.Marshaler = (NodeErrors)(nil)

type NodeErrors []*NodeError

type NodeError struct {
	ErrorLevel ErrorLevel
	Fset       *token.FileSet `json:"-"`
	Node       ast.Node       `json:"-"`
	Message    string
}

func (nErrs NodeErrors) Error() string {
	var buf bytes.Buffer
	for idx, nErr := range nErrs {
		if idx != 0 {
			_, _ = buf.Write([]byte("\n"))
		}
		_, _ = buf.Write([]byte(nErr.Error()))
	}

	return buf.String()
}

func (nErrs NodeErrors) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte("["))
	for idx, nErr := range nErrs {
		if idx != 0 {
			buf.Write([]byte(","))
		}
		b, err := nErr.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.Write([]byte("]"))
	return buf.Bytes(), nil
}

func (nErr *NodeError) Error() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pos := nErr.Fset.Position(nErr.Node.Pos())
	errPath, err := filepath.Rel(cwd, pos.Filename)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s:%d:%d: %s", errPath, pos.Line, pos.Column, nErr.Message)
}

func (nErr *NodeError) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte("{"))
	{
		buf.Write([]byte(strconv.Quote("level")))
		buf.Write([]byte(":"))
		switch nErr.ErrorLevel {
		case ErrorLevelError:
			buf.Write([]byte(strconv.Quote("error")))
		case ErrorLevelWarning:
			buf.Write([]byte(strconv.Quote("warning")))
		case ErrorLevelNotice:
			buf.Write([]byte(strconv.Quote("notice")))
		case ErrorLevelDebug:
			buf.Write([]byte(strconv.Quote("debug")))
		default:
			buf.Write([]byte(strconv.Quote("unknown")))
		}
	}
	buf.Write([]byte(","))
	{
		buf.Write([]byte(strconv.Quote("message")))
		buf.Write([]byte(":"))
		buf.Write([]byte(strconv.Quote(nErr.Error())))
	}

	buf.Write([]byte("}"))
	return buf.Bytes(), nil
}
