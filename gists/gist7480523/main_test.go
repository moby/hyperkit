package gist7480523_test

import (
	"fmt"

	"github.com/shurcooL/go/gists/gist7480523"
)

func ExampleGetRepoImportPathPattern() {
	fmt.Println(gist7480523.GetRepoImportPathPattern("/home/User/Go/src/github.com/owner/repo", "/home/User/Go/src"))
	fmt.Println(gist7480523.GetRepoImportPathPattern("/home/user/go/src/github.com/owner/repo", "/home/User/Go/src"))

	// Output:
	//github.com/owner/repo/...
	//github.com/owner/repo/...
}
