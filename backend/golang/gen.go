package golang

import (
	"fmt"

	"github.com/seal-script/sealing/ast"
)

type Generator[Ir any] interface {
	Gen(ast.File) (Ir, error)
}

type GenString struct {
	TEnv map[*ast.Name]ast.Type
	FEnv map[*ast.Name]*ast.FuncDecl
}

func (g *GenString) Gen(file *ast.File) (string, error) {
	decls := file.DeclList
	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.TypeDecl:
			g.TEnv[d.Name] = d.Type
		case *ast.FuncDecl:
			d.Type = g.TEnv[d.Name]
			g.FEnv[d.Name] = d
		default:
			return "", fmt.Errorf("Error of generator: Gen %#v", decl)
		}
	}
	ans := ""
	for _, fDecl := range g.FEnv {
		f, err := GenFunc(fDecl)
		if err != nil {
			return "", err
		}
		ans = ans + "\n\n" + f
	}
	return ans, nil
}

func GenFunc(fDecl *ast.FuncDecl) (string, error) {
	ans := fmt.Sprintf(`func %s`, fDecl.Name.Value)
	// utils.Todo()
	body, err := GenExpr(fDecl.Body)
	if err != nil {
		return "", err
	}
	ans += body
	return ans, nil
}

func GenExpr(expr ast.Expr) (string, error) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return GenFuncCall(e)
	case *ast.Name:
		return e.Value, nil
	default:
		return "", fmt.Errorf("Error of generator: GenExpr: Unknown expr: %#v", e)
	}
}

func GenFuncCall(funcCall *ast.CallExpr) (string, error) {
	args := ""
	for _, arg := range funcCall.ArgList {
		switch param := arg.(type) {
		case *ast.CallExpr:
			if args == "" {
				args = fmt.Sprintf("%v", param)
			} else {
				args = args + ", " + fmt.Sprintf("%v", param)
			}
		case *ast.Name:
			if args == "" {
				args = param.Value
			} else {
				args = args + ", " + param.Value
			}
		case *ast.Integer:
			args += ", " + fmt.Sprintf("%v", param)
		default:
			return "", fmt.Errorf(
				"Error of generator: GenFuncCall: Unimplemented pattern matching: %v",
				arg,
			)
		}
	}
	f, err := GenExpr(funcCall.Fun)
	if err != nil {
		return "", err
	}
	ans := fmt.Sprintf("%s(%s)", f, args)
	return ans, nil
}
