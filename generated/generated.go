// Package generated provides a function that parses a Go file and reports
// whether it contains a "// Code generated ... DO NOT EDIT." line comment.
//
// It's intended to stay up to date with the specification proposal in
// https://golang.org/issues/13560.
//
// The first priority is correctness (no false negatives, no false positives).
// It must return accurate results even if the input Go source code is not gofmted.
// The second priority is performance. The current version is implemented
// via go/parser, but it may be possible to improve performance via an
// alternative implementation. That can be explored later.
//
// The exact API is undecided and can change. The current API style is somewhat
// based on go/parser, but that may not be the best approach.
package generated

import (
	"go/parser"
	"go/token"
	"strings"
)

// ParseFile parses the source code of a single Go source file
// specified by filename, and reports whether the file contains
// a "// Code generated ... DO NOT EDIT." line comment
// matching the specification proposal in
// https://golang.org/issues/13560#issuecomment-277804473:
//
// 	The text must appear as the first line of a properly formatted Go // comment,
// 	and that comment must appear before but not be attached to the package clause
// 	and before any /* */ comment. This is similar to the rules for build tags.
//
// 	The comment line must match the case-sensitive regular expression (in Go syntax):
//
// 		`^// Code generated .* DO NOT EDIT\.$`
//
// 	The .* means the tool can put whatever folderol it wants in there,
// 	but the comment must be a single line and must start with `Code generated`
// 	and end with `DO NOT EDIT.`, with a period.
//
// If the source couldn't be read, the error indicates the specific
// failure. If the source was read but syntax errors were found,
// the result is estimated on a best effort basis from a partial AST.
//
// TODO: Decide on best policy of what to do in case of syntax errors
//       being encountered during parsing.
func ParseFile(filename string) (hasGeneratedComment bool, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly|parser.ParseComments)
	if f == nil { // Can only happen when err != nil.
		return false, err
	}
Outer:
	for _, cg := range f.Comments {
		if cg == f.Doc {
			// If we've reached the package comment, don't look any further,
			// because the generated comment must be before that.
			break
		}
		// Check if this comment group is a match.
		// The text must appear as the first line of a properly formatted line comment (//-style).
		if len(cg.List[0].Text) >= smallestMatchingComment &&
			strings.HasPrefix(cg.List[0].Text, "// Code generated ") &&
			strings.HasSuffix(cg.List[0].Text, " DO NOT EDIT.") &&
			fset.Position(cg.List[0].Pos()).Column == 1 {

			return true, nil
		}
		// Ensure none of the comments in this comment group are general comments (/*-style).
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "/*") {
				// If we've reached a general comment (/*-style), don't look any further,
				// because the generated comment must be before that.
				break Outer
			}
		}
	}
	return false, nil
}

const smallestMatchingComment = len("// Code generated  DO NOT EDIT.")
