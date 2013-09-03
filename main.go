package main

import (
	. "gist.github.com/5258650.git"
	. "gist.github.com/5259939.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5707298.git"
	"gist.github.com/6418290.git"
	"github.com/shurcooL/go-goon"
	"go/parser"
	"go/token"
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
	{
		println(string(debug.Stack()))
		buf := make([]byte, 1024)
		runtime.Stack(buf, true)
		println(string(buf))
	}

	{
		val := reflect.ValueOf(f)
		println(val.Type().Kind().String())
	}

	return "???"
}

func main() {
	f := func() { println("Hello from anon func!") }
	f()

	{
		doc := GetDocPackageAll(BuildPackageFromSrcDir(GetThisGoSourceDir()))

		println(doc.Filenames[0])
		//goon.Dump(doc)
	}

	{
		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, GetThisGoSourceFilepath(), nil, 0)
		CheckError(err)

		println(f.Name.String())
	}

	println(GetSourceAsString(f))
	var thisIsAFunkyVarName int
	//var name string = (GetVarName)    (thisIsAFunkyVarName)
	println("Name of var:", (gist6418290.GetVarName)(thisIsAFunkyVarName))
}
