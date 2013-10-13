package gist6418290

import (
	. "gist.github.com/5258650.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5639599.git"
	. "gist.github.com/5707298.git"
	. "gist.github.com/6445065.git"
	"go/ast"
	"runtime/debug"
	"strings"
)

// Gets the expression as a string.
func ExprToString(interface{}) string {
	// TODO: Replace use of debug.Stack() with direct use of runtime package...
	str := GetLine(string(debug.Stack()), 3)
	str = str[strings.Index(str, ": ")+len(": "):]
	p, err := ParseStmt(str)
	CheckError(err)

	innerQuery := func(i interface{}) bool {
		if ident, ok := i.(*ast.Ident); ok && ident.Name == "ExprToString" {
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
		return "<expr not found>"
	}

	return SprintAstBare(callExpr.Args[0])
}

func main() {
	var thisIsAFunkyVarName int
	println("Name of var:", ExprToString(thisIsAFunkyVarName))
	var name string = ExprToString(thisIsAFunkyVarName)
	println("Name of var:", name)
	println("Some func name:", ExprToString(strings.HasPrefix))
}
