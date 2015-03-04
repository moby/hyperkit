package gist5408860_test

import (
	"fmt"
	"sort"

	"github.com/shurcooL/go/gists/gist5408860"
)

func Example() {
	s := []int{5, 2, 6, 3, 1, 4} // Unsorted.

	sort.Sort(sort.IntSlice(s))
	fmt.Println(s)

	sort.Sort(gist5408860.Reverse{Interface: sort.IntSlice(s)})
	fmt.Println(s)

	// Output:
	// [1 2 3 4 5 6]
	// [6 5 4 3 2 1]
}
