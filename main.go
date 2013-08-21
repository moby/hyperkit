package gist5707298

import (
	"go/ast"
	"go/parser"
	"go/token"

	//"github.com/shurcooL/go-goon"
)

// ParseStmt is a convenience function for obtaining the AST of a statement x.
// The position information recorded in the AST is undefined.
//
func ParseStmt(x string) (ast.Stmt, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p;func _(){\n//line :1\n"+x+"\n;}", 0)
	if err != nil {
		return nil, err
	}
	return file.Decls[0].(*ast.FuncDecl).Body.List[0], nil
}

func ParseDecl(x string) (ast.Decl, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p\n//line :1\n"+x+"\n", 0)
	if err != nil {
		return nil, err
	}
	return file.Decls[0], nil
}

func main() {
	//goon.Dump(ParseStmt("var x int"))
}
