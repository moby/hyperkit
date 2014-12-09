package trim_test

import (
	"fmt"

	"github.com/shurcooL/go/trim"
)

func Example() {
	fmt.Printf("%q\n", trim.LastNewline("String\n"))
	fmt.Printf("%q\n", trim.LastNewline("String"))
	fmt.Printf("%q\n", trim.LastNewline(""))
	fmt.Printf("%q\n", trim.LastNewline("\n"))

	// Output:
	//"String"
	//"String"
	//""
	//""
}
