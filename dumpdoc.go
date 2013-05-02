package main

import (
	"github.com/davecgh/go-spew/spew"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("%s importPath", os.Args[0])
	}
	bpkg, err := build.Import(os.Args[1], "", 0)
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	for _, name := range append(bpkg.GoFiles, bpkg.CgoFiles...) {
		file, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, name), nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		files[name] = file
	}
	apkg := &ast.Package{Name: bpkg.Name, Files: files}
	dpkg := doc.New(apkg, bpkg.ImportPath, 0)
	spew.Dump(dpkg)
}
