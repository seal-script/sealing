package ast

import "fmt"

type Type interface {
	Expr
}

type (
	FuncType struct {
		Types []Type
		atype
	}
)

type atype struct{}

func (*atype) aType() {}

// func (*atype) aExpr() {}

func (*CallExpr) aType() {}

// func (t *BaseType) String() string {
// 	return fmt.Sprintf("%v", t.expr)
// }

func (t *FuncType) String() string {
	return fmt.Sprintf("(-> %s)", t.Types)
}
