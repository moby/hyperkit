package vfs_util

import (
	"errors"
	"os"
	"path"
	"strings"

	"sourcegraph.com/sourcegraph/go-vcs/vcs/util"

	"golang.org/x/tools/godoc/vfs"
)

type prefixFileSystem struct {
	real   vfs.FileSystem
	prefix string
}

func NewPrefixFS(real vfs.FileSystem, prefix string) *prefixFileSystem {
	return &prefixFileSystem{real: real, prefix: path.Clean(prefix)}
}

func (p *prefixFileSystem) Open(name string) (vfs.ReadSeekCloser, error) {
	if strings.HasPrefix(name, p.prefix) {
		return p.real.Open(p.innerPath(name))
	}
	return nil, errors.New(name + " doesn't exist")
}

func (p *prefixFileSystem) Lstat(name string) (os.FileInfo, error) {
	return p.Stat(name)
}

func (p *prefixFileSystem) Stat(name string) (os.FileInfo, error) {
	if strings.HasPrefix(name, p.prefix) {
		return p.real.Stat(p.innerPath(name))
	}

	if !strings.HasPrefix(p.prefix, name) {
		return nil, errors.New(name + " doesn't exist")
	}

	// TODO.
	return &util.FileInfo{
		Name_: path.Base(name),
		Mode_: os.ModeDir,
		/*Size_: 0,
		ModTime_ : time.Now().UTC(),
		Sys_: nil,*/
	}, nil
}

func (p *prefixFileSystem) ReadDir(name string) ([]os.FileInfo, error) {
	if strings.HasPrefix(name, p.prefix) {
		return p.real.ReadDir(p.innerPath(name))
		/*fis, err := p.real.ReadDir(p.innerPath(name))
		goon.DumpExpr(len(fis))
		goon.DumpExpr(fis[0].Name())
		goon.DumpExpr(fis[0].Size())
		goon.DumpExpr(fis[0].Mode())
		goon.DumpExpr(fis[0].ModTime())
		goon.DumpExpr(fis[0].IsDir())
		goon.DumpExpr(fis[0].Sys())
		return fis, err*/
	}

	if !strings.HasPrefix(p.prefix, name) {
		return nil, errors.New(name + " doesn't exist")
	}

	// TODO.
	return []os.FileInfo{&util.FileInfo{
		Name_: antibase(strings.TrimPrefix(p.prefix, name)),
		Mode_: os.ModeDir,
		/*Size_: 0,
		ModTime_ : time.Now().UTC(),
		Sys_: nil,*/
	}}, nil
}

func (p *prefixFileSystem) String() string {
	return "prefixFileSystem{" + p.real.String() + "}"
}

// innerPath returns an inner path for outer path.
//
// Precondition is that strings.HasPrefix(outerPath, p.prefix) is true.
func (p *prefixFileSystem) innerPath(outerPath string) string {
	path := outerPath[len(p.prefix):]
	if path == "" {
		path = "/"
	}
	return path
}

func antibase(name string) string {
	name = strings.TrimPrefix(name, "/")
	if i := strings.Index(name, "/"); i != -1 {
		name = name[:i]
	}
	return name
}
