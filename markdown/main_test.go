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

Note that --- not considering the asterisk --- the actual text
content starts at 4-columns in.

> Block quotes are
> written like so.
>
> They can span multiple paragraphs,
> if you like.

Use 3 dashes for an em-dash. Use 2 dashes for ranges (ex. "it's all in
chapters 12--14"). Three dots ... will be converted to an ellipsis.

An h2 header
------------

Here's a numbered list:

1. first item
2. second item
3. third item

Note again how the actual text starts at 4 columns in (4 characters from the left side). Here's a code sample:

	# Let me re-iterate ...
	for i in 1 .. 10 { do-something(i) }

As you probably guessed, indented 4 spaces. By the way, instead of indenting the block, you can use delimited blocks, if you like:

~~~
define foobar() {
	print "Welcome to flavor country!";
}
~~~

(which makes copying & pasting easier). You can optionally mark the
delimited block for Pandoc to syntax highlight it:

~~~python
import time
# Quick, count to ten!
for i in range(10):
	# (but not *too* quick)
	time.sleep(0.5)
	print i
~~~

### An h3 header

Now a nested list:

1. First, get these ingredients:

      * carrots
      * celery
      * lentils

2. Boil some water.

3. Dump everything in the pot and follow this algorithm:

        find wooden spoon
        uncover pot
        stir
        cover pot
        balance wooden spoon precariously on pot handle
        wait 10 minutes
        goto first step (or shut off burner when done)

    Do not bump wooden spoon or it will fall.

Notice again how text always lines up on 4-space indents (including
that last line which continues item 3 above). Here's a link to [a
website](http://foo.bar). Here's a link to a [local
doc](local-doc.html). Here's a footnote [^1].

[^1]: Footnote text goes here.

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
