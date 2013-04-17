package main

import (
	. "gist.github.com/5258650.git"
	"runtime/debug"
	"strings"
	"path/filepath"
)

// Gets the full path of the Go source file where this function was called from
func GetThisGoSourceFilepath() string {
	x := GetLine(string(debug.Stack()), 2)
	x = x[0:strings.Index(x, ":")]
	return x
}

// Gets the parent directory of the Go source file where this function was called from
func GetThisGoSourceDir() string {
	x := GetThisGoSourceFilepath()
	x, _ = filepath.Split(x)
	return x
}

func main() {
	println(GetThisGoSourceFilepath())
	print(GetThisGoSourceDir())
}