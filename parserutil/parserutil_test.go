package parserutil_test

import (
	"go/ast"
	"os"

	"github.com/shurcooL/go/parserutil"
)

func Example() {
	stmt, err := parserutil.ParseStmt("var x int")
	if err != nil {
		panic(err)
	}

	ast.Fprint(os.Stdout, nil, stmt, nil)

	// Output:
	//      0  *ast.DeclStmt {
	//      1  .  Decl: *ast.GenDecl {
	//      2  .  .  Doc: nil
	//      3  .  .  TokPos: 31
	//      4  .  .  Tok: var
	//      5  .  .  Lparen: 0
	//      6  .  .  Specs: []ast.Spec (len = 1) {
	//      7  .  .  .  0: *ast.ValueSpec {
	//      8  .  .  .  .  Doc: nil
	//      9  .  .  .  .  Names: []*ast.Ident (len = 1) {
	//     10  .  .  .  .  .  0: *ast.Ident {
	//     11  .  .  .  .  .  .  NamePos: 35
	//     12  .  .  .  .  .  .  Name: "x"
	//     13  .  .  .  .  .  .  Obj: *ast.Object {
	//     14  .  .  .  .  .  .  .  Kind: var
	//     15  .  .  .  .  .  .  .  Name: "x"
	//     16  .  .  .  .  .  .  .  Decl: *(obj @ 7)
	//     17  .  .  .  .  .  .  .  Data: 0
	//     18  .  .  .  .  .  .  .  Type: nil
	//     19  .  .  .  .  .  .  }
	//     20  .  .  .  .  .  }
	//     21  .  .  .  .  }
	//     22  .  .  .  .  Type: *ast.Ident {
	//     23  .  .  .  .  .  NamePos: 37
	//     24  .  .  .  .  .  Name: "int"
	//     25  .  .  .  .  .  Obj: nil
	//     26  .  .  .  .  }
	//     27  .  .  .  .  Values: nil
	//     28  .  .  .  .  Comment: nil
	//     29  .  .  .  }
	//     30  .  .  }
	//     31  .  .  Rparen: 0
	//     32  .  }
	//     33  }
}
