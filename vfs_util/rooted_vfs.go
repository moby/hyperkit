package vfs_util

import (
	"os"
	"strings"

	"code.google.com/p/go.tools/godoc/vfs"
)

type rootedFileSystem struct {
	real vfs.FileSystem
}

func NewRootedFS(real vfs.FileSystem) vfs.FileSystem {
	return &rootedFileSystem{real: real}
}

func (p *rootedFileSystem) Open(name string) (vfs.ReadSeekCloser, error) {
	return p.real.Open(unroot(name))
}

func (p *rootedFileSystem) Lstat(path string) (os.FileInfo, error) {
	return p.real.Lstat(unroot(path))
}

func (p *rootedFileSystem) Stat(path string) (os.FileInfo, error) {
	return p.real.Stat(unroot(path))
}

func (p *rootedFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	return p.real.ReadDir(unroot(path))
}

func (p *rootedFileSystem) String() string {
	return "rootedFileSystem{" + p.real.String() + "}"
}

func unroot(path string) string {
	if path == "/" {
		return "."
	}
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}
