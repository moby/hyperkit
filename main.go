package gist7480523

import (
	"go/build"
	"strings"

	"github.com/shurcooL/go/vcs"

	. "gist.github.com/5504644.git"
	. "gist.github.com/7519227.git"
)

type GoPackageStringer func(*GoPackage) string

// A GoPackage describes a single package found in a directory.
// This is partially a copy of "cmd/go".Package, except it can be imported and reused. =.=
// https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#24
type GoPackage struct {
	Bpkg     *build.Package
	Standard bool // is this package part of the standard Go library?

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
	return goPackageFromBuildPackage(bpkg)
}

func GoPackageFromImportPath(importPath string) *GoPackage {
	bpkg, err := BuildPackageFromImportPath(importPath)
	if err != nil {
		return nil
	}
	return goPackageFromBuildPackage(bpkg)
}

func goPackageFromBuildPackage(bpkg *build.Package) *GoPackage {
	goPackage := &GoPackage{
		Bpkg:     bpkg,
		Standard: bpkg.Goroot && bpkg.ImportPath != "" && !strings.Contains(bpkg.ImportPath, "."), // https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#110
	}
	if goPackage.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		goPackage.Vcs = vcs.New(goPackage.Bpkg.Dir)
	}

	return goPackage
}

func (this *GoPackage) UpdateVcsFields() {
	if this.Vcs != nil {
		this.Status = this.Vcs.GetStatus()
		this.LocalBranch = this.Vcs.GetLocalBranch()
		this.Local = this.Vcs.GetLocalRev()
		this.Remote = this.Vcs.GetRemoteRev()
	}
}
