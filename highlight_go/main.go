// Package highlight_go provides a syntax highlighter for Go, using go/scanner.
package highlight_go

import (
	"go/scanner"
	"go/token"
	"io"

	"github.com/sourcegraph/annotate"
	"github.com/sourcegraph/syntaxhighlight"
)

func tokenKind(tok token.Token, lit string) int {
	switch {
	case tok.IsKeyword() || (tok.IsOperator() && tok < token.LPAREN):
		return syntaxhighlight.KEYWORD

	// Literals.
	case tok == token.INT || tok == token.FLOAT || tok == token.IMAG:
		return syntaxhighlight.DECIMAL
	case tok == token.STRING || tok == token.CHAR:
		return syntaxhighlight.STRING
	case lit == "true" || lit == "false" || lit == "iota" || lit == "nil":
		return syntaxhighlight.LITERAL

	case tok == token.COMMENT:
		return syntaxhighlight.COMMENT
	default:
		return syntaxhighlight.PLAINTEXT
	}
}

func Print(src []byte, w io.Writer, p syntaxhighlight.Printer) error {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	var lastOffset int

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		offset := int(fset.Position(pos).Offset)

		var tokString string
		if lit != "" {
			tokString = lit
		} else {
			tokString = tok.String()
		}

		// TODO: Clean this up.
		//if tok == token.SEMICOLON {
		if tok == token.SEMICOLON && lit == "\n" {
			continue
		}

		// TODO: Clean this up.
		whitespace := string(src[lastOffset:offset])
		lastOffset = offset + len(tokString)
		tokString = whitespace + tokString

		err := p.Print(w, tokenKind(tok, lit), tokString)
		if err != nil {
			return err
		}
	}

	return nil
}

func Annotate(src []byte, a syntaxhighlight.Annotator) (annotate.Annotations, error) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	var anns annotate.Annotations

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		offset := int(fset.Position(pos).Offset)

		var tokString string
		if lit != "" {
			tokString = lit
		} else {
			tokString = tok.String()
		}

		// TODO: Clean this up.
		//if tok == token.SEMICOLON {
		if tok == token.SEMICOLON && lit == "\n" {
			continue
		}

		ann, err := a.Annotate(offset, tokenKind(tok, lit), tokString)
		if err != nil {
			return nil, err
		}
		if ann != nil {
			anns = append(anns, ann)
		}
	}

	return anns, nil
}
