package gist5504644

import (
	. "gist.github.com/5286084.git"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	//"github.com/davecgh/go-spew/spew"
)

func astPackage(bpkg *build.Package) *ast.Package {
	files := make(map[string]*ast.File)
	fset := token.NewFileSet()
	for _, name := range append(bpkg.GoFiles, bpkg.CgoFiles...) {
		file, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, name), nil, parser.ParseComments)
		CheckError(err)
		files[name] = file
	}
	return &ast.Package{Name: bpkg.Name, Files: files}
}

func BuildPackageFromImportPath(ImportPath string) *build.Package {
	bpkg, err := build.Import(ImportPath, "", 0)
	CheckError(err)
	return bpkg
}

func BuildPackageFromSrcDir(SrcDir string) *build.Package {
	bpkg, err := build.Import(".", SrcDir, 0)
	CheckError(err)
	return bpkg
}

func GetDocPackage(bpkg *build.Package) *doc.Package {
	apkg := astPackage(bpkg)
	return doc.New(apkg, bpkg.ImportPath, 0)
}

func GetDocPackageAll(bpkg *build.Package) *doc.Package {
	apkg := astPackage(bpkg)
	return doc.New(apkg, bpkg.ImportPath, doc.AllDecls) // TODO: Is doc.AllMethods needed also?
}

func main() {
	dpkg := GetDocPackage(BuildPackageFromImportPath("os"))
	println(dpkg.Consts[0].Names[0])
	println(dpkg.Types[0].Name)
	println(dpkg.Vars[0].Names[0])
	println(dpkg.Funcs[0].Name)
}
