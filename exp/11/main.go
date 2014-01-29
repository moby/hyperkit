// Displays Go package source code with dot imports inlined.
package exp11

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"code.google.com/p/go.tools/go/loader"
	//"code.google.com/p/go.tools/go/types"
	//"honnef.co/go/importer"

	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"
)

var _ = AstPackageFromBuildPackage
var _ = PrintlnAst

const parserMode = parser.ParseComments
const astMergeMode = 0*ast.FilterFuncDuplicates | ast.FilterUnassociatedComments | ast.FilterImportDuplicates

var imports map[string]*loader.PackageInfo
var dotImports []*loader.PackageInfo

func findDotImports(pi *loader.PackageInfo) {
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
	/*imp2 := importer.New()
	imp2.Config.UseGcFallback = true
	cfg := types.Config{Import: imp2.Import}
	_ = cfg*/

	conf := loader.Config{
		//TypeChecker:   cfg,
		SourceImports: true,
	}

	conf.Import(importPath)

	prog, err := conf.Load()
	CheckError(err)

	/*pi, err := imp.ImportPackage(importPath)
	CheckError(err)
	_ = pi*/

	pi := prog.Imported[importPath]

	// Create ImportPath -> *PackageInfo map
	imports = make(map[string]*loader.PackageInfo, len(prog.AllPackages))
	for _, pi := range prog.AllPackages {
		imports[pi.Pkg.Path()] = pi
	}

	findDotImports(pi)

	files := make(map[string]*ast.File)
	{
		// This package
		for _, file := range pi.Files {
			filename := prog.Fset.File(file.Package).Name()
			files[filename] = file
		}

		// All dot imports
		for _, pi := range dotImports {
			for _, file := range pi.Files {
				filename := prog.Fset.File(file.Package).Name()
				files[filename] = file
			}
		}
	}

	apkg := &ast.Package{Name: pi.Pkg.Name(), Files: files}

	merged := ast.MergePackageFiles(apkg, astMergeMode)

	fmt.Fprintln(w, "package "+SprintAst(prog.Fset, merged.Name))
	fmt.Fprintln(w)
	fmt.Fprintln(w, `import (`)
	for _, importSpec := range merged.Imports {
		if importSpec.Name != nil && importSpec.Name.Name == "." {
			continue
		}
		fmt.Fprintln(w, "\t"+SprintAst(prog.Fset, importSpec))
	}
	fmt.Fprintln(w, `)`)
	fmt.Fprintln(w)

	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
			continue
		}

		fmt.Fprintln(w, SprintAst(prog.Fset, decl))
		fmt.Fprintln(w)
	}

	// TODO: Make this work equivalent to above
	//fmt.Fprintln(w, SprintAst(imp.Fset, merged))

	//goon.Dump(merged)
}
