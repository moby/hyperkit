// Package for representing paths in a structured manner.
package path

import (
	"path"
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
