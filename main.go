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
		// TODO: Abstract-out the overlap check
		if f, ok := i.(*ast.FuncLit); ok && ((startIndex <= int(f.Pos())-1 && int(f.Pos())-1 <= endIndex) || (int(f.Pos())-1 <= startIndex && endIndex <= int(f.End())-1)) {
			return true
		}
		return false
	}
	funcLit := FindFirst(fileAst, query)

	return SprintAst(fs, funcLit)
}

var f2 = func() { panic(1337) }

func main() {
	f := func() {
		println("Hello from anon func!") // Comments are currently not preserved
	}
	if 5*5 > 26 {
		f = f2
	}

	print(GetSourceAsString(f))

	// Output:
	// func() {
	// 	println("Hello from anon func!")
	// }
}
