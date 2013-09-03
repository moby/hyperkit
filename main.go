package gist6418290

import (
	. "gist.github.com/5258650.git"
	. "gist.github.com/5259939.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5707298.git"
	"github.com/shurcooL/go-goon"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

var _ = GetThisGoSourceDir
var _ = BuildPackageFromImportPath
var _ = CheckError
var _ = debug.FreeOSMemory
var _ = reflect.Copy
var _ = goon.Dump
var _ = ParseDecl
var _ = runtime.BlockProfile

// Gets the name of the variable.
// TODO: Finish...
//
//(*ast.CallExpr)(&ast.CallExpr{
//	Fun: (*ast.ParenExpr)(&ast.ParenExpr{
//		Lparen: (token.Pos)(55),
//		X: (*ast.Ident)(&ast.Ident{
//			NamePos: (token.Pos)(56),
//			Name:    (string)("GetVarName"),
//			Obj:     (*ast.Object)(nil),
//		}),
//		Rparen: (token.Pos)(66),
//	}),
//	Lparen: (token.Pos)(71),
//	Args: ([]ast.Expr)([]ast.Expr{
//		(*ast.Ident)(&ast.Ident{
//			NamePos: (token.Pos)(72),
//			Name:    (string)("thisIsAFunkyVarName"),
//			Obj:     (*ast.Object)(nil),
//		}),
//	}),
//	Ellipsis: (token.Pos)(0),
//	Rparen:   (token.Pos)(91),
//}),
//
//(*ast.CallExpr)(&ast.CallExpr{
//	Fun: (*ast.ParenExpr)(&ast.ParenExpr{
//		Lparen: (token.Pos)(55),
//		X: (*ast.SelectorExpr)(&ast.SelectorExpr{
//			X: (*ast.Ident)(&ast.Ident{
//				NamePos: (token.Pos)(56),
//				Name:    (string)("gist6418290"),
//				Obj:     (*ast.Object)(nil),
//			}),
//			Sel: (*ast.Ident)(&ast.Ident{
//				NamePos: (token.Pos)(68),
//				Name:    (string)("GetVarName"),
//				Obj:     (*ast.Object)(nil),
//			}),
//		}),
//		Rparen: (token.Pos)(78),
//	}),
//	Lparen: (token.Pos)(83),
//	Args: ([]ast.Expr)([]ast.Expr{
//		(*ast.Ident)(&ast.Ident{
//			NamePos: (token.Pos)(84),
//			Name:    (string)("thisIsAFunkyVarName"),
//			Obj:     (*ast.Object)(nil),
//		}),
//	}),
//	Ellipsis: (token.Pos)(0),
//	Rparen:   (token.Pos)(103),
//}),
//
// Need to find that *ast.CallExpr within AST, return Args[0].(*ast.Ident).Name
func GetVarName(v interface{}) string {
	str := GetLine(string(debug.Stack()), 3)
	str = str[strings.Index(str, ": ")+len(": "):]
	p, err := ParseStmt(str)
	CheckError(err)
	return goon.Sdump(p)
}

func main() {
	var thisIsAFunkyVarName int
	println("Name of var:", GetVarName(thisIsAFunkyVarName))
	var name string = GetVarName(thisIsAFunkyVarName)
	println("Name of var:", name)
}
