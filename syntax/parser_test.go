package syntax

import (
	"bytes"
	"io"
	"testing"
)

// func NewParser(t *testing.T, in io.Reader) Parser {
// 	p := Parser{}
// 	p.Init(in, func(err error) {
// 		t.Log(err.Error())
// 	})
// 	p.next()
// 	return p
// }

// func TestParser(t *testing.T) {
// 	// Slice reader
// 	data := []byte(`(let x 测试gdfh 烤红薯烤豆腐)`)
// 	var in io.Reader = bytes.NewReader(data)

// 	p := NewParser(t, in)
// 	ast, err := p.ParseFile()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	t.Logf("%v", ast)
// }

// func TestParsingFuncDecl(t *testing.T) {
// 	data := []byte(`let x 7)`)
// 	var in io.Reader = bytes.NewReader(data)

// 	p := NewParser(t, in)
// 	p.next()
// 	res, err := p.ParseFuncDecl()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	t.Log("-------------------------\n")
// 	t.Logf("%#v", res)
// 	// fd := ast.FuncDecl(res)
// 	// t.Logf("%v", res.Name.Value)
// }

func TestParseFuncCallExpr(t *testing.T) {
	// data := []byte(`fact n = (add (fact ((-) n 1)) (fact (minus n 2)))`)
	data := []byte(`fact n = (*) n (fact ((-) n 1))`)
	var in io.Reader = bytes.NewReader(data)

	p := NewParser(t, in)
	// p.next()
	res, err := p.ParseDecl()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("-------------------------\n")
	t.Logf("%v", res)
	// fd := ast.FuncDecl(res)
	// t.Logf("%#v", res.Fun)
	// t.Logf("%#v", res.ArgList[0])
	// t.Logf("%#v", res.ArgList[0].(*ast.CallExpr).Fun)
}

func TestParseType(t *testing.T) {
	data := []byte(`map : (a -> b) -> f a -> f b`)
	var in io.Reader = bytes.NewReader(data)

	p := NewParser(t, in)
	// p.next()
	res, err := p.ParseDecl()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("-------------------------\n")
	t.Logf("%v", res)
}
