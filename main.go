package gist5645828

import (
	"fmt"
	"go/ast"
	"go/doc"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
)

func PrintPackageFullSummary(dpkg *doc.Package) {
	FprintPackageFullSummary(os.Stdout, dpkg)
}

func FprintPackageFullSummary(w io.Writer, dpkg *doc.Package) {
	for _, v := range dpkg.Vars {
		for _, spec := range v.Decl.Specs {
			spec.(*ast.ValueSpec).Comment = nil
		}
		fmt.Fprintln(w, SprintAstBare(v.Decl))
	}
	for _, t := range dpkg.Types {
		for _, v := range t.Vars {
			for _, spec := range v.Decl.Specs {
				spec.(*ast.ValueSpec).Comment = nil
			}
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
		for _, spec := range c.Decl.Specs {
			spec.(*ast.ValueSpec).Values = nil
			spec.(*ast.ValueSpec).Comment = nil
		}
		fmt.Fprintln(w, SprintAstBare(c.Decl))
		//fmt.Fprintln(w, "const", strings.Join(c.Names, "\n"))
	}
	for _, t := range dpkg.Types {
		for _, c := range t.Consts {
			for _, spec := range c.Decl.Specs {
				spec.(*ast.ValueSpec).Values = nil
				spec.(*ast.ValueSpec).Comment = nil
			}
			fmt.Fprintln(w, SprintAstBare(c.Decl))
			//fmt.Fprintln(w, "const", strings.Join(c.Names, "\n"))
		}
	}
	fmt.Fprintln(w)
	for _, t := range dpkg.Types {
		//fmt.Fprintln(w, SprintAstBare(t.Decl))
		fmt.Fprintln(w, "type", t.Name)
	}
}

func printPackageSummary(dpkg *doc.Package) {
	fmt.Println(`import . "` + dpkg.ImportPath + `"`)
	for _, f := range dpkg.Funcs {
		fmt.Print("\t")
		PrintlnAstBare(f.Decl)
	}
	for _, t := range dpkg.Types {
		for _, f := range t.Funcs {
			fmt.Print("\t")
			PrintlnAstBare(f.Decl)
		}
		// THINK: Do I want to include methods?
		/*for _, m := range t.Methods {
			fmt.Print("\t")
			PrintlnAstBare(m.Decl)
		}*/
	}
	fmt.Println()
}

func hasAnyFuncs(dpkg *doc.Package) bool {
	if len(dpkg.Funcs) > 0 {
		return true
	}

	for _, t := range dpkg.Types {
		if len(t.Funcs) > 0 {
			return true
		}
	}

	return false
}

func PrintPackageSummary(ImportPath string) {
	dpkg, err := GetDocPackage(BuildPackageFromImportPath(ImportPath))
	if err != nil {
		return
	}
	if !hasAnyFuncs(dpkg) {
		return
	}
	printPackageSummary(dpkg)
}

func PrintPackageSummaryWithPath(ImportPath, fullPath string) {
	dpkg, err := GetDocPackage(BuildPackageFromImportPath(ImportPath))
	if err != nil {
		return
	}
	if !hasAnyFuncs(dpkg) {
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
		if err != nil {
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
	dpkg, err := GetDocPackageAll(BuildPackageFromImportPath("github.com/shurcooL/Conception-go"))
	if err != nil {
		panic(err)
	}
	PrintPackageFullSummary(dpkg)
}
