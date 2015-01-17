package highlight_go_test

import (
	"fmt"
	"io"

	"github.com/shurcooL/go/gists/gist6418290"
	"github.com/shurcooL/go/highlight_go"
	"github.com/sourcegraph/syntaxhighlight"
)

// debugPrinter implements syntaxhighlight.Printer and prints the parameters its given.
type debugPrinter struct{}

func (debugPrinter) Print(w io.Writer, kind syntaxhighlight.Kind, tokText string) error {
	fmt.Println(gist6418290.GetParentFuncArgsAsString(w, kind, tokText))

	return nil
}

func ExamplePrint() {
	src := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hey there, Go.")
}
`)

	// debugPrinter implements syntaxhighlight.Printer and prints the parameters its given.
	p := debugPrinter{}

	highlight_go.Print(src, nil, p)

	// Output:
	// Print(<nil>, syntaxhighlight.Keyword, "package")
	// Print(<nil>, syntaxhighlight.Whitespace, " ")
	// Print(<nil>, syntaxhighlight.Plaintext, "main")
	// Print(<nil>, syntaxhighlight.Whitespace, "\n\n")
	// Print(<nil>, syntaxhighlight.Keyword, "import")
	// Print(<nil>, syntaxhighlight.Whitespace, " ")
	// Print(<nil>, syntaxhighlight.String, "\"fmt\"")
	// Print(<nil>, syntaxhighlight.Whitespace, "\n\n")
	// Print(<nil>, syntaxhighlight.Keyword, "func")
	// Print(<nil>, syntaxhighlight.Whitespace, " ")
	// Print(<nil>, syntaxhighlight.Plaintext, "main")
	// Print(<nil>, syntaxhighlight.Plaintext, "(")
	// Print(<nil>, syntaxhighlight.Plaintext, ")")
	// Print(<nil>, syntaxhighlight.Whitespace, " ")
	// Print(<nil>, syntaxhighlight.Plaintext, "{")
	// Print(<nil>, syntaxhighlight.Whitespace, "\n\t")
	// Print(<nil>, syntaxhighlight.Plaintext, "fmt")
	// Print(<nil>, syntaxhighlight.Plaintext, ".")
	// Print(<nil>, syntaxhighlight.Plaintext, "Println")
	// Print(<nil>, syntaxhighlight.Plaintext, "(")
	// Print(<nil>, syntaxhighlight.String, "\"Hey there, Go.\"")
	// Print(<nil>, syntaxhighlight.Plaintext, ")")
	// Print(<nil>, syntaxhighlight.Whitespace, "\n")
	// Print(<nil>, syntaxhighlight.Plaintext, "}")
	// Print(<nil>, syntaxhighlight.Whitespace, "\n")
}

func ExamplePrintWhitespace() {
	src := []byte("  package    main      \n\t\n")

	highlight_go.Print(src, nil, debugPrinter{})

	// Output:
	// Print(<nil>, syntaxhighlight.Whitespace, "  ")
	// Print(<nil>, syntaxhighlight.Keyword, "package")
	// Print(<nil>, syntaxhighlight.Whitespace, "    ")
	// Print(<nil>, syntaxhighlight.Plaintext, "main")
	// Print(<nil>, syntaxhighlight.Whitespace, "      \n\t\n")
}
