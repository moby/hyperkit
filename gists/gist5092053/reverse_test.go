package gist5092053

import (
	"fmt"
	"sort"
)

func Example_reverseAdapter() {
	s := []int{5, 2, 6, 3, 1, 4} // Unsorted.

	sort.Sort(sort.IntSlice(s))
	fmt.Println(s)

	sort.Sort(reverseAdapter{Interface: sort.IntSlice(s)})
	fmt.Println(s)

	// Output:
	// [1 2 3 4 5 6]
	// [6 5 4 3 2 1]
}
