package sanitized_anchor_name_test

import (
	"fmt"

	"github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name"
)

func ExampleCreate() {
	anchorName := sanitized_anchor_name.Create("This is a header")

	fmt.Println(anchorName)

	// Output:
	//this-is-a-header
}
