// Package vfsioutil implements some I/O utility functions for vfs.FileSystem.
package vfsioutil

import (
	"bytes"
	"io"

	"golang.org/x/tools/godoc/vfs"
)

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func ReadFile(fs vfs.FileSystem, filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
