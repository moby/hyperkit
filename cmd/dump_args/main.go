// Dumps the command-line arguments.
package main

import (
	"os"

	"github.com/shurcooL/go-goon"
)

func main() {
	goon.DumpExpr(os.Args)
}
