package ast

import (
	"fmt"
)

func ToString(node Node) string {
	if e, ok := node.(Expr); ok {
		return ExprToString(e)
	} else if s, ok := node.(Stmt); ok {
		return StmtToString(s)
	} else {
		panic(fmt.Sprintf("ast stringer: can't deal with %T %v", node, node))
	}
}
