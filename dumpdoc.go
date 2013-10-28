package gist5504644

import (
	. "gist.github.com/5286084.git"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
)

func astPackage(bpkg *build.Package) *ast.Package {
	// TODO: Either find a way to use code.google.com/p/go.tools/importer directly, or do file AST parsing in parallel like it does
	files := make(map[string]*ast.File)
	fset := token.NewFileSet()
	for _, name := range append(bpkg.GoFiles, bpkg.CgoFiles...) {
		file, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, name), nil, parser.ParseComments)
		CheckError(err)
		files[name] = file
	}
	return &ast.Package{Name: bpkg.Name, Files: files}
}

func BuildPackageFromImportPath(ImportPath string) (bpkg *build.Package, err error) {
	return build.Import(ImportPath, "", 0)
}

func BuildPackageFromSrcDir(SrcDir string) (bpkg *build.Package, err error) {
	return build.Import(".", SrcDir, 0)
}

func GetDocPackage(bpkg *build.Package, err error) *doc.Package {
	CheckError(err)
	apkg := astPackage(bpkg)
	return doc.New(apkg, bpkg.ImportPath, 0)
}

func GetDocPackageAll(bpkg *build.Package, err error) *doc.Package {
	CheckError(err)
	apkg := astPackage(bpkg)
	return doc.New(apkg, bpkg.ImportPath, doc.AllDecls) // TODO: Is doc.AllMethods needed also?
}

func GetDocPackageFromFiles(filePaths ...string) *doc.Package {
	// TODO: Either find a way to use code.google.com/p/go.tools/importer directly, or do file AST parsing in parallel like it does
	files := make(map[string]*ast.File)
	fset := token.NewFileSet()
	for _, name := range filePaths {
		file, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		CheckError(err)
		files[name] = file
	}
	// TODO: Figure out an import path, package name
	apkg := &ast.Package{Name: "bpkg.Name", Files: files}
	return doc.New(apkg, "ImportPath", 0)
}

func main() {
	dpkg := GetDocPackage(BuildPackageFromImportPath("os"))
	println(dpkg.Consts[0].Names[0])
	println(dpkg.Types[0].Name)
	println(dpkg.Vars[0].Names[0])
	println(dpkg.Funcs[0].Name)
}
