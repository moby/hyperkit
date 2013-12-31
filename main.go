package gist7480523

import (
	"go/build"
	"os/exec"

	. "gist.github.com/5504644.git"
	. "gist.github.com/5892738.git"
	. "gist.github.com/7519227.git"
)

func GetGitRepoRoot(path string) (isGitRepo bool, rootPath string) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path

	if out, err := cmd.CombinedOutput(); err == nil {
		return true, MustTrimLastNewline(string(out))
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
		return MustTrimLastNewline(string(out))
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

// ---

// TODO: Rename this struct to something meaningful, maybe "VersionedGoPackage"... I still can't think of a good name
//       that describes what this represents (basically, a locally checked out Go Package in your GOPATH)
type Something struct {
	Bpkg          *build.Package
	Path          string
	IsGitRepo     bool
	Status        string
	LocalBranch   string
	Remote, Local string
}

func SomethingFromImportPathFound(importPathFound ImportPathFound) *Something {
	bpkg, err := BuildPackageFromSrcDir(importPathFound.FullPath())
	if err != nil {
		return nil
	}

	w := &Something{Bpkg: bpkg, Path: importPathFound.FullPath()}
	return w
}

func SomethingFromImportPath(importPath string) *Something {
	bpkg, err := BuildPackageFromImportPath(importPath)
	if err != nil {
		return nil
	}

	w := &Something{Bpkg: bpkg, Path: bpkg.Dir}
	return w
}

func (w *Something) Update() {
	w.IsGitRepo, w.Status = IsFolderGitRepo(w.Path)
	if w.IsGitRepo {
		w.LocalBranch = CheckGitRepoLocalBranch(w.Path)
		w.Remote = CheckGitRepoRemote(w.Path)
		w.Local = CheckGitRepoLocal(w.Path)
	}
}

func (w *Something) String() string {
	out := ""

	if w.IsGitRepo {
		out += "@"
		if w.LocalBranch != "master" {
			out += "b"
		} else {
			out += " "
		}
		if w.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if w.Remote != w.Local {
			out += "+"
		} else {
			out += " "
		}
	} else {
		out += "    "
	}
	if w.Bpkg.IsCommand() {
		out += "/ "
	} else {
		out += "  "
	}
	out += w.Bpkg.ImportPath

	return out
}
