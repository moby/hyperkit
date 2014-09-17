package u6

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/shurcooL/go/exp/13"
	"github.com/shurcooL/go/pipe_util"
	"github.com/shurcooL/go/vcs"
	"gopkg.in/pipe.v2"

	. "github.com/shurcooL/go/gists/gist5892738"
	. "github.com/shurcooL/go/gists/gist7480523"
)

// Show the difference between the working directory and the most recent commit.
// Precondition is that goPackage.Dir.Repo is not nil, and VcsLocal is updated.
// TODO: Support for non-git.
func GoPackageWorkingDiff(goPackage *GoPackage) string {
	// git diff
	if goPackage.Dir.Repo.VcsLocal.Status != "" {
		switch goPackage.Dir.Repo.Vcs.Type() {
		case vcs.Git:
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

// Branches returns a Markdown table of branches with ahead/behind information relative to master branch.
func Branches(repo *exp13.VcsState) string {
	switch repo.Vcs.Type() {
	case vcs.Git:
		branchInfo := func(line []byte) []byte {
			branch := TrimLastNewline(string(line))

			cmd := exec.Command("git", "rev-list", "--count", "--left-right", "master..."+branch)
			cmd.Dir = repo.Vcs.RootPath()
			out, err := cmd.Output()
			if err != nil && len(out) == 0 {
				return []byte(err.Error())
			}

			behindAhead := bytes.Split(out, []byte("\t"))

			if branch == repo.VcsLocal.LocalBranch {
				return []byte(fmt.Sprintf("**%s** | %s | %s", branch, string(behindAhead[0]), string(behindAhead[1])))
			} else {
				return []byte(fmt.Sprintf("%s | %s | %s", branch, string(behindAhead[0]), string(behindAhead[1])))
			}
		}

		p := pipe.Script(
			pipe.Println("Branch | Behind | Ahead"),
			pipe.Println("-------|-------:|:-----"),
			pipe.Line(
				pipe.Exec("git", "for-each-ref", "--format=%(refname:short)", "refs/heads"),
				pipe.Replace(branchInfo),
			),
		)

		out, err := pipe_util.OutputDir(p, repo.Vcs.RootPath())
		if err != nil {
			return err.Error()
		}
		return string(out)
	default:
		return ""
	}
}
