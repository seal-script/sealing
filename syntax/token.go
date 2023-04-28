package syntax

import "fmt"

type tokenTag int

const (
	_ParentLeft  tokenTag = iota // Left '('
	_ParentRight                 // Right ')'
	_Ident                       // Identifier
	_Integer                     // Integer lit
	_Float                       // Float number lit
	_Complex                     // Complex number lit
	_Rune                        // A single unicode character
	_String                      // String with some encoding
	_Let                         // Let binding
	_Type                        // Type declaration
	_Quote                       // ' single quote
	_DoubleQuote                 // Double quote "
	_Comment                     // Comment
	_Semi                        // ';' or '\n'
	_EOF                         // End Of File
)

func (tag tokenTag) String() string {
	switch tag {
	case _ParentLeft:
		return "ParentLeft"

	case _ParentRight:
		return "ParentRight"

	case _Ident:
		return "Ident"

	case _Integer:
		return "Integer"

	case _Float:
		return "Float"

	case _Complex:
		return "Complex"

	case _Let:
		return "Let"

	case _Type:
		return "Type"

	case _Comment:
		return "Comment"

	case _Semi:
		return "Semicolon"

	case _Quote:
		return "Quote"

	case _DoubleQuote:
		return "DoubleQuote"

	case _EOF:
		return "EOF"

	case _Rune:
		return "Rune"

	case _String:
		return "String"

	default:
		return "Unknown"
	}
}

type Token struct {
	tag tokenTag
	lit string
}

// func (tok Token) String() string {
// 	return fmt.Sprintf("{%s, \"%s\"}", tok.tag.String(), tok.lit)
// }

func (tok *Token) String() string {
	return fmt.Sprintf("{%s, \"%s\"}", tok.tag.String(), tok.lit)
}
