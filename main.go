package gist8018045

import "github.com/shurcooL/go-goon"
import "path/filepath"
import "os"
import "io/ioutil"
import "strings"
import "fmt"

//import . "gist.github.com/5504644.git"
//import "os/exec"
//import . "gist.github.com/5694308.git"

import . "gist.github.com/7480523.git"
import . "gist.github.com/7519227.git"

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

var skipGopath = map[string]bool{"/Users/Dmitri/Local/Ongoing/Conception/GoLand": false, "/Users/Dmitri/Dropbox/Work/2013/GoLanding": false}

func GetGoPackages(out chan<- ImportPathFound) {
	gopathEntries := filepath.SplitList(os.Getenv("GOPATH"))
	//goon.DumpExpr(gopathEntries)
	//goon.DumpExpr(build.Default.SrcDirs())
	//return
	for _, gopathEntry := range gopathEntries {
		if skipGopath[gopathEntry] {
			continue
		}

		//println("---", gopathEntry, "---\n")
		rec(out, NewImportPathFound(".", gopathEntry))
	}
	close(out)
}

func main() {
	out := make(chan ImportPathFound)
	go GetGoPackages(out)

	for importPathFound := range out {
		println(importPathFound.ImportPath())
	}
	//rec("/Users/Dmitri/Dropbox/Work/2013/GoLand/src/gist.github.com/7176504.git")
	//rec("/Users/Dmitri/Dropbox/Work/2013/GoLand/src/github.com/chsc/gogl")
	//rec("/Users/Dmitri/Dropbox/Work/2013/GoLand/src/honnef.co/go/importer")
}
