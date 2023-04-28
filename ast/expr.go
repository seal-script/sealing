package ast

type typeInfo[T any] interface {
	SetTypeInfo(T)
	GetTypeInfo() T
}

type (
	Expr interface {
		Node
		typeInfo[int]
		aExpr() // hack again
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
