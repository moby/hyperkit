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
func GetExprAsString(interface{}) string {
	return GetParentArgExprAsString(0)
}

// Gets the argIndex argument expression of parent func call as a string.
func GetParentArgExprAsString(argIndex uint32) string {
	// TODO: Replace use of debug.Stack() with direct use of runtime package...
	parentName := GetLine(string(debug.Stack()), 3)
	parentName = parentName[1:strings.Index(parentName, ": ")]

	str := GetLine(string(debug.Stack()), 5)
	str = str[strings.Index(str, ": ")+len(": "):]
	p, err := ParseStmt(str)
	CheckError(err)

	innerQuery := func(i interface{}) bool {
		if ident, ok := i.(*ast.Ident); ok && ident.Name == parentName {
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

	return SprintAstBare(callExpr.Args[argIndex])
}

func getMySecondArgExprAsString(int, int) string {
	return GetParentArgExprAsString(1)
}

func main() {
	var thisIsAFunkyVarName int
	println("Name of var:", GetExprAsString(thisIsAFunkyVarName))
	var name string = GetExprAsString(thisIsAFunkyVarName)
	println("Name of var:", name)
	println("Some func name:", GetExprAsString(strings.HasPrefix))
	println("Name of second arg:", getMySecondArgExprAsString(5, thisIsAFunkyVarName))
}
