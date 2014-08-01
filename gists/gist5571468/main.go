// Package gist5571468 reads the content of a file, panics on error.
package gist5571468

import (
	"io/ioutil"

	. "github.com/shurcooL/go/gists/gist5286084"
)

func MustReadFile(filename string) string {
	return string(MustReadFileB(filename))
}

func MustReadFileB(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	CheckError(err)
	return b
}

func TryReadFile(filename string) string {
	return string(TryReadFileB(filename))
}

func TryReadFileB(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return b
}
