// Package for representing paths in a structured manner.
package path

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Path struct {
	Elements []string
	Relative bool
}

func New(p string) *Path {
	p = path.Clean(p)
	if path.IsAbs(p) {
		return &Path{Elements: strings.Split(p[1:], "/")}
	} else {
		return &Path{Elements: strings.Split(p, "/"), Relative: true}
	}
}

func (this *Path) String() (out string) {
	if !this.Relative {
		return "/" + path.Join(this.Elements...)
	} else {
		return path.Join(this.Elements...)
	}
}

// Converts this path to a host OS native path format.
// TODO: Currently, it won't work on Windows.
func (this *Path) hostPath() (out string) {
	if !this.Relative {
		return string(filepath.Separator) + filepath.Join(this.Elements...)
	} else {
		return filepath.Join(this.Elements...)
	}
}

// Opens the file at this path as a ReadCloser.
func (this *Path) ToReadCloser() (io.ReadCloser, error) {
	return os.Open(this.hostPath())
}
