package syntax

import (
	"bytes"
	"io"
	"testing"
)

func NewParser(t *testing.T, in io.Reader) Parser {
	p := Parser{}
	p.Init(in, func(err error) {
		t.Log(err.Error())
	})
	return p
}

func TestParser(t *testing.T) {
	// Slice reader
	data := []byte(`(let x 测试gdfh 烤红薯烤豆腐)`)
	var in io.Reader = bytes.NewReader(data)

	p := NewParser(t, in)
	ast, err := p.ParseFile()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", ast)
}
