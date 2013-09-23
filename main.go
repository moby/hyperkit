package gist5645828

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

/*
$ doc doc.Package
http://golang.org/pkg/go/doc/#Package
/usr/local/go/src/pkg/go/doc/doc.go:14:
// Package is the documentation for an entire package.
type Package struct {
	Doc		string
	Name		string
	ImportPath	string
	Imports		[]string
	Filenames	[]string
	Notes		map[string][]*Note
	// DEPRECATED. For backward compatibility Bugs is still populated,
	// but all new code should use Notes instead.
	Bugs	[]string

	// declarations
	Consts	[]*Value
	Types	[]*Type
	Vars	[]*Value
	Funcs	[]*Func
}
*/
func PrintPackageFullSummary(dpkg *doc.Package) {
	for _, v := range dpkg.Vars {
		PrintlnAstBare(v.Decl)
	}
	fmt.Println()
	for _, f := range dpkg.Funcs {
		PrintlnAstBare(f.Decl)
	}
	fmt.Println()
	for _, c := range dpkg.Consts {
		PrintlnAstBare(c.Decl)
	}
	fmt.Println()
	for _, t := range dpkg.Types {
		PrintlnAstBare(t.Decl)
	}
}

func printPackageSummary(dpkg *doc.Package) {
	fmt.Println(`import . "` + dpkg.ImportPath + `"`)
	for _, f := range dpkg.Funcs {
		fmt.Print("\t")
		PrintlnAstBare(f.Decl)
	}
	fmt.Println()
}

func PrintPackageSummary(ImportPath string) {
	printPackageSummary(GetDocPackage(BuildPackageFromImportPath(ImportPath)))
}

func PrintPackageSummaryWithPath(ImportPath, fullPath string) {
	dpkg := GetDocPackage(BuildPackageFromImportPath(ImportPath))
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
	//PrintPackageSummariesInDir("gist.github.com")
	PrintPackageFullSummary(GetDocPackageAll(BuildPackageFromImportPath("gist.github.com/5694308.git")))
}
