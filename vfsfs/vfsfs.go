// Package vfsfs implements vfs.FileSystem using a http.FileSystem.
package vfsfs

import (
	"net/http"
	"os"

	"golang.org/x/tools/godoc/vfs"
)

// New returns a vfs.FileSystem adapter for the provided http.FileSystem.
func New(fs http.FileSystem) vfs.FileSystem {
	return &vfsFS{fs: fs}
}

type vfsFS struct {
	fs http.FileSystem
}

func (v *vfsFS) Open(name string) (vfs.ReadSeekCloser, error) {
	return v.fs.Open(name)
}

func (v *vfsFS) Lstat(path string) (os.FileInfo, error) {
	return v.Stat(path)
}

func (v *vfsFS) Stat(path string) (os.FileInfo, error) {
	fi, err := v.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return fi.Stat()
}

func (v *vfsFS) ReadDir(path string) ([]os.FileInfo, error) {
	fi, err := v.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return fi.Readdir(0)
}

func (v *vfsFS) String() string { return "vfsfs" }
