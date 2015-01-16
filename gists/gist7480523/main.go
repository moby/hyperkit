package gist7480523

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/go/exp/12"
	"golang.org/x/tools/go/vcs"

	. "github.com/shurcooL/go/gists/gist5504644"
	. "github.com/shurcooL/go/gists/gist7802150"
)

type GoPackageStringer func(*GoPackage) string

// A GoPackage describes a single package found in a directory.
// This is partially a copy of "cmd/go".Package, except it can be imported and reused. =.=
// https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#24
type GoPackage struct {
	Bpkg    *build.Package
	BpkgErr error

	Dir *exp12.Directory
}

func GoPackageFromImportPathFound(importPathFound ImportPathFound) *GoPackage {
	bpkg, err := BuildPackageFromSrcDir(importPathFound.FullPath())
	return goPackageFromBuildPackage(bpkg, err)
}

func GoPackageFromImportPath(importPath string) *GoPackage {
	bpkg, err := BuildPackageFromImportPath(importPath)
	return goPackageFromBuildPackage(bpkg, err)
}

func GoPackageFromPath(path, srcDir string) (*GoPackage, error) {
	bpkg, err := BuildPackageFromPath(path, srcDir)
	if err != nil {
		if _, noGo := err.(*build.NoGoError); noGo || bpkg.Dir == "" {
			return nil, err
		}
	}
	return goPackageFromBuildPackage(bpkg, err), nil
}

func goPackageFromBuildPackage(bpkg *build.Package, bpkgErr error) *GoPackage {
	if bpkgErr != nil {
		if _, noGo := bpkgErr.(*build.NoGoError); noGo || bpkg.Dir == "" {
			return nil
		}
	}

	if bpkg.ConflictDir != "" {
		fmt.Fprintf(os.Stderr, "warning: ConflictDir=%q (Dir=%q)\n", bpkg.ConflictDir, bpkg.Dir)
		return nil
	}

	goPackage := &GoPackage{
		Bpkg:    bpkg,
		BpkgErr: bpkgErr,

		Dir: exp12.LookupDirectory(bpkg.Dir),
	}

	/*if goPackage.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		// TODO: markAsNotNeedToUpdate() because of external insight?
	}*/

	return goPackage
}

// This is okay to call concurrently (a mutex is used internally).
// Actually, not completely okay because MakeUpdated technology is not thread-safe.
func (this *GoPackage) UpdateVcs() {
	if this.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		MakeUpdated(this.Dir)
	}
}

func (this *GoPackage) UpdateVcsFields() {
	if this.Dir.Repo != nil {
		MakeUpdated(this.Dir.Repo.VcsLocal)
		MakeUpdated(this.Dir.Repo.VcsRemote)

		repoImportPath := GetRepoImportPath(this.Dir.Repo.Vcs.RootPath(), this.Bpkg.SrcRoot)
		if repoRoot, err := vcs.RepoRootForImportPath(repoImportPath, false); err == nil {
			this.Dir.Repo.RepoRoot = repoRoot
		}
	}
}

func GetRepoImportPath(repoPath, srcRoot string) string {
	if s, err := filepath.EvalSymlinks(repoPath); err == nil {
		repoPath = s
	} else {
		fmt.Fprintln(os.Stderr, "warning: GetRepoImportPath: can't resolve symlink")
	}
	if s, err := filepath.EvalSymlinks(srcRoot); err == nil {
		srcRoot = s
	} else {
		fmt.Fprintln(os.Stderr, "warning: GetRepoImportPath: can't resolve symlink")
	}

	// Detect and handle case mismatch in prefix.
	if prefixLen := len(srcRoot + "/"); len(repoPath) >= prefixLen && srcRoot+"/" != repoPath[:prefixLen] && strings.EqualFold(srcRoot+"/", repoPath[:prefixLen]) {
		fmt.Fprintln(os.Stderr, "warning: GetRepoImportPath: prefix case doesn't match")
		return repoPath[prefixLen:]
	}

	return strings.TrimPrefix(repoPath, srcRoot+"/")
}
func GetRepoImportPathPattern(repoPath, srcRoot string) (repoImportPathPattern string) {
	return GetRepoImportPath(repoPath, srcRoot) + "/..."
}

func (this *GoPackage) String() string {
	return this.Bpkg.ImportPath
}

// =====

// GoPackageRepo represents a collection of Go packages contained by one VCS repository.
type GoPackageRepo struct {
	rootPath   string
	goPackages []*GoPackage
}

func NewGoPackageRepo(rootPath string, goPackages []*GoPackage) GoPackageRepo {
	return GoPackageRepo{rootPath, goPackages}
}

// RepoImportPath returns what would be the import path of the root folder of the repository. It may or may not
// be an actual Go package. E.g.,
//
//	"github.com/owner/repo"
func (repo GoPackageRepo) RepoImportPath() string {
	return GetRepoImportPath(repo.rootPath, repo.goPackages[0].Bpkg.SrcRoot)
}

// ImportPathPattern returns an import path pattern that matches all of the Go packages in this repo.
// E.g.,
//
//	"github.com/owner/repo/..."
func (repo GoPackageRepo) ImportPathPattern() string {
	return GetRepoImportPathPattern(repo.rootPath, repo.goPackages[0].Bpkg.SrcRoot)
}

// RootPath returns the path to the root workdir folder of the repository.
func (repo GoPackageRepo) RootPath() string         { return repo.rootPath }
func (repo GoPackageRepo) GoPackages() []*GoPackage { return repo.goPackages }

// ImportPaths returns a newline separated list of all import paths.
func (repo GoPackageRepo) ImportPaths() string {
	var importPaths []string
	for _, goPackage := range repo.goPackages {
		importPaths = append(importPaths, goPackage.Bpkg.ImportPath)
	}
	return strings.Join(importPaths, "\n")
}
