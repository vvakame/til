package metago

import "go/ast"

type NodeError struct {
	Node    ast.Node
	Message string
}
