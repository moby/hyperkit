// Package vfs_util provides a debug Virtual File System. It prints the method names and parameters as they're called.
package vfs_util

import (
	"log"
	"os"

	"golang.org/x/tools/godoc/vfs"

	"github.com/shurcooL/go/gists/gist6418290"
)

type debugFileSystem struct {
	real vfs.FileSystem
}

func NewDebugFS(real vfs.FileSystem) *debugFileSystem {
	return &debugFileSystem{real: real}
}

func (dfs *debugFileSystem) Open(name string) (vfs.ReadSeekCloser, error) {
	log.Println(gist6418290.GetParentFuncArgsAsString(name))
	return dfs.real.Open(name)
}

func (dfs *debugFileSystem) Lstat(path string) (os.FileInfo, error) {
	log.Println(gist6418290.GetParentFuncArgsAsString(path))
	return dfs.real.Lstat(path)
}

func (dfs *debugFileSystem) Stat(path string) (os.FileInfo, error) {
	log.Println(gist6418290.GetParentFuncArgsAsString(path))
	return dfs.real.Stat(path)
}

func (dfs *debugFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	log.Println(gist6418290.GetParentFuncArgsAsString(path))
	return dfs.real.ReadDir(path)
}

func (dfs *debugFileSystem) String() string {
	log.Println(gist6418290.GetParentFuncArgsAsString())
	return dfs.real.String()
}
