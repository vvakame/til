package metago

import "go/ast"

var _ ast.Visitor = (astVisitorFunc)(nil)

type astVisitorFunc func(node ast.Node) bool

func (v astVisitorFunc) Visit(node ast.Node) (w ast.Visitor) {
	if v(node) {
		return v
	}

	return nil
}
