package gist5645828

import (
	"io/ioutil"
	//. "gist.github.com/5286084.git"
	"fmt"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
	"go/doc"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PrintPackageFullSummary(dpkg *doc.Package) {
	FprintPackageFullSummary(os.Stdout, dpkg)
}

func FprintPackageFullSummary(w io.Writer, dpkg *doc.Package) {
	for _, v := range dpkg.Vars {
		fmt.Fprintln(w, SprintAstBare(v.Decl))
	}
	for _, t := range dpkg.Types {
		for _, v := range t.Vars {
			fmt.Fprintln(w, SprintAstBare(v.Decl))
		}
	}
	fmt.Fprintln(w)
	for _, f := range dpkg.Funcs {
		fmt.Fprintln(w, SprintAstBare(f.Decl))
	}
	for _, t := range dpkg.Types {
		for _, f := range t.Funcs {
			fmt.Fprintln(w, SprintAstBare(f.Decl))
		}
		for _, m := range t.Methods {
			fmt.Fprintln(w, SprintAstBare(m.Decl))
		}
	}
	fmt.Fprintln(w)
	for _, c := range dpkg.Consts {
		//fmt.Fprintln(w, SprintAstBare(c.Decl))
		fmt.Fprintln(w, strings.Join(c.Names, "\n"))
	}
	for _, t := range dpkg.Types {
		for _, c := range t.Consts {
			//fmt.Fprintln(w, SprintAstBare(c.Decl))
			fmt.Fprintln(w, strings.Join(c.Names, "\n"))
		}
	}
	fmt.Fprintln(w)
	for _, t := range dpkg.Types {
		//fmt.Fprintln(w, SprintAstBare(t.Decl))
		fmt.Fprintln(w, t.Name)
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
	dpkg := GetDocPackage(BuildPackageFromImportPath(ImportPath))
	if len(dpkg.Funcs) == 0 {
		return
	}
	printPackageSummary(dpkg)
}

func PrintPackageSummaryWithPath(ImportPath, fullPath string) {
	dpkg := GetDocPackage(BuildPackageFromImportPath(ImportPath))
	if len(dpkg.Funcs) == 0 {
		return
	}
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
			if v.IsDir() { // TODO: Build a build.Package to figure out if this is a valid Go package; rather than assuming all dirs are
				//PrintPackageSummaryWithPath(filepath.Join(dirname, v.Name()), filepath.Join(path0, dirname, v.Name()))
				PrintPackageSummary(filepath.Join(dirname, v.Name()))
			}
		}
	}
}

func main() {
	//PrintPackageSummary("gist.github.com/5639599.git"); return
	//PrintPackageSummariesInDir("gist.github.com")
	PrintPackageFullSummary(GetDocPackageAll(BuildPackageFromImportPath("gist.github.com/5694308.git")))
}
