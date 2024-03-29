package syntax

import (
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/seal-script/sealing/ast"
	"github.com/seal-script/sealing/utils"
)

// A Location is a triple: (<File>, <line>, <col>)
type Location = utils.Location
type TokenStream = scanner

// An interface for parsing
type Parsing interface {
	ParseDecl() (ast.Decl, error)
	ParseExpr() (ast.Expr, error)
}

// Error type for parsing
type ParsingError struct {
	error
	Location Location
}

// Error wrapper
func errorOf(location Location, format string, args ...any) ParsingError {
	return ParsingError{
		error:    fmt.Errorf(format, args...),
		Location: location,
	}
}

func (pErr *ParsingError) String() string {
	return fmt.Sprintf(
		`
ParsingError {
	%v,
	location: %v,
}
`, pErr.error, pErr.Location,
	)
}

// better way to debug?
const debug = false
const trace = false

// The main parser
type Parser struct {
	scanner
	filePath string
}

func NewParser(t *testing.T, in io.Reader) Parser {
	p := Parser{}
	p.Init(in, func(err error) {
		t.Log(err.Error())
	})
	p.next() // Fill buffer
	return p
}

func (p *Parser) Init(r io.Reader, errHandler func(error)) {
	p.scanner.init(r, func(r, c uint, msg string) {
		errHandler(fmt.Errorf("Syntax error: (%d, %d) %s", r, c, msg))
	}, 0)
}

func (p *Parser) errorOf(format string, args ...any) ParsingError {
	return errorOf(p.Locate(), format, args...)
}

// Parsing a file
func (p *Parser) ParseFile() (*ast.File, error) {
	f := new(ast.File)
	f.Location = p.Locate()

	// While not end of file
	for p.token.tag != _EOF {
		// fmt.Println(p.token)
		decl, err := p.ParseDecl()
		if err != nil {
			return nil, err
		}
		f.DeclList = append(f.DeclList, decl)
		p.next() // seperator?
	}
	return f, nil
}

func (p *Parser) ParseDecl() (ast.Decl, error) {
	if p.token.tag == _Ident {
		fName, err := p.ParseNameExpr()
		if err != nil {
			return nil, err
		}
		switch p.token.tag {
		case _Colon:
			return p.ParseTypeDecl(fName)
		case _Ident, _Assign, _ParentLeft, _Integer:
			return p.ParseFuncDecl(fName)
		default:
			return nil, p.errorOf(
				"Expected ':' | '=' | identifier, found %#v\n",
				p.token,
			)
		}
	}
	return nil, p.errorOf(
		"Expected identifier, found %#v\n",
		p.token,
	)
}

// `let x <expression>)`
// `let (f x...) <expression>)`
// f | x = <expression>
func (p *Parser) ParseFuncDecl(fName *ast.Name) (*ast.FuncDecl, error) {
	decl := new(ast.FuncDecl)
	decl.Name = fName
	// decl.Location =
	// p.next()
	switch p.token.tag {

	// Function w/o parameters
	case _Assign, _Ident, _ParentLeft, _Integer:
		args := []ast.Pattern{}
		var err error
		for err == nil {
			var arg ast.Pattern
			arg, err = p.ParsePatternExpr()
			if err == nil {
				args = append(args, arg)
			}
		}
		decl.Params = args
		if p.token.tag == _Assign {
			p.next()
		}
		body, err := p.ParseFuncCallExpr()
		if err != nil {
			return nil, err
		}
		decl.Body = body

	default:
		return nil, p.errorOf("Error while parsing function declaration")
	}
	return decl, nil
}

// x : Int
// f : Int -> Int
func (p *Parser) ParseTypeDecl(fName *ast.Name) (*ast.TypeDecl, error) {
	p.next()
	t, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	return &ast.TypeDecl{
		Name: fName,
		Type: t,
	}, nil
}

// Int
// Int -> Int
// (Int -> Int) -> Int
func (p *Parser) ParseType() (ast.Type, error) {
	switch p.token.tag {
	case _Ident:
		t, err := p.ParseFuncCallExpr()
		if err != nil {
			return nil, err
		}
		if p.token.tag != _Arrow {
			return t, nil
		}
		p.next()
		ts, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		return &ast.FuncType{
			Context: []ast.Field{},
			Types:   []ast.Type{t, ts},
		}, nil

	case _ParentLeft:
		p.next()
		t, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		if p.token.tag == _ParentRight {
			p.next()
		}
		if p.token.tag != _Arrow {
			return t, nil
		}
		p.next()
		ts, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		return &ast.FuncType{
			Context: []ast.Field{},
			Types:   []ast.Type{t, ts},
		}, nil
	default:
		return nil, p.errorOf("ParseType: Unexpected token: %#v\n", p.token)
	}
}

