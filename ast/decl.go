package ast

type (
	// Top-level Declarations
	Decl interface {
		Node
		aDecl() // Just for constraint... golang hack!
	}

	//              Path
	// LocalPkgName Path
	ImportDecl struct {
		Alias *Name
		Path  string
		decl
	}

	// Name Type
	TypeDecl struct {
		Name *Name
		Type *Type
		decl
	}

	// func          Name Type { Body }
	// func          Name Type
	// func Receiver Name Type { Body }
	// func Receiver Name Type
	FuncDecl struct {
		Name   *Name
		Type   *Type
		Params *[]Field
		Body   *Expr // Purely functional! nil means no body (forward declaration)
		decl
	}
)

type decl struct{ node }

func (*decl) aDecl() {}
