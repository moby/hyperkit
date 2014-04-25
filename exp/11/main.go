// Displays Go package source code with dot imports inlined.
package exp11

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"code.google.com/p/go.tools/go/loader"
	"code.google.com/p/go.tools/imports"

	//"code.google.com/p/go.tools/go/types"
	//"honnef.co/go/importer"

	. "gist.github.com/5286084.git"
	. "gist.github.com/5504644.git"
	. "gist.github.com/5639599.git"

	"github.com/shurcooL/go/exp/15"
)

var _ = AstPackageFromBuildPackage
var _ = PrintlnAst
var _ = exp15.SortImports

const parserMode = parser.ParseComments
const astMergeMode = 0*ast.FilterFuncDuplicates | ast.FilterUnassociatedComments | ast.FilterImportDuplicates

var dotImports []*loader.PackageInfo

func findDotImports(prog *loader.Program, pi *loader.PackageInfo) {
	for _, file := range pi.Files {
		for _, importSpec := range file.Imports {
			if importSpec.Name != nil && importSpec.Name.Name == "." {
				dotImportImportPath := strings.Trim(importSpec.Path.Value, `"`)
				dotImportPi := prog.AllPackages[prog.ImportMap[dotImportImportPath]]
				dotImports = append(dotImports, dotImportPi)
				findDotImports(prog, dotImportPi)
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

	findDotImports(prog, pi)

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

	switch 3 {
	case 1:
		fmt.Fprintln(w, "package "+SprintAst(prog.Fset, merged.Name))
		fmt.Fprintln(w)
		fmt.Fprintln(w, `import (`)
		// TODO: SortImports (ala goimports).
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
	case 2:
		sortDecls(merged)

		//fmt.Fprintln(w, SprintAst(token.NewFileSet(), merged))

		//ast.SortImports(prog.Fset, merged)
		exp15.SortImports2(prog.Fset, merged)

		fmt.Fprintln(w, SprintAst(prog.Fset, merged))
	case 3:
		sortDecls(merged)

		// TODO: Clean up this mess...
		fset2, f2 := exp15.SortImports2(token.NewFileSet(), merged)

		fmt.Fprintln(w, "package "+SprintAst(prog.Fset, merged.Name))
		for _, decl := range f2.Decls {
			if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
				fmt.Fprintln(w)
				fmt.Fprintln(w, SprintAst(fset2, decl))
			}
		}
		for _, decl := range merged.Decls {
			if x, ok := decl.(*ast.GenDecl); ok && (x.Tok == token.IMPORT || x.Tok == token.PACKAGE) {
				continue
			}

			fmt.Fprintln(w)
			fmt.Fprintln(w, SprintAst(prog.Fset, decl))
		}
	case 4:
		sortDecls(merged)

		src := []byte(SprintAst(prog.Fset, merged))

		out, err := imports.Process("", src, nil)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(out)
		fmt.Println()
	}
}

func sortDecls(merged *ast.File) {
	var sortedDecls []ast.Decl
	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.PACKAGE {
			sortedDecls = append(sortedDecls, decl)
		}
	}
	/*for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && x.Tok == token.IMPORT {
			sortedDecls = append(sortedDecls, decl)
			goon.DumpExpr(decl)
		}
	}*/
	var specs []ast.Spec
	for _, importSpec := range merged.Imports {
		if importSpec.Name != nil && importSpec.Name.Name == "." {
			continue
		}
		importSpec.EndPos = 0
		specs = append(specs, importSpec)
	}
	sortedDecls = append(sortedDecls, &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: (token.Pos)(1), // Needs to be non-zero to be considered as a group.
		Specs:  specs,
	})
	//goon.DumpExpr(sortedDecls[len(sortedDecls)-1])
	for _, decl := range merged.Decls {
		if x, ok := decl.(*ast.GenDecl); ok && (x.Tok == token.IMPORT || x.Tok == token.PACKAGE) {
			continue
		}
		sortedDecls = append(sortedDecls, decl)
	}
	merged.Decls = sortedDecls
}