// `(f x...)`
func (p *Parser) ParseExpr() (ast.Expr, error) {
	switch p.token.tag {
	case _Integer:
		return p.ParseIntegerExpr()
	case _Ident:
		fCall := new(ast.CallExpr)
		fName := &ast.Name{Value: p.token.lit}
		fCall.Fun = fName
		p.next()
		return fCall, nil
	case _ParentLeft:
		p.next()
		if p.token.tag == _Symbol {
			fCall := new(ast.CallExpr)
			name := &ast.Name{Value: p.token.lit}
			fCall.Fun = name
			p.next()
			if p.token.tag == _ParentRight {
				p.next()
			}

			for param, err := p.ParseExpr(); err == nil; {
				fCall.ArgList = append(fCall.ArgList, param)
				param, err = p.ParseExpr()
			}
			return fCall, nil
		}

		// Function call
		expr, err := p.ParseFuncCallExpr()
		if err != nil {
			return nil, err
		}
		if p.token.tag == _ParentRight {
			p.next()
		}
		return expr, nil
	default:
		return nil, p.errorOf("ParseExpr error: encounter %#v", p.token)
	}
}

// `x`
func (p *Parser) ParseNameExpr() (*ast.Name, error) {
	switch p.token.tag {
	case _Ident:
		name := &ast.Name{Value: p.token.lit}
		p.next()
		return name, nil
	default:
		return nil, p.errorOf("ParseNameExpr error: encounter %#v", p.token)
	}
}

func (p *Parser) ParsePatternExpr() (ast.Pattern, error) {
	switch p.token.tag {
	case _Ident:
		return p.ParseNameExpr()
	case _ParentLeft:
		p.next()
		pattern, err := p.ParseFuncCallExpr()
		if err != nil {
			return nil, err
		}
		if p.token.tag == _ParentRight {
			p.next()
		}
		return pattern, nil
	case _Integer:
		return p.ParseIntegerExpr()

	// case _Float:
	// return p.Parse
	default:
		return nil, p.errorOf("ParsePatternExpr: %#v\n", p.token)
	}
}

// `7`
func (p *Parser) ParseIntegerExpr() (*ast.Integer, error) {
	switch p.token.tag {
	case _Integer:
		i, err := strconv.Atoi(p.token.lit)
		if err != nil {
			return nil, err
		}
		j := &ast.Integer{Value: i}
		p.next()
		return j, nil
	default:
		return nil, p.errorOf("ParseIntegerExpr error: encounter %#v", p.token)
	}
}

// `f x...`
func (p *Parser) ParseFuncCallExpr() (*ast.CallExpr, error) {
	if !(p.token.tag == _Ident || p.token.tag == _ParentLeft) {
		return nil, p.errorOf("ParseFuncCallExpr error: encounter %#v", p.token)
	}
	fCall := new(ast.CallExpr)
	switch p.token.tag {
	// f x
	case _Ident:
		fName := &ast.Name{Value: p.token.lit}
		fCall.Fun = fName
		p.next()
		// fmt.Println(p.token)
		for param, err := p.ParseExpr(); err == nil; {
			fCall.ArgList = append(fCall.ArgList, param)
			param, err = p.ParseExpr()
		}
		// if p.token.tag != _ParentRight {
		// 	return nil, fmt.Errorf("ParseFuncCallExpr error: encounter %#v", p.token)
		// }
		// p.next()
		return fCall, nil

	// case _Symbol:

	// (f x) y
	case _ParentLeft:
		p.next()

		// fmt.Println(">>>", p.token)
		if p.token.tag == _Symbol {
			// (-) 3 2
			name := &ast.Name{Value: p.token.lit}
			fCall.Fun = name
			p.next()
		} else {
			f, err := p.ParseExpr()
			// Normal function call
			if err != nil {
				return nil, fmt.Errorf("Error of parser: %v\n", err)
			}
			fCall.Fun = f
		}

		if p.token.tag == _ParentRight {
			p.next()
		}
		for param, err := p.ParseExpr(); err == nil; {
			fCall.ArgList = append(fCall.ArgList, param)
			param, err = p.ParseExpr()
		}
		return fCall, nil
	default:
		return nil, p.errorOf("ParseFuncCallExpr error: encounter %#v", p.token)
	}
}

func (p *Parser) Locate() Location {
	return Location{
		FilePath: p.filePath,
		Line:     p.line,
		Col:      p.col,
	}
}

// map : (a b : Type) => (a -> b) -> [a] -> [b]
