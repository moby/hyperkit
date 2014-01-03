package gist7480523

import (
	"go/build"

	"github.com/shurcooL/go/vcs"

	. "gist.github.com/5504644.git"
	. "gist.github.com/7519227.git"
)

type GoPackageStringer func(*GoPackage) string

// A GoPackage represents an instance of a Go package.
type GoPackage struct {
	Bpkg *build.Package

	Vcs vcs.Vcs // TODO: Reuse same Vcs for all Go packages that are subfolders of same repo, etc.

	// TODO: These cached values should be a part of a Vcs struct or something, etc.
	Status        string
	LocalBranch   string
	Local, Remote string
}

func GoPackageFromImportPathFound(importPathFound ImportPathFound) *GoPackage {
	bpkg, err := BuildPackageFromSrcDir(importPathFound.FullPath())
	if err != nil {
		return nil
	}

	w := &GoPackage{Bpkg: bpkg}
	return w
}

func GoPackageFromImportPath(importPath string) *GoPackage {
	bpkg, err := BuildPackageFromImportPath(importPath)
	if err != nil {
		return nil
	}

	w := &GoPackage{Bpkg: bpkg}
	return w
}

func (w *GoPackage) Path() string {
	return w.Bpkg.Dir
}

func (w *GoPackage) CheckIfUnderVcs() bool {
	w.Vcs = vcs.New(w.Path())
	return w.Vcs != nil
}

func (w *GoPackage) UpdateVcsFields() {
	if w.Vcs != nil {
		w.Status = w.Vcs.GetStatus()
		w.LocalBranch = w.Vcs.GetLocalBranch()
		w.Local = w.Vcs.GetLocalRev()
		w.Remote = w.Vcs.GetRemoteRev()
	}
}
