package github_flavored_markdown

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// In this test, nothing should be sanitized away.
func TestSanitize1(t *testing.T) {
	var text = []byte(`### GitHub Flavored Markdown rendered locally using go gettable native Go code

` + "```Go" + `
package main

import "fmt"

func main() {
	// This is a comment!
	/* so is this */
	fmt.Println("Hello, playground", 123, 1.336)
}
` + "```" + `

` + "```diff" + `
diff --git a/main.go b/main.go
index dc83bf7..5260a7d 100644
--- a/main.go
+++ b/main.go
@@ -1323,10 +1323,10 @@ func (this *GoPackageSelecterAdapter) GetSelectedGoPackage() *GoPackage {
 }

 // TODO: Move to the right place.
-var goPackages = &exp14.GoPackages{SkipGoroot: false}
+var goPackages = &exp14.GoPackages{SkipGoroot: true}

 func NewGoPackageListingWidget(pos, size mathgl.Vec2d) *SearchableListWidget {
 	goPackagesSliceStringer := &goPackagesSliceStringer{goPackages}
` + "```" + `
`)

	htmlFlags := 0
	//htmlFlags |= blackfriday.HTML_SANITIZE_OUTPUT
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	renderer := &renderer{blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}

	// Parser extensions for GitHub Flavored Markdown.
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	//extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK

	unsanitized := blackfriday.Markdown(text, renderer, extensions)

	// GitHub Flavored Markdown sanitization.
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")

	output := []byte(p.Sanitize(string(unsanitized)))

	output = append(output, '\n')
	diff, err := diff(unsanitized, output)
	if err != nil {
		log.Fatalln(err)
	}

	if len(diff) != 0 {
		t.Errorf("Difference of %d lines:\n%s", bytes.Count(diff, []byte("\n")), string(diff))
	}
}

// Make sure that <script> tag is sanitized away.
func TestSanitize2(t *testing.T) {
	text := []byte("Hello <script>alert();</script> world.")

	if expected, got := "<p>Hello  world.</p>", string(Markdown(text)); expected != got {
		t.Errorf("expected: %q, got: %q\n", expected, got)
	}
}

// Make sure that "class" attribute values that are not sane get sanitized away.
func TestSanitize3a(t *testing.T) {
	// Just a normal class name, should be preserved.
	text := []byte(`Hello <span class="foo bar bash">there</span> world.`)

	if expected, got := `<p>Hello <span class="foo bar bash">there</span> world.</p>`, string(Markdown(text)); expected != got {
		t.Errorf("expected: %q, got: %q\n", expected, got)
	}
}
func TestSanitize3b(t *testing.T) {
	// JavaScript in class name, should be sanitized away.
	text := []byte(`Hello <span class="javascript:alert('XSS')">there</span> world.`)

	if expected, got := "<p>Hello <span>there</span> world.</p>", string(Markdown(text)); expected != got {
		t.Errorf("expected: %q, got: %q\n", expected, got)
	}
}
func TestSanitize3c(t *testing.T) {
	// Script injection attempt, should be sanitized away.
	text := []byte(`Hello <span class="><script src='http://hackers.org/XSS.js'></script>">there</span> world.`)

	if expected, got := "<p>Hello ", string(Markdown(text)); expected != got {
		t.Errorf("expected: %q, got: %q\n", expected, got)
	}
}

// TODO: Factor out.
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "")
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
