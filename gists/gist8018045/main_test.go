package gist8018045

import (
	"fmt"
	"time"

	"github.com/shurcooL/go/gists/gist7480523"
)

func ExampleGetGoPackages() {
	started := time.Now()

	out := make(chan *gist7480523.GoPackage)
	go GetGoPackages(out)
	for goPackage := range out {
		fmt.Println(goPackage.Bpkg.ImportPath)
	}

	fmt.Println("time taken:", time.Since(started).Seconds()*1000, "ms")
}
