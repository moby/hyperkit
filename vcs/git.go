package vcs

import (
	"os/exec"

	. "gist.github.com/5892738.git"
)

type gitVcs struct {
	commonVcs
}

func (this *gitVcs) Type() Type { return Git }

func (this *gitVcs) GetStatus() string {
	_, status := IsFolderGitRepo(this.rootPath)
	return status
}

func (this *gitVcs) GetRemote() string {
	cmd := exec.Command("git", "ls-remote", "--get-url", "origin")
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil {
		return TrimLastNewline(string(out))
	} else {
		return ""
	}
}

func (this *gitVcs) GetDefaultBranch() string {
	return "master"
}

func (this *gitVcs) GetLocalBranch() string {
	return CheckGitRepoLocalBranch(this.rootPath)
}

func (this *gitVcs) GetLocalRev() string {
	return CheckGitRepoLocal(this.rootPath, this.GetDefaultBranch())
}

func (this *gitVcs) GetRemoteRev() string {
	return CheckGitRepoRemote(this.rootPath, this.GetDefaultBranch())
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

func CheckGitRepoLocal(path, branch string) string {
	cmd := exec.Command("git", "rev-parse", branch)
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil && len(out) >= gitRevisionLength {
		return string(out[:gitRevisionLength])
	} else {
		return ""
	}
}

func CheckGitRepoRemote(path, branch string) string {
	// true here is not a boolean value, but a command /bin/true that will make git think it asked for a password,
	// and prevent potential interactive password prompts (opting to return failure exit code instead).
	cmd := exec.Command("git", "-c", "core.askpass=true", "ls-remote", "--heads", "origin", branch)
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil && len(out) >= gitRevisionLength {
		return string(out[:gitRevisionLength])
	} else {
		return ""
	}
}
