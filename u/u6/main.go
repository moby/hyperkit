package u6

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
				if len(out) > 0 {
					// diff exits with a non-zero status when the files don't match.
					// Ignore that failure as long as we get output.
					err = nil
				}
				if err != nil {
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
			if err != nil {
				log.Printf("error running %v: %v\n", cmd.Args, err)
				return []byte(fmt.Sprintf("%s | ? | ?\n", branch))
			}

			behindAhead := strings.Split(TrimLastNewline(string(out)), "\t")

			if branch == repo.VcsLocal.LocalBranch {
				return []byte(fmt.Sprintf("**%s** | %s | %s\n", branch, behindAhead[0], behindAhead[1]))
			} else {
				return []byte(fmt.Sprintf("%s | %s | %s\n", branch, behindAhead[0], behindAhead[1]))
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

// Branches returns a Markdown table of branches with ahead/behind information relative to upstream.
func BranchesRemote(repo *exp13.VcsState) string {
	switch repo.Vcs.Type() {
	case vcs.Git:
		branchInfo := func(line []byte) []byte {
			branchUpstream := strings.Split(TrimLastNewline(string(line)), "\t")
			if len(branchUpstream) != 2 {
				return []byte("error: len(branchUpstream) != 2")
			}

			branch := branchUpstream[0]
			upstream := branchUpstream[1]
			if upstream == "" {
				return []byte(fmt.Sprintf("%s | | | \n", branch))
			}

			cmd := exec.Command("git", "rev-list", "--count", "--left-right", upstream+"..."+branch)
			cmd.Dir = repo.Vcs.RootPath()
			out, err := cmd.Output()
			if err != nil {
				// This usually happens when the upstream branch is gone.
				return []byte(fmt.Sprintf("%s | ~~%s~~ | | \n", branch, upstream))
			}

			behindAhead := strings.Split(TrimLastNewline(string(out)), "\t")

			if branch == repo.VcsLocal.LocalBranch {
				return []byte(fmt.Sprintf("**%s** | %s | %s | %s\n", branch, upstream, behindAhead[0], behindAhead[1]))
			} else {
				return []byte(fmt.Sprintf("%s | %s | %s | %s\n", branch, upstream, behindAhead[0], behindAhead[1]))
			}
		}

		p := pipe.Script(
			pipe.Println("Branch | Upstream | Behind | Ahead"),
			pipe.Println("-------|----------|-------:|:-----"),
			pipe.Line(
				pipe.Exec("git", "for-each-ref", "--format=%(refname:short)\t%(upstream:short)", "refs/heads"),
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
