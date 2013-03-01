package gist4727543

import (
	"testing"
	"fmt"
)

func TestSqrt(t *testing.T) {
	Goal := `var _ = gist4670289.GoKeywords
var _ = gist4670289.GoKeywords
var _ = NewPkgName.GoKeywords
var _ = GoKeywords
`

	// TODO: This test case should be automatically generated from dev environment...
	var Out string
	Out += fmt.Sprintln(GetForcedUse("gist.github.com/4670289.git"))
	Out += fmt.Sprintln(GetForcedUseRenamed("gist.github.com/4670289.git", ""))
	Out += fmt.Sprintln(GetForcedUseRenamed("gist.github.com/4670289.git", "NewPkgName"))
	Out += fmt.Sprintln(GetForcedUseRenamed("gist.github.com/4670289.git", "."))

	if Out != Goal {
		t.Errorf("%s is not %s", Goal, Out)
	}
}