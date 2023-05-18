package ast

import (
	"fmt"
	"reflect"
)

type (
	// Top-level Declarations
	Decl interface {
		Node
		aDecl() // Just for constraint... golang hack!
	}

	//              Path
	// LocalPkgName Path
	// import Data.Functor (fmap)
	ImportDecl struct {
		Alias *Name
		Path  string
		decl
	}

	ModuleDecl struct {
		Name       *Name
		ExportList []Expr
		decl
	}

	// Name Type
	// double : Int -> Int
	// List : Type -> Type
	TypeDecl struct {
		Name *Name
		Type Type
		decl
	}

	// func          Name Type { Body }
	// func          Name Type
	// func Receiver Name Type { Body }
	// func Receiver Name Type
	// f a = a + 1
	// f a = x + 1 where {
	//     x = a * a
	// }
	// f a = let {
	//     double :: Int
	//     double x = x * x
	//     x = double a
	// } in x + 1
	FuncDecl struct {
		Name   *Name
		Type   Type
		Params []Pattern
		Body   Expr // Purely functional! nil means no body (forward declaration)
		decl
	}

	// data Nat {
	//     Zero : Nat
	//     Succ : Nat -> Nat
	// }
	EnumDecl struct {
		Name *Name
		Cons []TypeDecl
		decl
	}

	// seal Monoid a {
	//     Empty : a
	//     (<>)  : a -> a -> a
	// }
	//
	// seal Functor f {
	//     map : (a -> b) -> f a -> f b
	//     (<$>) = map
	// }
	SealDecl struct {
		Name   *Name
		fields []TypeDecl
		decl
	}
)

type decl struct{ node }

func (*decl) aDecl() {}

// Format declarations
func (funcDecl *FuncDecl) String() string {
	return fmt.Sprintf(
		`Function Decl {
			Name: %s,
			Type: %s : %v,
			Params: %s,
			Body: %s : %v,
		}`,
		funcDecl.Name,
		funcDecl.Type,
		reflect.TypeOf(funcDecl.Type),
		funcDecl.Params,
		funcDecl.Body,
		reflect.TypeOf(funcDecl.Body),
	)
}

func (typeDecl *TypeDecl) String() string {
	return fmt.Sprintf(
		`Type Decl {
			Name: %s,
			Type: %s,
		}`,
		typeDecl.Name,
		typeDecl.Type,
	)
}
