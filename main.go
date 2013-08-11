package main

import (
	"io/ioutil"
	//. "gist.github.com/5286084.git"
	"fmt"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
	"go/doc"
	"os"
	"path/filepath"
)

func printPackageSummary(dpkg *doc.Package) {
	fmt.Println(`import . "` + dpkg.ImportPath + `"`)
	for _, f := range dpkg.Funcs {
		fmt.Print("\t")
		PrintlnAstBare(f.Decl)
	}
	fmt.Println()
}

func PrintPackageSummary(ImportPath string) {
	printPackageSummary(GetDocPackage(ImportPath))
}

func PrintPackageSummaryWithPath(ImportPath, fullPath string) {
	dpkg := GetDocPackage(ImportPath)
	fmt.Println(filepath.Join(fullPath, dpkg.Filenames[0]))
	printPackageSummary(dpkg)
}

func PrintPackageSummariesInDir(dirname string) {
	gopathEntries := filepath.SplitList(os.Getenv("GOPATH"))
	for _, gopathEntry := range gopathEntries {
		path0 := filepath.Join(gopathEntry, "src")
		entries, err := ioutil.ReadDir(filepath.Join(path0, dirname))
		//CheckError(err)
		if nil != err {
			continue
		}
		//for _, v := range entries {
		for i := len(entries) - 1; i >= 0; i-- {
			v := entries[i]
			if v.IsDir() {
				PrintPackageSummaryWithPath(filepath.Join(dirname, v.Name()), filepath.Join(path0, dirname, v.Name()))
				//PrintPackageSummary(filepath.Join(dirname, v.Name()))
			}
		}
	}
}

func main() {
	//PrintPackageSummary("gist.github.com/5639599.git"); return
	PrintPackageSummariesInDir("gist.github.com")
}
