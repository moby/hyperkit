// Dumps the command-line arguments.
package main

import (
	"os"

	"github.com/shurcooL/go-goon"
)

func main() {
	goon.DumpExpr(os.Args[0])  // Program name.
	goon.DumpExpr(os.Args[1:]) // Program arguments.
	goon.DumpExpr(os.Getwd())  // Current working directory.
}
