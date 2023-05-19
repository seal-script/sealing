package golang

import (
	"bytes"
	"io"
	"testing"

	"github.com/seal-script/sealing/ast"
	"github.com/seal-script/sealing/syntax"
)

func TestGen(t *testing.T) {
	data := []byte(`
		test : Int -> Int;
		test x = add x 1
	`)
	var in io.Reader = bytes.NewReader(data)

	p := syntax.NewParser(t, in)
	file, err := p.ParseFile()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("-------------------------\n")
	t.Logf("%v", file)

	var g GenString = GenString{
		TEnv: map[string]ast.Type{},
		FEnv: map[string]*ast.FuncDecl{},
	}
	s, err := g.Gen(file)
	if err != nil {
		t.Error(err)
	}
	t.Logf(s)
}
