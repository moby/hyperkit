package markdown_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

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

	output, err := markdown.Process("", input, nil)
	if err != nil {
		log.Fatalln(err)
	}

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

	output, err := markdown.Process("", input, nil)
	if err != nil {
		log.Fatalln(err)
	}

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

const reference = `An h1 header
============

Paragraphs are separated by a blank line.

2nd paragraph. *Italic*, **bold**, ` + "`monospace`" + `. Itemized lists look like:

- this one
- that one
- the other one

Nothing to note here.

> Block quotes are written like so.
>
> > They can be nested.
>
> They can span multiple paragraphs, if you like.

- Item 1
- Item 2
	- Item 2a
		- Item 2a
	- Item 2b
- Item 3

Hmm.

1. Item 1
2. Item 2
	1. Blah.
	2. Blah.
3. Item 3
	- Item 3a
	- Item 3b

Last paragraph here.

An h2 header
------------

[Visit GitHub!](www.github.com)

![Hmm](http://example.org/image.png)

~~Mistaken text.~~

http://example.com

### An h3 header

Here's a numbered list:

1. first item
2. second item
3. third item

Note again how the actual text starts at 4 columns in (4 characters from the left side). Here's a code sample:

` + "```" + `
# Let me re-iterate ...
for i in 1 .. 10 { do-something(i) }
` + "```" + `

As you probably guessed, indented 4 spaces. By the way, instead of indenting the block, you can use delimited blocks, if you like:

` + "```" + `
define foobar() {
    print "Welcome to flavor country!";
}
` + "```" + `

(which makes copying & pasting easier). You can optionally mark the delimited block for Pandoc to syntax highlight it:

` + "```" + `Go
func main() {
    println("Hi.")
}
` + "```" + `

Done.
`

func Test(t *testing.T) {

	output, err := markdown.Process("", []byte(reference), nil)
	if err != nil {
		log.Fatalln(err)
	}

	diff, err := diff(output, []byte(reference))
	if err != nil {
		log.Fatalln(err)
	}

	if len(diff) != 0 {
		t.Errorf("Difference of %d lines:\n%s", bytes.Count(diff, []byte("\n")), string(diff))
	}
}

// TODO: Factor out.
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "mdfmt")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "mdfmt")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
