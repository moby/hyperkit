package gist5259939

import (
	. "gist.github.com/5258650.git"
	"path/filepath"
	"runtime/debug"
	"strings"
	//. "gist.github.com/5504644.git"
)

// Gets the full path of the Go source file where this function was called from
func GetThisGoSourceFilepath() string {
	path := GetLine(string(debug.Stack()), 2)
	path = path[0:strings.Index(path, ":")]
	return path
}

// Gets the parent directory of the Go source file where this function was called from
func GetThisGoSourceDir() string {
	path := GetLine(string(debug.Stack()), 2)
	path = path[0:strings.Index(path, ":")]
	path, _ = filepath.Split(path)
	return path
}

func main() {
	println(GetThisGoSourceFilepath())
	println(GetThisGoSourceDir())
	//bpkg := BuildPackageFromSrcDir(GetThisGoSourceDir())
	//println(bpkg.ImportPath, bpkg.Name)
}
