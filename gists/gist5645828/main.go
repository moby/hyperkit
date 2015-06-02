package gist5645828

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shurcooL/go/gists/gist5504644"
	"github.com/shurcooL/go/gists/gist5639599"
)

type sectionWriter struct {
	Writer         io.Writer
	writtenSection bool
}

func (sw *sectionWriter) Write(p []byte) (n int, err error) {
	sw.writtenSection = true
	return sw.Writer.Write(p)
}

func (sw *sectionWriter) WriteBreak() {
	if sw.writtenSection {
		io.WriteString(sw.Writer, "\n")
		sw.writtenSection = false
	}
}

func PrintPackageFullSummary(dpkg *doc.Package) {
	FprintPackageFullSummary(os.Stdout, dpkg)
}

func FprintPackageFullSummary(wr io.Writer, dpkg *doc.Package) {
	w := &sectionWriter{Writer: wr}
	for _, v := range dpkg.Vars {
		for _, spec := range v.Decl.Specs {
			spec.(*ast.ValueSpec).Doc = nil
			spec.(*ast.ValueSpec).Comment = nil
		}
		fmt.Fprintln(w, gist5639599.SprintAstBare(v.Decl))
	}
	for _, t := range dpkg.Types {
		for _, v := range t.Vars {
			for _, spec := range v.Decl.Specs {
				spec.(*ast.ValueSpec).Doc = nil
				spec.(*ast.ValueSpec).Comment = nil
			}
			fmt.Fprintln(w, gist5639599.SprintAstBare(v.Decl))
		}
	}
	w.WriteBreak()
	for _, f := range dpkg.Funcs {
		fmt.Fprintln(w, gist5639599.SprintAstBare(f.Decl))
	}
	for _, t := range dpkg.Types {
		for _, f := range t.Funcs {
			fmt.Fprintln(w, gist5639599.SprintAstBare(f.Decl))
		}
		for _, m := range t.Methods {
			fmt.Fprintln(w, gist5639599.SprintAstBare(m.Decl))
		}
	}
	w.WriteBreak()
	for _, c := range dpkg.Consts {
		for _, spec := range c.Decl.Specs {
			spec.(*ast.ValueSpec).Values = nil
			spec.(*ast.ValueSpec).Doc = nil
			spec.(*ast.ValueSpec).Comment = nil
		}
		fmt.Fprintln(w, gist5639599.SprintAstBare(c.Decl))
		//fmt.Fprintln(w, "const", strings.Join(c.Names, "\n"))
	}
	for _, t := range dpkg.Types {
		for _, c := range t.Consts {
			for _, spec := range c.Decl.Specs {
				spec.(*ast.ValueSpec).Values = nil
				spec.(*ast.ValueSpec).Doc = nil
				spec.(*ast.ValueSpec).Comment = nil
			}
			fmt.Fprintln(w, gist5639599.SprintAstBare(c.Decl))
			//fmt.Fprintln(w, "const", strings.Join(c.Names, "\n"))
		}
	}
	w.WriteBreak()
	for _, t := range dpkg.Types {
		//fmt.Fprintln(w, gist5639599.SprintAstBare(t.Decl))
		fmt.Fprintln(w, "type", t.Name)
	}
}

func printPackageSummary(dpkg *doc.Package) {
	fmt.Println(`import "` + dpkg.ImportPath + `"`)
	for _, f := range dpkg.Funcs {
		fmt.Print("\t")
		gist5639599.PrintlnAstBare(f.Decl)
	}
	for _, t := range dpkg.Types {
		for _, f := range t.Funcs {
			fmt.Print("\t")
			gist5639599.PrintlnAstBare(f.Decl)
		}
		// THINK: Do I want to include methods?
		/*for _, m := range t.Methods {
			fmt.Print("\t")
			gist5639599.PrintlnAstBare(m.Decl)
		}*/
	}
	fmt.Println()
}

// hasAnyFuncs returns true if the package has any funcs or methods.
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

func PrintPackageSummary(importPath string) {
	dpkg, err := gist5504644.GetDocPackage(gist5504644.BuildPackageFromImportPath(importPath))
	if err != nil {
		return
	}
	if !hasAnyFuncs(dpkg) {
		return
	}
	printPackageSummary(dpkg)
}

func PrintPackageSummaryWithPath(importPath, fullPath string) {
	dpkg, err := gist5504644.GetDocPackage(gist5504644.BuildPackageFromImportPath(importPath))
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
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
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
	dpkg, err := gist5504644.GetDocPackageAll(gist5504644.BuildPackageFromImportPath("github.com/microcosm-cc/bluemonday"))
	if err != nil {
		panic(err)
	}
	PrintPackageFullSummary(dpkg)
}
