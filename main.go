package main

import (
	. "gist.github.com/5258650.git"
	"runtime/debug"
	"strings"
)

func GetThisGoSourceFilepath() string {
	x := GetLine(string(debug.Stack()), 2)
	x = x[0:strings.Index(x, ":")]
	return x
}

func main() {
	print(GetThisGoSourceFilepath())
}