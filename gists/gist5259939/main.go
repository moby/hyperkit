package gist5259939

import (
	"path/filepath"
	"runtime"
	//. "gist.github.com/5504644.git"
)

// Gets the full path of the Go source file where this function was called from
func GetThisGoSourceFilepath() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// Gets the parent directory of the Go source file where this function was called from
func GetThisGoSourceDir() string {
	_, file, _, _ := runtime.Caller(1)
	dir, _ := filepath.Split(file)
	return dir
}

func main() {
	println(GetThisGoSourceFilepath())
	println(GetThisGoSourceDir())
	//bpkg := BuildPackageFromSrcDir(GetThisGoSourceDir())
	//println(bpkg.ImportPath, bpkg.Name)
}
