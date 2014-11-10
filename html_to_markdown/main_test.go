package html_to_markdown_test

import (
	"bytes"
	"testing"

	"github.com/russross/blackfriday"
	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/go/html_to_markdown"
	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	const markdownIn = `This is Markdown.

- Item
- Item
- Item

And here are some ordered items.

1. Item 1
1. Item 2
1. Item 3

Cool.
`

	renderedHtml := blackfriday.MarkdownBasic([]byte(markdownIn))

	t.Log(string(renderedHtml))

	parsedHtml, err := html.Parse(bytes.NewReader(renderedHtml))
	if err != nil {
		panic(err)
	}

	markdownOut := html_to_markdown.Unnamed1(parsedHtml)

	if markdownIn != markdownOut {
		goon.DumpExpr(markdownIn, markdownOut)
		t.Fail()
	}
}

/* More test material:

<strong>Blah blah: </strong>blah blah blah <strong>blah blah.</strong> Blah blah.

*/
