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
	_, status := IsFolderGitRepo(this.path)
	return status
}

func (this *gitVcs) GetDefaultBranch() string {
	return "master"
}

func (this *gitVcs) GetLocalBranch() string {
	return CheckGitRepoLocalBranch(this.path)
}

func (this *gitVcs) GetLocalRev() string {
	return CheckGitRepoLocal(this.path)
}

func (this *gitVcs) GetRemoteRev() string {
	return CheckGitRepoRemote(this.path)
}

// ---

func GetGitRepoRoot(path string) (isGitRepo bool, rootPath string) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err == nil {
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

	if out, err := cmd.CombinedOutput(); err == nil {
		return true, string(out)
	} else {
		return false, ""
	}
}

func CheckGitRepoLocalBranch(path string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err == nil {
		return TrimLastNewline(string(out)) // Since rev-parse is considered porcelain and may change, need to error-check its output
	} else {
		return ""
	}
}

func CheckGitRepoLocal(path string) string {
	cmd := exec.Command("git", "rev-parse", "master")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err == nil && len(out) >= 40 {
		return string(out[:40])
	} else {
		return ""
	}
}

func CheckGitRepoRemote(path string) string {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", "master")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err == nil && len(out) >= 40 {
		return string(out[:40])
	} else {
		return ""
	}
}
