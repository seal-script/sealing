package syntax

import (
	"fmt"
	"io"

	"github.com/seal-script/sealing/ast"
	"github.com/seal-script/sealing/utils"
)

type Location = utils.Location

const debug = false
const trace = false

type Parser struct {
	scanner
	location Location
}

func (p *Parser) Init(r io.Reader, errHandler func(error)) {
	p.scanner.init(r, func(r, c uint, msg string) {
		errHandler(fmt.Errorf("Syntax error: (%d, %d) %s", r, c, msg))
	}, 0)
}

// Parsing a file
func (p *Parser) ParseFile() (*ast.File, error) {
	f := new(ast.File)
	f.Location = p.Locate(p.line, p.col)

	// While not end of file
	for p.token.tag != _EOF {
		switch p.token.tag {
		case _ParentLeft:
			decl, err := p.ParseDecl()
			if err != nil {
				return nil, err
			}
			f.DeclList = append(f.DeclList, decl)
		default:
			p.next()
			continue
		}
	}
	return f, nil
}

func (p *Parser) ParseDecl() (ast.Decl, error) {
	p.next()

	// Let binding
	if p.token.tag == _Let {
		return p.ParseFuncDecl()
	}

	// Type declaration
	if p.token.tag == _Type {
		return p.ParseTypeDecl()
	}
	return nil, nil
}

func (p *Parser) ParseFuncDecl() (*ast.FuncDecl, error) {
	decl := new(ast.FuncDecl)
	p.next()
	if p.token.tag == _ParentLeft {

	}
	return decl, nil
}

func (p *Parser) ParseList() ([]Token, error) {
	return nil, nil
}

func (p *Parser) ParseTypeDecl() (*ast.TypeDecl, error) {
	return nil, nil
}

func (p *Parser) Locate(line, col uint) Location {
	return Location{
		FilePath: p.location.FilePath,
		Line:     line,
		Col:      col,
	}
}
