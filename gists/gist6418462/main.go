package gist6418462

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"runtime"

	. "github.com/shurcooL/go/gists/gist5286084"
	. "github.com/shurcooL/go/gists/gist5639599"
	. "github.com/shurcooL/go/gists/gist6433744"
	. "github.com/shurcooL/go/gists/gist6445065"
)

// GetSourceAsString returns the source of the func f.
func GetSourceAsString(f interface{}) string {
	v := reflect.ValueOf(f)
	if v.IsNil() {
		return "nil"
	}
	pc := v.Pointer()
	function := runtime.FuncForPC(pc)
	if function == nil {
		return "nil"
	}
	file, line := function.FileLine(pc)

	var startIndex, endIndex int
	{
		b, err := ioutil.ReadFile(file)
		CheckError(err)
		startIndex, endIndex = GetLineStartEndIndicies(b, line-1)
	}

	fs := token.NewFileSet()
	fileAst, err := parser.ParseFile(fs, file, nil, 0*parser.ParseComments)
	CheckError(err)

	// TODO: Consider using ast.Walk() instead of custom FindFirst()
	query := func(i interface{}) bool {
		// TODO: Factor-out the unusual overlap check
		if f, ok := i.(*ast.FuncLit); ok && ((startIndex <= int(f.Pos())-1 && int(f.Pos())-1 <= endIndex) || (int(f.Pos())-1 <= startIndex && startIndex <= int(f.End())-1)) {
			return true
		}
		return false
	}
	funcAst := FindFirst(fileAst, query)

	// If func literal wasn't found, try again looking for func declaration
	if funcAst == nil {
		query := func(i interface{}) bool {
			// TODO: Factor-out the unusual overlap check
			if f, ok := i.(*ast.FuncDecl); ok && ((startIndex <= int(f.Pos())-1 && int(f.Pos())-1 <= endIndex) || (int(f.Pos())-1 <= startIndex && startIndex <= int(f.End())-1)) {
				return true
			}
			return false
		}
		funcAst = FindFirst(fileAst, query)
	}

	if funcAst == nil {
		return fmt.Sprintf("<func src not found at %v:%v>", file, line)
	}

	return SprintAst(fs, funcAst)
}
