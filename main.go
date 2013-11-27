package gist7519227

import (
	"path/filepath"
)

type ImportPathFound struct {
	importPath  string
	gopathEntry string
}

func NewImportPathFound(importPath, gopathEntry string) ImportPathFound {
	return ImportPathFound{
		importPath:  importPath,
		gopathEntry: gopathEntry,
	}
}

func (w *ImportPathFound) ImportPath() string {
	return w.importPath
}

func (w *ImportPathFound) GopathEntry() string {
	return w.gopathEntry
}

func (w *ImportPathFound) FullPath() string {
	return filepath.Join(w.gopathEntry, "src", w.importPath)
}
