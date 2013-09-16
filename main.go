package gist6418462

import (
	. "gist.github.com/5286084.git"
	. "gist.github.com/5639599.git"
	. "gist.github.com/6433744.git"
	. "gist.github.com/6445065.git"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"runtime"

	"fmt"
)

// Returns the source of the func f.
func GetSourceAsString(f interface{}) string {
	pc := reflect.ValueOf(f).Pointer()
	file, line := runtime.FuncForPC(pc).FileLine(pc)

	var startIndex, endIndex int
	{
		b, err := ioutil.ReadFile(file)
		CheckError(err)
		startIndex, endIndex = GetLineStartEndIndicies(b, line-1)
	}

	fs := token.NewFileSet()
	fileAst, err := parser.ParseFile(fs, file, nil, 0*parser.ParseComments)
	CheckError(err)

	query := func(i interface{}) bool {
		// TODO: Factor-out the unusual overlap check
		if f, ok := i.(*ast.FuncLit); ok && ((startIndex <= int(f.Pos())-1 && int(f.Pos())-1 <= endIndex) || (int(f.Pos())-1 <= startIndex && startIndex <= int(f.End())-1)) {
			return true
		} else if f, ok := i.(*ast.FuncDecl); ok && ((startIndex <= int(f.Pos())-1 && int(f.Pos())-1 <= endIndex) || (int(f.Pos())-1 <= startIndex && startIndex <= int(f.End())-1)) {
			return true
		}
		return false
	}
	funcAst := FindFirst(fileAst, query)

	if funcAst == nil {
		return fmt.Sprintf("<func src not found at %v:%v>", file, line)
	}

	return SprintAst(fs, funcAst)
}

var f2 = func() { panic(1337) }

func main() {
	f := func() {
		println("Hello from anon func!") // Comments are currently not preserved
	}
	if 5*5 > 26 {
		f = f2
	}

	println(GetSourceAsString(f))

	// Output:
	// func() {
	// 	println("Hello from anon func!")
	// }

	f2 := func(a int, b int) int {
		c := a + b
		return c
	}

	println(GetSourceAsString(f2))

	// Output:
	// func(a int, b int) int {
	// 	c := a + b
	// 	return c
	// }

	println(GetSourceAsString(GetSourceAsString))
}
