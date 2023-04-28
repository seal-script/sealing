package syntax

import (
	"bytes"
	"io"
	"testing"
)

func newScanner(t *testing.T, in io.Reader) scanner {
	s := scanner{}
	errh := func(r, c uint, msg string) {
		t.Logf("[error] Source error at (%d, %d): %s\n", r, c, msg)
	}
	s.init(in, errh, 0)
	return s
}

func TestIntegers(t *testing.T) {
	return
}

func TestScanner(t *testing.T) {
	// Slice reader
	data := []byte(`(let x 测试gdfh 烤红薯烤豆腐); gg 123.4 (let y 666)`)
	var in io.Reader = bytes.NewReader(data)

	// Build scanner
	s := newScanner(t, in)

	// Scan!
	for s.token.tag != _EOF {
		err := s.next()
		if err != nil {
			break
		}
		t.Logf("%v\n", &s.token)
	}
}
