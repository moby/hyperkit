package u6

import (
	"os/exec"

	"github.com/shurcooL/go/vcs"

	. "gist.github.com/7480523.git"
)

// TODO: Currently, it shows modified files, but not new files. Fix that.
// TODO: Support for non-git.
func GoPackageWorkingDiff(goPackage *GoPackage) string {
	// git diff
	if goPackage.Dir.Repo.VcsLocal.Status != "" {
		if goPackage.Dir.Repo.Vcs.Type() == vcs.Git {
			cmd := exec.Command("git", "diff", "--no-ext-diff")
			cmd.Dir = goPackage.Dir.Repo.Vcs.RootPath()
			if outputBytes, err := cmd.CombinedOutput(); err == nil {
				return string(outputBytes)
			} else {
				return err.Error()
			}
		}
	}
	return "gitDiff error"
}
