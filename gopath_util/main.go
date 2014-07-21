package gopath_util

import (
	"errors"
	"strings"

	. "gist.github.com/7480523.git"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/go/trash"
	"github.com/shurcooL/gostatus/status"
)

// Moves go gettable repo with no local changes into trash.
// importPathPattern must match exactly with the repo root.
// For example, "github.com/user/repo/...".
func RemoveRepo(importPathPattern string) error {
	// TODO: Use an official Go package for `go list` functionality whenever possible.
	importPaths := gotool.ImportPaths([]string{importPathPattern})
	if len(importPaths) == 0 {
		return errors.New("no packages to remove")
	}

	var firstGoPackage *GoPackage
	for i, importPath := range importPaths {
		goPackage := GoPackageFromImportPath(importPath)
		if goPackage == nil {
			return errors.New("Import Path not found: " + importPath)
		}

		if goPackage.Bpkg.Goroot {
			return errors.New("can't remove packages from GOROOT")
		}

		goPackage.UpdateVcs()

		if goPackage.Dir.Repo == nil {
			return errors.New("can't get repo status")
		}

		if i == 0 {
			firstGoPackage = goPackage
		} else if firstGoPackage.Dir.Repo != goPackage.Dir.Repo {
			return errors.New("matched Go Packages span more than 1 repo: " + firstGoPackage.Dir.Repo.Vcs.RootPath() + " != " + goPackage.Dir.Repo.Vcs.RootPath())
		} else if !strings.HasPrefix(goPackage.Bpkg.Dir, firstGoPackage.Dir.Repo.Vcs.RootPath()) { // TODO: This is probably not neccessary...
			return errors.New("Go Package not inside repo: " + goPackage.Bpkg.Dir + " doesn't have prefix " + firstGoPackage.Dir.Repo.Vcs.RootPath())
		}
	}

	if repoImportPathPattern := GetRepoImportPathPattern(firstGoPackage.Dir.Repo.Vcs.RootPath(), firstGoPackage.Bpkg.SrcRoot); repoImportPathPattern != importPathPattern {
		return errors.New("importPathPattern not exact repo root match: " + importPathPattern + " != " + repoImportPathPattern)
	}

	firstGoPackage.UpdateVcsFields()

	cleanStatus := func(goPackage *GoPackage) bool {
		packageStatus := status.PlumbingPresenterV2(goPackage)[:4]
		return packageStatus == "    " || packageStatus == "  + " // Updates are okay to ignore.
	}

	if !cleanStatus(firstGoPackage) {
		return errors.New("non-clean status: " + status.PorcelainPresenter(firstGoPackage))
	}

	err := trash.MoveToTrash(firstGoPackage.Dir.Repo.Vcs.RootPath())
	return err
}
