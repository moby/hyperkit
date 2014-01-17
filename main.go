package gist5504644

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"

	. "gist.github.com/5286084.git"
)

func AstPackageFromBuildPackage(bpkg *build.Package) (apkg *ast.Package, err error) {
	// TODO: Either find a way to use code.google.com/p/go.tools/importer directly, or do file AST parsing in parallel like it does
	filenames := append(bpkg.GoFiles, bpkg.CgoFiles...)
	files := make(map[string]*ast.File, len(filenames))
	fset := token.NewFileSet()
	for _, filename := range filenames {
		fileAst, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, filename), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[filename] = fileAst // TODO: Figure out if filename or full path are to be used (the key of this map doesn't seem to be used anywhere!)
	}
	return &ast.Package{Name: bpkg.Name, Files: files}, nil
}

func BuildPackageFromImportPath(importPath string) (bpkg *build.Package, err error) {
	return build.Import(importPath, "", 0)
}

func BuildPackageFromSrcDir(srcDir string) (bpkg *build.Package, err error) {
	return build.ImportDir(srcDir, 0)
}

func getDocPackageMode(bpkg *build.Package, err error, mode doc.Mode) (dpkg *doc.Package, err2 error) {
	if err != nil {
		return nil, err
	}
	apkg, err := AstPackageFromBuildPackage(bpkg)
	if err != nil {
		return nil, err
	}
	return doc.New(apkg, bpkg.ImportPath, mode), nil
}

func GetDocPackage(bpkg *build.Package, err error) (dpkg *doc.Package, err2 error) {
	return getDocPackageMode(bpkg, err, 0)
}

func GetDocPackageAll(bpkg *build.Package, err error) (dpkg *doc.Package, err2 error) {
	return getDocPackageMode(bpkg, err, doc.AllDecls|doc.AllMethods)
}

/* Commented out because it's not in use anywhere, candidate for removal
func GetDocPackageFromFiles(paths ...string) (dpkg *doc.Package) {
	// TODO: Either find a way to use code.google.com/p/go.tools/importer directly, or do file AST parsing in parallel like it does
	files := make(map[string]*ast.File, len(paths))
	fset := token.NewFileSet()
	for _, path := range paths {
		fileAst, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		CheckError(err)
		files[filepath.Base(path)] = fileAst // TODO: Figure out if filename or full path are to be used (the key of this map doesn't seem to be used anywhere!)
	}
	// TODO: Figure out an import path, package name
	apkg := &ast.Package{Name: "bpkg.Name", Files: files}
	return doc.New(apkg, "ImportPath", 0)
}*/

func main() {
	dpkg, err := GetDocPackage(BuildPackageFromImportPath("os"))
	CheckError(err)
	println(dpkg.Consts[0].Names[0])
	println(dpkg.Types[0].Name)
	println(dpkg.Vars[0].Names[0])
	println(dpkg.Funcs[0].Name)
}
