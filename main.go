package gist6418290

import (
	. "gist.github.com/5258650.git"
	//. "gist.github.com/5259939.git"
	. "gist.github.com/5286084.git"
	//. "gist.github.com/5504644.git"
	. "gist.github.com/5707298.git"
	. "gist.github.com/6445065.git"
	//"github.com/shurcooL/go-goon"
	"go/ast"
	//"reflect"
	//"runtime"
	"runtime/debug"
	"strings"
)

/*var _ = GetThisGoSourceDir
var _ = BuildPackageFromImportPath
var _ = CheckError
var _ = debug.FreeOSMemory
var _ = reflect.Copy
var _ = goon.Dump
var _ = ParseDecl
var _ = runtime.BlockProfile*/

// Gets the name of the variable.
func GetVarName(interface{}) string {
	// TODO: Replace use of debug.Stack() with direct use of runtime package...
	str := GetLine(string(debug.Stack()), 3)
	str = str[strings.Index(str, ": ")+len(": "):]
	p, err := ParseStmt(str)
	CheckError(err)

	innerQuery := func(i interface{}) bool {
		if ident, ok := i.(*ast.Ident); ok && ident.Name == "GetVarName" {
			return true
		}
		return false
	}

	query := func(i interface{}) bool {
		if c, ok := i.(*ast.CallExpr); ok && nil != FindFirst(c.Fun, innerQuery) {
			return true
		}
		return false
	}
	callExpr, _ := FindFirst(p, query).(*ast.CallExpr)

	if callExpr == nil {
		return "<var name not found>"
	}

	return callExpr.Args[0].(*ast.Ident).Name
}

func main() {
	var thisIsAFunkyVarName int
	println("Name of var:", GetVarName(thisIsAFunkyVarName))
	var name string = GetVarName(thisIsAFunkyVarName)
	println("Name of var:", name)
}
