package main

import (
	//"fmt"
	. "gist.github.com/5258650.git"
	. "gist.github.com/5259939.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
	. "gist.github.com/5707298.git"
	. "gist.github.com/6433744.git"
	. "gist.github.com/6445065.git"
	"github.com/shurcooL/go-goon"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"runtime"
	"runtime/debug"
)

var _ = GetThisGoSourceDir
var _ = BuildPackageFromImportPath
var _ = CheckError
var _ = debug.FreeOSMemory
var _ = reflect.Copy
var _ = goon.Dump
var _ = ParseDecl
var _ = runtime.BlockProfile
var _ = GetLine
var _ = FindFirst

// Returns source of anon func string.
// TODO: Finish...
func GetSourceAsString(f interface{}) string {
	pc := reflect.ValueOf(f).Pointer()
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	//println("runtime.FuncForPC(pc).FileLine(pc):", file, line)

	var startIndex, endIndex int
	{
		b, err := ioutil.ReadFile(file)
		CheckError(err)
		startIndex, endIndex = GetLineStartEndIndicies(b, line-1)
	}

	fs := token.NewFileSet()
	fileAst, err := parser.ParseFile(fs, file, nil, 0*parser.ParseComments)
	CheckError(err)
	if 0 == 1 {
		goon.Dump(fileAst)
	}

	//return fmt.Sprintf("%s: %d -> [%d, %d]", file, line, startIndex, endIndex)

	query := func(i interface{}) bool {
		if f, ok := i.(*ast.FuncLit); ok && startIndex <= int(f.Type.Func)-1 && int(f.Type.Func)-1 <= endIndex {
			return true
		}
		return false
	}
	flit := FindFirst(fileAst, query)

	return SprintAst(fs, flit)
}

var f2 = func() { panic(1337) }

func main() {
	f := func() {
		println("Hello from anon func!")
	}
	if 5*5 > 26 {
		f = f2
	}

	print(GetSourceAsString(f))
}

/* Need to find v of type (*ast.FuncLit) with startIndex <= v.(*ast.FuncLit).Type.Func-1 <= endIndex

Rhs: ([]ast.Expr)([]ast.Expr{
	(*ast.FuncLit)(&ast.FuncLit{
		Type: (*ast.FuncType)(&ast.FuncType{
			Func: (token.Pos)(1112),
			Params: (*ast.FieldList)(&ast.FieldList{
				Opening: (token.Pos)(1116),
				List:    ([]*ast.Field)([]*ast.Field{}),
				Closing: (token.Pos)(1117),
			}),
			Results: (*ast.FieldList)(nil),
		}),
		Body: (*ast.BlockStmt)(&ast.BlockStmt{
			Lbrace: (token.Pos)(1119),
			List: ([]ast.Stmt)([]ast.Stmt{
				(*ast.ExprStmt)(&ast.ExprStmt{
					X: (*ast.CallExpr)(&ast.CallExpr{
						Fun: (*ast.Ident)(&ast.Ident{
							NamePos: (token.Pos)(1123),
							Name:    (string)("println"),
							Obj:     (*ast.Object)(nil),
						}),
						Lparen: (token.Pos)(1130),
						Args: ([]ast.Expr)([]ast.Expr{
							(*ast.BasicLit)(&ast.BasicLit{
								ValuePos: (token.Pos)(1131),
								Kind:     (token.Token)(9),
								Value:    (string)("\"Hello from anon func!\""),
							}),
						}),
						Ellipsis: (token.Pos)(0),
						Rparen:   (token.Pos)(1154),
					}),
				}),
			}),
			Rbrace: (token.Pos)(1157),
		}),
	}),
}),*/
