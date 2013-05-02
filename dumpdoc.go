package main

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

func GetDocPackage(ImportPath string) *doc.Package {
	bpkg, err := build.Import(ImportPath, "", 0)
	CheckError(err)
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	for _, name := range append(bpkg.GoFiles, bpkg.CgoFiles...) {
		file, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, name), nil, parser.ParseComments)
		CheckError(err)
		files[name] = file
	}
	apkg := &ast.Package{Name: bpkg.Name, Files: files}
	return doc.New(apkg, bpkg.ImportPath, 0)
}

func main() {
	dpkg := GetDocPackage("os")
	println(dpkg.Consts[0].Names[0])
	println(dpkg.Types[0].Name)
	println(dpkg.Vars[0].Names[0])
	println(dpkg.Funcs[0].Name)
}