package gist5504644

import (
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	"sync"
)

// ParseFiles parses the Go source files files within directory dir
// and returns their ASTs, or the first parse error if any.
// TODO: This was made unexported in Go library, need to find a good replacement.
func ParseFiles(fset *token.FileSet, dir string, files ...string) ([]*ast.File, error) {
	var wg sync.WaitGroup
	n := len(files)
	parsed := make([]*ast.File, n, n)
	errors := make([]error, n, n)
	for i, file := range files {
		if !filepath.IsAbs(file) {
			file = filepath.Join(dir, file)
		}
		wg.Add(1)
		go func(i int, file string) {
			parsed[i], errors[i] = parser.ParseFile(fset, file, nil, 0)
			wg.Done()
		}(i, file)
	}
	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return parsed, nil
}

// ---

func AstPackageFromBuildPackage(bpkg *build.Package) (apkg *ast.Package, err error) {
	// TODO: Either find a way to use golang.org/x/tools/importer directly, or do file AST parsing in parallel like it does
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
	return build.Import(importPath, "", build.ImportComment)
}

func BuildPackageFromSrcDir(srcDir string) (bpkg *build.Package, err error) {
	return build.ImportDir(srcDir, build.ImportComment)
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
	// TODO: Either find a way to use golang.org/x/tools/importer directly, or do file AST parsing in parallel like it does
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
	if err != nil {
		panic(err)
	}
	println(dpkg.Consts[0].Names[0])
	println(dpkg.Types[0].Name)
	println(dpkg.Vars[0].Names[0])
	println(dpkg.Funcs[0].Name)
}
