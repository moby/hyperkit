package path

import "fmt"

func Example() {
	paths := []*Path{
		&Path{[]string{"path", "to", "somewhere"}, false},
		&Path{[]string{"src", "github.com"}, true},
	}

	for _, p := range paths {
		fmt.Println(p.String())
	}

	// Output:
	// /path/to/somewhere
	// src/github.com
}
