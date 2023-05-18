package backand

import "github.com/seal-script/sealing/ast"

type Ir interface{}

type CodeGen interface {
	Gen(ast ast.Node) Ir
}
