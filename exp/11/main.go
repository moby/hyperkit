// Displays Go package source code with dot imports inlined.
package exp11

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"code.google.com/p/go.tools/go/types"
	"code.google.com/p/go.tools/importer"
	importer2 "honnef.co/go/importer"

	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
)

var _ = AstPackageFromBuildPackage
var _ = PrintlnAst

const parserMode = parser.ParseComments
const astMergeMode = 0*ast.FilterFuncDuplicates | ast.FilterUnassociatedComments | ast.FilterImportDuplicates

var imports map[string]*importer.PackageInfo
var dotImports []*importer.PackageInfo

func findDotImports(pi *importer.PackageInfo) {
	for _, file := range pi.Files {
		for _, importSpec := range file.Imports {
			if importSpec.Name != nil && importSpec.Name.Name == "." {
				importPath := strings.Trim(importSpec.Path.Value, `"`)
				dotImports = append(dotImports, imports[importPath])
				findDotImports(imports[importPath])
			}
		}
	}
}

func InlineDotImports(w io.Writer, importPath string) {
	imp2 := importer2.New()
	imp2.Config.UseGcFallback = true
	cfg := types.Config{Import: imp2.Import}
	_ = cfg

	imp := importer.New(&importer.Config{
		//TypeChecker:   cfg,
		SourceImports: true,
	})

	pi, err := imp.ImportPackage(importPath)
	CheckError(err)
	_ = pi

	// Create ImportPath -> *PackageInfo map
	imports = make(map[string]*importer.PackageInfo, len(imp.AllPackages()))
	for _, pi := range imp.AllPackages() {
		imports[pi.Pkg.Path()] = pi
	}

	findDotImports(pi)

	files := make(map[string]*ast.File)
	{
		// This package
		for _, file := range pi.Files {
			filename := imp.Fset.File(file.Package).Name()
			files[filename] = file
		}

		// All dot imports
		for _, pi := range dotImports {
			for _, file := range pi.Files {
				filename := imp.Fset.File(file.Package).Name()
				files[filename] = file
			}
		}
	}

	apkg := &ast.Package{Name: pi.Pkg.Name(), Files: files}

	merged := ast.MergePackageFiles(apkg, astMergeMode)

	fmt.Fprintln(w, "package "+SprintAst(imp.Fset, merged.Name))
	fmt.Fprintln(w)
	fmt.Fprintln(w, `import (`)
	for _, importSpec := range merged.Imports {
		if importSpec.Name != nil && importSpec.Name.Name == "." {
			continue
		}
		fmt.Fprintln(w, "\t"+SprintAst(imp.Fset, importSpec))
	}
	fmt.Fprintln(w, `)`)
	fmt.Fprintln(w)

	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
			continue
		}

		fmt.Fprintln(w, SprintAst(imp.Fset, decl))
		fmt.Fprintln(w)
	}

	// TODO: Make this work equivalent to above
	//fmt.Fprintln(w, SprintAst(imp.Fset, merged))

	//goon.Dump(merged)
}
