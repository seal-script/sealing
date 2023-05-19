package golang

import (
	"fmt"

	"github.com/seal-script/sealing/ast"
	"github.com/seal-script/sealing/utils"
)

type Generator[Ir any] interface {
	Gen(ast.File) (Ir, error)
}

type GenString struct {
	TEnv map[string]ast.Type
	FEnv map[string]*ast.FuncDecl
}

func (g *GenString) Gen(file *ast.File) (string, error) {
	decls := file.DeclList
	for i, decl := range decls {
		switch d := decl.(type) {
		case *ast.TypeDecl:
			g.TEnv[d.Name.Value] = d.Type
		case *ast.FuncDecl:
			decls[i].(*ast.FuncDecl).Type = g.TEnv[d.Name.Value]
			g.FEnv[d.Name.Value] = d
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
	utils.Todo()
	ts := fDecl.Type.(*ast.FuncType).Types
	ps := fDecl.Params
	params := ""
	for i, p := range ps {
		e, err := GenExpr(p.(ast.Expr))
		if err != nil {
			return "", err
		}
		t, err := GenType(ts[i])
		if err != nil {
			return "", nil
		}
		pair := fmt.Sprintf("%s %s", e, t)
		if params == "" {
			params = pair
		} else {
			params += ", " + pair
		}
	}
	ans += fmt.Sprintf("(%s)", params)
	ans += " {\n"
	body, err := GenExpr(fDecl.Body)
	if err != nil {
		return "", err
	}
	ans += "return " + body
	ans += "\n}"
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

func GenType(tpe ast.Type) (string, error) {
	switch t := tpe.(type) {
	case *ast.FuncType:
		// fmt.Println(t.Types)
		rest, err := GenType(&ast.FuncType{
			Context: t.Context,
			Types:   t.Types[1:],
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("func(%s) %s", t.Types[0], rest), nil
	case *ast.CallExpr:
		ans := fmt.Sprintf("%v", t.Fun)
		param := ""
		for _, arg := range t.ArgList {
			x, err := GenType(arg)
			if err != nil {
				return "", err
			}
			if param == "" {
				param = x
			} else {
				param += ", " + x
			}
		}
		if param != "" {
			ans += fmt.Sprintf("[%s]", param)
		}
		return ans, nil
	default:
		return "", fmt.Errorf("Error of generator: Unknown type: %#v\n", t)
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
