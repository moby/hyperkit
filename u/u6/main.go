package u6

import (
	"os/exec"

	"github.com/shurcooL/go/pipe_util"
	"github.com/shurcooL/go/vcs"
	"gopkg.in/pipe.v2"

	. "gist.github.com/5892738.git"
	. "gist.github.com/7480523.git"
)

// Show the difference between the working directory and the most recent commit.
// Precondition is that goPackage.Dir.Repo is not nil, and VcsLocal is updated.
// TODO: Support for non-git.
func GoPackageWorkingDiff(goPackage *GoPackage) string {
	// git diff
	if goPackage.Dir.Repo.VcsLocal.Status != "" {
		if goPackage.Dir.Repo.Vcs.Type() == vcs.Git {
			/*cmd := exec.Command("git", "diff", "--no-ext-diff", "HEAD")
			cmd.Dir = goPackage.Dir.Repo.Vcs.RootPath()

			if out, err := cmd.CombinedOutput(); err == nil {
				return string(out)
			} else {
				return err.Error()
			}*/

			newFileDiff := func(line []byte) []byte {
				cmd := exec.Command("git", "diff", "--no-ext-diff", "--", "/dev/null", TrimLastNewline(string(line)))
				cmd.Dir = goPackage.Dir.Repo.Vcs.RootPath()
				out, err := cmd.Output()
				if err != nil && len(out) == 0 {
					return []byte(err.Error())
				}
				return out
			}

			p := pipe.Script(
				pipe.Exec("git", "diff", "--no-ext-diff", "HEAD"),
				pipe.Line(
					pipe.Exec("git", "ls-files", "--others", "--exclude-standard"),
					pipe.Replace(newFileDiff),
				),
			)

			out, err := pipe_util.OutputDir(p, goPackage.Dir.Repo.Vcs.RootPath())
			if err != nil {
				return err.Error()
			}
			return string(out)
		}
	}
	return ""
}
