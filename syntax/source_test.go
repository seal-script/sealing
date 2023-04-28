package syntax

import (
	"bytes"
	"io"
	"testing"
)

func TestSource(t *testing.T) {
	// Slice reader
	data := []byte(`(let x 1)`)
	var in io.Reader = bytes.NewReader(data)
	// Build source
	s := source{}
	errh := func(r, c uint, msg string) {
		t.Logf("[error] Source error at (%d, %d): %s\n", r, c, msg)
	}
	s.init(in, errh)

	s.start()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	s.nextch()
	// t.Logf("%#v\n", s)
	seg := s.segment()
	t.Log(string(seg))
}
