// Package gist5259939 gets the filepath of go source file that calls the function.
package gist5259939

import (
	"path/filepath"
	"runtime"
)

// ThisGoSourceFile returns the full path of the Go source file where this function was called from.
func ThisGoSourceFile() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// ThisGoSourceDir returns the parent directory of the Go source file where this function was called from.
func ThisGoSourceDir() string {
	_, file, _, _ := runtime.Caller(1)
	dir, _ := filepath.Split(file)
	return dir
}
