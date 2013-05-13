package main

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

func main() {
}