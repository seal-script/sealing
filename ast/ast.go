package ast

import "github.com/seal-script/sealing/utils"

type Location = utils.Location

type Node interface {
	// Locate() returns the position associated with the node as follows:
	// 1) The position of a node representing a terminal syntax production
	//    (Name, BasicLit, etc.) is the position of the respective production
	//    in the source.
	// 2) The position of a node representing a non-terminal production
	//    (IndexExpr, IfStmt, etc.) is the position of a token uniquely
	//    associated with that production; usually the left-most one
	//    ('[' for IndexExpr, 'if' for IfStmt, etc.)
	Locate() Location
	aNode() // Just for constraint... golang hack!
}

type node struct {
	// commented out for now since not yet used
	// doc  *Comment // nil means no comment(s) attached
	Location Location
}

func (n *node) Locate() Location { return n.Location }
func (*node) aNode()             {}

// package PkgName; DeclList[0], DeclList[1], ...
type File struct {
	// Pragma   Pragma
	PkgName  *Name
	DeclList []Decl
	EOF      Location
	node
}
