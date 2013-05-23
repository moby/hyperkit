package main

import (
	"go/token"
	"go/printer"
	"bytes"
	"fmt"
)

func SprintAst(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return buf.String()
}

func SprintAstBare(node interface{}) string {
	fset := token.NewFileSet()
	return SprintAst(fset, node)
}

func PrintlnAst(fset *token.FileSet, node interface{}) {
	fmt.Println(SprintAst(fset, node))
}

func PrintlnAstBare(node interface{}) {
	fset := token.NewFileSet()
	PrintlnAst(fset, node)
}

func main() {
}