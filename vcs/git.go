package vcs

import (
	"os/exec"

	. "gist.github.com/5892738.git"
)

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		addVcsProvider(func(path string) Vcs {
			if isRepo, rootPath := GetGitRepoRoot(path); isRepo {
				return &gitVcs{commonVcs{rootPath: rootPath}}
			}
			return nil
		})
	}
}

type gitVcs struct {
	commonVcs
}

func (this *gitVcs) Type() Type { return Git }

func (this *gitVcs) GetStatus() string {
	_, status := IsFolderGitRepo(this.rootPath)
	return status
}

func (this *gitVcs) GetDefaultBranch() string {
	return "master"
}

func (this *gitVcs) GetLocalBranch() string {
	return CheckGitRepoLocalBranch(this.rootPath)
}

func (this *gitVcs) GetLocalRev() string {
	return CheckGitRepoLocal(this.rootPath)
}

func (this *gitVcs) GetRemoteRev() string {
	return CheckGitRepoRemote(this.rootPath)
}

// ---

func GetGitRepoRoot(path string) (isGitRepo bool, rootPath string) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil {
		return true, TrimLastNewline(string(out)) // Since rev-parse is considered porcelain and may change, need to error-check its output
	} else {
		return false, ""
	}
}

func IsFolderGitRepo(path string) (isGitRepo bool, status string) {
	// Alternative: git rev-parse
	// For individual files: git ls-files --error-unmatch -- 'Filename', return code == 0
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil {
		return true, string(out)
	} else {
		return false, ""
	}
}

func CheckGitRepoLocalBranch(path string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil {
		return TrimLastNewline(string(out)) // Since rev-parse is considered porcelain and may change, need to error-check its output
	} else {
		return ""
	}
}

// Length of a git revision hash.
const gitRevisionLength = 40

func CheckGitRepoLocal(path string) string {
	cmd := exec.Command("git", "rev-parse", "master")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil && len(out) >= gitRevisionLength {
		return string(out[:gitRevisionLength])
	} else {
		return ""
	}
}

func CheckGitRepoRemote(path string) string {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", "master")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil && len(out) >= gitRevisionLength {
		return string(out[:gitRevisionLength])
	} else {
		return ""
	}
}
