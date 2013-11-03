package gist6418290

import (
	"go/ast"
	"runtime/debug"
	"strings"

	. "gist.github.com/5258650.git"
	. "gist.github.com/5639599.git"
	. "gist.github.com/5707298.git"
	. "gist.github.com/6445065.git"
)

// Gets the expression as a string.
func GetExprAsString(interface{}) string {
	return GetParentArgExprAsString(0)
}

func getParent2ArgExprAllAsAst() []ast.Expr {
	// TODO: Replace use of debug.Stack() with direct use of runtime package...
	stack := string(debug.Stack())

	parentName := GetLine(stack, 5)
	parentName = parentName[1:strings.Index(parentName, ": ")]
	if dotPos := strings.LastIndex(parentName, "."); dotPos != -1 { // Trim package prefix
		parentName = parentName[dotPos+1:]
	}

	str := GetLine(stack, 7)
	str = str[strings.Index(str, ": ")+len(": "):]
	p, err := ParseStmt(str)
	if err != nil {
		return nil
	}

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
		return nil
	}
	return callExpr.Args
}

// Gets the argIndex argument expression of parent func call as a string.
func GetParentArgExprAsString(argIndex uint32) string {
	args := getParent2ArgExprAllAsAst()
	if args == nil {
		return "<expr not found>"
	}
	if argIndex >= uint32(len(args)) {
		return "<out of range>"
	}

	return SprintAstBare(args[argIndex])
}

// Gets all argument expressions of parent func call as a string.
func GetParentArgExprAllAsString() []string {
	args := getParent2ArgExprAllAsAst()
	if args == nil {
		return nil
	}

	out := make([]string, len(args))
	for i := range args {
		out[i] = SprintAstBare(args[i])
	}
	return out
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
