package gist7802150_test

import (
	"fmt"

	"github.com/shurcooL/go/gists/gist7802150"
)

func ExampleFileUri_Path() {
	u := gist7802150.FileUri("file:///usr/local/go/bin/go")

	fmt.Println(u.Path())

	// Output:
	///usr/local/go/bin/go
}
