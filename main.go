package main

import (
	"fmt"
	. "gist.github.com/5258650.git"
	. "gist.github.com/5259939.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5707298.git"
	. "gist.github.com/6433744.git"
	"github.com/shurcooL/go-goon"
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

// Returns source of anon func string.
// TODO: Finish...
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
	a, err := parser.ParseFile(fs, file, nil, 0)
	CheckError(err)
	if 0 == 1 {
		goon.Dump(a)
	}

	return fmt.Sprintf("%s: %d -> [%d, %d]", file, line, startIndex, endIndex)
}

func main() {
	f := func() {
		println("Hello from anon func!")
	}

	println(GetSourceAsString(f))
}
