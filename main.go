package main

import (
	"go/token"
	"go/printer"
	"bytes"
	"fmt"
)

func SprintAstBare(node interface{}) string {
	fset := token.NewFileSet()
	return SprintAst(fset, node)
}

func SprintAst(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return buf.String()
}

func PrintlnAst(fset *token.FileSet, node interface{}) {
	fmt.Println(SprintAst(fset, node))
}

func main() {
}