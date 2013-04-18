package main

import (
	"fmt"
	"sort"
)

// Reverse embeds a sort.Interface value and implements a reverse sort over
// that value.
type Reverse struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	sort.Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func main() {
	s := []int{5, 2, 6, 3, 1, 4} // unsorted
	sort.Sort(Reverse{sort.IntSlice(s)})
	fmt.Println(s)
}