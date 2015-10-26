package analysis_test

import (
	"fmt"
	"os"

	"github.com/shurcooL/go/analysis"
)

func ExampleIsGeneratedFile() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(analysis.IsGeneratedFile(cwd, "testdata/generated_0.go"))
	fmt.Println(analysis.IsGeneratedFile(cwd, "testdata/handcrafted_0.go"))

	// Output:
	// true <nil>
	// false <nil>
}
