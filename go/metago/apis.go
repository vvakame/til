package metago

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func ValueOf(interface{}) Value {
	panic("in meta context!")
}

type Value interface {
	Fields() []Field
}

type Field interface {
	Name() string
	StructTagGet(string) string
	Value() interface{}
}

type Config struct {
	TargetPackages []string
}

type Result struct {
	CompileErrors []packages.Error
	Results       []*FileResult
}

type FileResult struct {
	Package           *packages.Package
	File              *ast.File
	BaseFilePath      string
	GeneratedFilePath string
	GeneratedCode     string
	Errors            NodeErrors
}

type Processor interface {
	Process(cfg *Config) (*Result, error)
}

func NewProcessor() (Processor, error) {
	return &metaProcessor{}, nil
}
