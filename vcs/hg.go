package vcs

import (
	"os/exec"
	"strings"

	trim "github.com/shurcooL/go/trim"
)

func getHgRepoRoot(path string) (isHgRepo bool, rootPath string) {
	cmd := exec.Command("hg", "root")
	cmd.Dir = path

	if out, err := cmd.Output(); err == nil {
		return true, trim.LastNewline(string(out))
	} else {
		return false, ""
	}
}

type hgVcs struct {
	commonVcs
}

func (this *hgVcs) Type() Type { return Hg }

func (this *hgVcs) GetStatus() string {
	cmd := exec.Command("hg", "status")
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil {
		return string(out)
	} else {
		return ""
	}
}

func (this *hgVcs) GetStash() string {
	// TODO: Does Mercurial have stashes? Figure it out, add support, etc.
	return ""
}

func (this *hgVcs) GetRemote() string {
	cmd := exec.Command("hg", "paths", "default")
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil {
		return trim.LastNewline(string(out))
	} else {
		return ""
	}
}

func (this *hgVcs) GetDefaultBranch() string {
	return "default"
}

func (this *hgVcs) GetLocalBranch() string {
	cmd := exec.Command("hg", "branch")
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil {
		return trim.LastNewline(string(out))
	} else {
		return ""
	}
}

// Length of a Mercurial revision hash.
const hgRevisionLength = 40

func (this *hgVcs) GetLocalRev() string {
	// Alternative: hg parent --template '{node}'
	cmd := exec.Command("hg", "--debug", "identify", "-i", "--rev", this.GetDefaultBranch())
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil && len(out) >= hgRevisionLength {
		return string(out[:hgRevisionLength])
	} else {
		return ""
	}
}

func (this *hgVcs) GetRemoteRev() string {
	// TODO: Make this more robust and proper, etc.
	cmd := exec.Command("hg", "--debug", "identify", "-i", "--rev", this.GetDefaultBranch(), "default")
	cmd.Dir = this.rootPath

	if out, err := cmd.Output(); err == nil {
		// Get the last line of output.
		lines := strings.Split(trim.LastNewline(string(out)), "\n") // Always returns at least 1 element.
		return lines[len(lines)-1]
	}
	return ""
}

func (this *hgVcs) IsContained(rev string) bool {
	// TODO.
	return false
}
