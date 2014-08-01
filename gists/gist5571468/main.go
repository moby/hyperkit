package gist5571468

import (
	"io/ioutil"

	. "gist.github.com/5286084.git"
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
