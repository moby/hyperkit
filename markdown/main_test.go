package markdown_test

import (
	"os"

	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/markdown"
)

func Example() {
	input := []byte(`Title
=

This is a new paragraph. I wonder    if I have too     many spaces.
What about new paragraph.
But the next one...

  Is really new.

1. Item one.
1. Item TWO.


Final paragraph.
`)

	output := blackfriday.Markdown(input, markdown.NewRenderer(), 0)

	os.Stdout.Write(output)

	// Output:
	//Title
	//=====
	//
	//This is a new paragraph. I wonder if I have too many spaces. What about new paragraph. But the next one...
	//
	//Is really new.
	//
	//1. Item one.
	//2. Item TWO.
	//
	//Final paragraph.
	//
}

func Example2() {
	input := []byte(`Title
==

Subtitle
---

How about ` + "`this`" + ` and other stuff like *italic*, **bold** and ***super    extra***.
`)

	output := blackfriday.Markdown(input, markdown.NewRenderer(), 0)

	os.Stdout.Write(output)

	// Output:
	//Title
	//=====
	//
	//Subtitle
	//--------
	//
	//How about `this` and other stuff like *italic*, **bold** and ***super extra***.
	//
}
