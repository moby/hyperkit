package gist5092053_test

import (
	"fmt"

	"github.com/shurcooL/go/gists/gist5092053"
)

func Example() {
	m := map[string]int{
		"blah": 5,
		"boo":  9,
		"yah":  1,
	}

	sm := gist5092053.SortMapByValue(m)

	for _, v := range sm {
		fmt.Println(v.Value, v.Key)
	}

	// Output:
	// 9 boo
	// 5 blah
	// 1 yah
}
