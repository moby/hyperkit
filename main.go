package gist8018045

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/shurcooL/go-goon"

	. "gist.github.com/5504644.git"
	. "gist.github.com/7480523.git"
	. "gist.github.com/7519227.git"
)

var _ = fmt.Print
var _ = goon.Dump

func rec(out chan<- ImportPathFound, importPathFound ImportPathFound) {
	if goPackage := GoPackageFromImportPathFound(importPathFound); goPackage != nil {
		out <- importPathFound
	}

	entries, err := ioutil.ReadDir(importPathFound.FullPath())
	if err == nil {
		for _, v := range entries {
			if v.IsDir() && !strings.HasPrefix(v.Name(), ".") {
				rec(out, NewImportPathFound(filepath.Join(importPathFound.ImportPath(), v.Name()), importPathFound.GopathEntry()))
			}
		}
	}
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

var skipGopath = map[string]bool{"/Users/Dmitri/Local/Ongoing/Conception/GoLand": false, "/Users/Dmitri/Dropbox/Work/2013/GoLanding": false}

func GetGoPackages(out chan<- ImportPathFound) {
	getGoPackagesB(out)
}

func getGoPackagesA(out chan<- ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	//goon.DumpExpr(gopathEntries)
	//goon.DumpExpr(build.Default.SrcDirs())
	//return

	for _, gopathEntry := range gopathEntries {
		/*if skipGopath[gopathEntry] {
			continue
		}*/

		//println("---", gopathEntry, "---\n")
		rec(out, NewImportPathFound(".", gopathEntry))
	}
	close(out)
}

func getGoPackagesB(out chan<- ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		root := filepath.Join(gopathEntry, "src")
		if !isDir(root) {
			continue
		}

		_ = filepath.Walk(root, func(path string, fi os.FileInfo, _ error) error {
			if !fi.IsDir() {
				return nil
			}
			if strings.HasPrefix(fi.Name(), ".") {
				return filepath.SkipDir
			}
			importPath, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			importPathFound := NewImportPathFound(importPath, gopathEntry)
			if goPackage := GoPackageFromImportPathFound(importPathFound); goPackage != nil {
				out <- importPathFound
			}
			return nil
		})
	}
	close(out)
}

func getGoPackagesC(out chan<- ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		root := filepath.Join(gopathEntry, "src")
		if !isDir(root) {
			continue
		}

		_ = filepath.Walk(root, func(path string, fi os.FileInfo, _ error) error {
			if !fi.IsDir() {
				return nil
			}
			if strings.HasPrefix(fi.Name(), ".") {
				return filepath.SkipDir
			}
			bpkg, err := BuildPackageFromSrcDir(path)
			if err != nil {
				return nil
			}
			/*if bpkg.Goroot {
				return nil
			}*/
			out <- NewImportPathFound(bpkg.ImportPath, bpkg.Root)
			return nil
		})
	}
	close(out)
}

func main() {
	started := time.Now()

	out := make(chan ImportPathFound)
	go GetGoPackages(out)

	for importPathFound := range out {
		_ = importPathFound
		println(importPathFound.ImportPath())
		//goon.Dump(importPathFound)
	}

	goon.Dump(time.Since(started).Seconds() * 1000)
}
