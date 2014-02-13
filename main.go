package gist7480523

import (
	"go/build"
	"strings"
	"github.com/shurcooL/go/exp/12"

	. "gist.github.com/5504644.git"
	. "gist.github.com/7519227.git"
	. "gist.github.com/7802150.git"
)

type GoPackageStringer func(*GoPackage) string

// A GoPackage describes a single package found in a directory.
// This is partially a copy of "cmd/go".Package, except it can be imported and reused. =.=
// https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#24
type GoPackage struct {
	Bpkg     *build.Package
	Standard bool // is this package part of the standard Go library?

	Dir *exp12.Directory
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

		Dir: exp12.NewDirectory(bpkg.Dir),
	}

	/*if goPackage.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		// TODO: markAsNotNeedToUpdate() because of external insight?
	}*/

	return goPackage
}

func (this *GoPackage) UpdateVcs() {
	if this.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		MakeUpdated(this.Dir)
	}
}

func (this *GoPackage) UpdateVcsFields() {
	if this.Dir.Repo != nil {
		MakeUpdated(this.Dir.Repo.VcsLocal)
		MakeUpdated(this.Dir.Repo.VcsRemote)
	}
}

func (this *GoPackage) String() string {
	return this.Bpkg.ImportPath
}
