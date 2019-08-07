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
	Package       *packages.Package
	File          *ast.File
	FilePath      string
	GeneratedCode string
	Errors        NodeErrors
}

type Processor interface {
	Process() (*Result, error)
}

func NewProcessor(cfg *Config) (Processor, error) {
	p := &metaProcessor{
		cfg:          cfg,
		removeNodes:  make(map[ast.Node]bool),
		replaceNodes: make(map[ast.Node]ast.Node),
		valueMapping: make(map[*ast.Object]ast.Expr),
		fieldMapping: make(map[*ast.Object]ast.Expr),
	}
	return p, nil
}
