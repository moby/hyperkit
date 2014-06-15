package gopath_util

import (
	"errors"

	. "gist.github.com/7480523.git"

	"github.com/shurcooL/go/trash"
	"github.com/shurcooL/gostatus/status"
)

// TODO: Support pattern matching like `go list`.
func Remove(importPath string) error {
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

	goPackage.UpdateVcsFields()

	notableStatus := func(goPackage *GoPackage) bool {
		// Check for notable status.
		packageStatus := status.PorcelainPresenter(goPackage)[:3] // Assumes status.PorcelainPresenter output is always at least 3 bytes.
		return packageStatus != "   " &&
			packageStatus != "  +" // Updates are okay to ignore.
	}

	if notableStatus(goPackage) {
		return errors.New("notable status: " + status.PorcelainPresenter(goPackage))
	}

	if goPackage.Dir.Repo.Vcs.RootPath() != goPackage.Bpkg.Dir {
		return errors.New("repo root path mismatch: " + goPackage.Dir.Repo.Vcs.RootPath() + " != " + goPackage.Bpkg.Dir)
	}

	err := trash.MoveToTrash(goPackage.Dir.Repo.Vcs.RootPath())
	return err
}
