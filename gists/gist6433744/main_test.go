package gist6433744_test

import (
	"fmt"

	"github.com/shurcooL/go/gists/gist6433744"
)

func ExampleGetLineStartEndIndicies() {
	b := []byte(`this

this is a longer line
and
stuff
last`)

	for lineIndex := 0; ; lineIndex++ {
		s, e := gist6433744.GetLineStartEndIndicies(b, lineIndex)
		fmt.Printf("%v: [%v, %v]\n", lineIndex, s, e)
		if s == -1 {
			break
		}
	}

	// Output:
	// 0: [0, 4]
	// 1: [5, 5]
	// 2: [6, 27]
	// 3: [28, 31]
	// 4: [32, 37]
	// 5: [38, 42]
	// 6: [-1, -1]
}
