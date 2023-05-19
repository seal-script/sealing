package ast

import "fmt"

type Type interface {
	Expr
}

type (
	FuncType struct {
		Context []Field
		Types   []Type
		atype
		expr
	}
)

type atype struct{}

func (*atype) aType() {}

func (*CallExpr) aType() {}

// func (t *BaseType) String() string {
// 	return fmt.Sprintf("%v", t.expr)
// }

func (t *FuncType) String() string {
	return fmt.Sprintf("(-> %s)", t.Types)
}
