package ast

import "fmt"

type typeInfo[T any] interface {
	SetTypeInfo(T)
	GetTypeInfo() T
}

type (
	Expr interface {
		Node
		// typeInfo[int]
		aExpr() // hack again
	}

	Pattern interface {
		Unify(expr Expr) map[*Name]Expr
	}

	CallExpr struct {
		Fun     Expr
		ArgList []Expr // nil means no arguments
		HasDots bool   // last argument is followed by ...
		expr
		// atype
	}

	// Placeholder for an expression that failed to parse
	// correctly and where we can't provide a better node.
	BadExpr struct {
		expr
	}

	// Value
	Name struct {
		Value string
		expr
	}

	Integer struct {
		Value int
		expr
	}

	Float struct {
		Value float64
		expr
	}

	// Name Type
	//      Type
	Field struct {
		Name *Name // nil means anonymous field/parameter (structs/parameters), or embedded element (interfaces)
		Type Expr  // field names declared in a list share the same Type (identical pointers)
		node
	}
)

type expr struct {
	node
	// typeAndValue // After typechecking, contains the results of typechecking this expression.
}

func (*expr) aExpr() {}

// Patterns
func (field *Field) Unify(expr Expr) map[*Name]Expr {
	vs := map[*Name]Expr{}
	vs[field.Name] = expr
	return vs
}

func (name *Name) Unify(expr Expr) map[*Name]Expr {
	vs := map[*Name]Expr{}
	vs[name] = expr
	return vs
}

func (callExpr *CallExpr) Unify(expr Expr) map[*Name]Expr {
	vs := map[*Name]Expr{}
	return vs
}

func (intExpr *Integer) Unify(expr Expr) map[*Name]Expr {
	return map[*Name]Expr{}
}

func (floatExpr *Float) Unify(expr Expr) map[*Name]Expr {
	return map[*Name]Expr{}
}

// Format print expressions
func (name *Name) String() string {
	return fmt.Sprintf("%s", name.Value)
}

func (callExpr *CallExpr) String() string {
	if len(callExpr.ArgList) == 0 {
		// fmt.Println("???")
		return fmt.Sprintf("%s", callExpr.Fun)
	}
	return fmt.Sprintf("(%s %s)", callExpr.Fun, callExpr.ArgList)
}

func (intExpr *Integer) String() string {
	return fmt.Sprintf("%d", intExpr.Value)
}

func (floatExpr *Float) String() string {
	return fmt.Sprintf("%f", floatExpr.Value)
}
