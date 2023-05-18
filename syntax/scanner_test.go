package syntax

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
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
	data := []byte(`123`)
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

func TestSamples(t *testing.T) {
	data := []byte{}
	samples := samples()
	for _, sample := range samples {
		data = append(data, []byte(sample.raw)...)
		data = append(data, []byte(" ")...)
	}
	var in io.Reader = bytes.NewReader(data)

	s := newScanner(t, in)
	for _, sample := range samples {
		err := s.next()
		if err != nil {
			t.Errorf("Error of scanner: %#v\n", err)
			break
		}
		token := s.token
		// t.Logf("%v, %v\n", &token, &sample.token)
		if !(token == sample.token) {
			t.Errorf("Error of scanner: expected %#v, found %#v\n", sample.token, token)
			return
		}
	}
	s.next()
	t.Logf("End with %#v\n", s.token)
}

// const (
//  	_ParentLeft  tokenTag = iota // Left '('
//  	_ParentRight                 // Right ')'
//  	_BraceLeft                   // Left '{'
//  	_BraceRight                  // Right '}'
//  	_Ident                       // Identifier
//  	_Integer                     // Integer lit
//  	_Float                       // Float number lit
//  	_Complex                     // Complex number lit
//  	_Rune                        // A single unicode character
//  	_String                      // String with some encoding
//  	_Let                         // Let binding
//  	_Type                        // Type declaration
//  	_Quote                       // ' single quote
//  	_DoubleQuote                 // Double quote "
//  	_Comment                     // Comment
//  	_Semi                        // ';' or '\n'
//  	_EOF                         // End Of File
// )

type sample struct {
	raw   string
	token Token
}

var signSamples = [...]sample{
	{"(", Token{_ParentLeft, "("}},
	{")", Token{_ParentRight, ")"}},
	{"{", Token{_BraceLeft, "{"}},
	{"}", Token{_BraceRight, "}"}},
	{")", Token{_ParentRight, ")"}},
	{"'", Token{_Quote, "'"}},
	{"\"", Token{_DoubleQuote, "\""}},
	{";", Token{_Semi, ";"}},
	{"=", Token{_Assign, "="}},
	{"->", Token{_Arrow, "->"}},
}

var keywordSamples = [...]sample{
	{"let", Token{_Let, "let"}},
}

func identifierSamples() []sample {
	ans := []sample{
		{"a", Token{_Ident, "a"}},
		{"a1", Token{_Ident, "a1"}},
		{"test", Token{_Ident, "test"}},
		{"test09", Token{_Ident, "test09"}},
		{"test_me", Token{_Ident, "test_me"}},
		{"变量", Token{_Ident, "变量"}},
		{"命运石之门", Token{_Ident, "命运石之门"}},
	}
	return ans
}

func symbolSamples() []sample {
	ans := []sample{
		{"<>", Token{_Symbol, "<>"}},
		{"-", Token{_Symbol, "-"}},
		{"++", Token{_Symbol, "++"}},
		{"@>", Token{_Symbol, "@>"}},
		{"==", Token{_Symbol, "=="}},
		{"!=", Token{_Symbol, "!="}},
		{"::", Token{_Symbol, "::"}},
	}
	return ans
}

func integerSamples(limit int) []sample {
	ans := []sample{}

	str := func(x any) string {
		return fmt.Sprintf("%v", x)
	}
	for i := 0; i < limit; i++ {
		ans = append(ans, sample{str(i), Token{_Integer, str(i)}})
	}
	return ans
}

func floatSamples(limit int) []sample {
	ans := []sample{}
	for i := 0; i < limit; i++ {
		num := rand.Float64()
		s := fmt.Sprintf("%v", num)
		ans = append(ans, sample{s, Token{_Float, s}})
	}
	return ans
}

func samples() []sample {
	ans := []sample{}
	ans = append(ans, signSamples[:]...)
	ans = append(ans, keywordSamples[:]...)
	ans = append(ans, identifierSamples()[:]...)
	ans = append(ans, symbolSamples()[:]...)
	ans = append(ans, integerSamples(100)[:]...)
	ans = append(ans, floatSamples(100)[:]...)
	return ans
}
