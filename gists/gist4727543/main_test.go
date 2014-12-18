package gist4727543_test

import (
	"fmt"

	"github.com/shurcooL/go/gists/gist4727543"
)

func Example() {
	fmt.Println(gist4727543.GetForcedUse("io/ioutil"))
	fmt.Println(gist4727543.GetForcedUseRenamed("io/ioutil", ""))
	fmt.Println(gist4727543.GetForcedUseRenamed("io/ioutil", "RenamedPkg"))
	fmt.Println(gist4727543.GetForcedUseRenamed("io/ioutil", "."))
	fmt.Println()
	fmt.Println(gist4727543.GetForcedUseFromImport(`github.com/shurcooL/go/gists/gist4727543`))
	fmt.Println(gist4727543.GetForcedUseFromImport(`"github.com/shurcooL/go/gists/gist4727543"`))
	fmt.Println(gist4727543.GetForcedUseFromImport("`github.com/shurcooL/go/gists/gist4727543`"))
	fmt.Println(gist4727543.GetForcedUseFromImport(`. "github.com/shurcooL/go/gists/gist4727543"`))
	fmt.Println(gist4727543.GetForcedUseFromImport(`renamed "github.com/shurcooL/go/gists/gist4727543"`))
	fmt.Println(gist4727543.GetForcedUseFromImport(`bad`))
	fmt.Println(gist4727543.GetForcedUseFromImport(`bad bad bad`))

	// Output:
	// var _ = ioutil.NopCloser
	// var _ = ioutil.NopCloser
	// var _ = RenamedPkg.NopCloser
	// var _ = NopCloser
	//
	// var _ = gist4727543.GetForcedUse
	// var _ = gist4727543.GetForcedUse
	// var _ = gist4727543.GetForcedUse
	// var _ = GetForcedUse
	// var _ = renamed.GetForcedUse
	// Package "bad" not valid (doesn't exist or can't be built).
	// Invalid import string.
}
