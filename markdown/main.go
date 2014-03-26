// Package markdown provides a Markdown renderer.
package markdown

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/russross/blackfriday"
)

type markdownRenderer struct {
	normalTextMarker   int
	orderedListCounter int
}

// TODO: Unfinished implementation.
// Block-level callbacks.
func (_ *markdownRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {}
func (_ *markdownRenderer) BlockQuote(out *bytes.Buffer, text []byte)             {}
func (_ *markdownRenderer) BlockHtml(out *bytes.Buffer, text []byte)              {}
func (_ *markdownRenderer) Header(out *bytes.Buffer, text func() bool, level int) {
	marker := out.Len()
	text()
	switch level {
	case 1:
		fmt.Fprint(out, "\n", strings.Repeat("=", out.Len()-marker), "\n")
	}
}
func (_ *markdownRenderer) HRule(out *bytes.Buffer) {}
func (m *markdownRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	m.orderedListCounter = 1
	text()
	fmt.Fprint(out, "\n")
}
func (m *markdownRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		fmt.Fprintf(out, "\n%d. %s", m.orderedListCounter, string(text))
		m.orderedListCounter++
	} else {
	}
}
func (_ *markdownRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	fmt.Fprint(out, "\n")
	text()
	fmt.Fprint(out, "\n")
}
func (_ *markdownRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (_ *markdownRenderer) TableRow(out *bytes.Buffer, text []byte)                               {}
func (_ *markdownRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             {}
func (_ *markdownRenderer) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (_ *markdownRenderer) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (_ *markdownRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}

// Span-level callbacks.
func (_ *markdownRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int)                 {}
func (_ *markdownRenderer) CodeSpan(out *bytes.Buffer, text []byte)                           {}
func (_ *markdownRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (_ *markdownRenderer) Emphasis(out *bytes.Buffer, text []byte)                           {}
func (_ *markdownRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)    {}
func (_ *markdownRenderer) LineBreak(out *bytes.Buffer)                                       {}
func (_ *markdownRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {}
func (_ *markdownRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte)                          {}
func (_ *markdownRenderer) TripleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (_ *markdownRenderer) StrikeThrough(out *bytes.Buffer, text []byte)                      {}
func (_ *markdownRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                 {}

// Low-level callbacks
func (_ *markdownRenderer) Entity(out *bytes.Buffer, entity []byte) { /*out.Write(entity)*/
}
func (m *markdownRenderer) NormalText(out *bytes.Buffer, text []byte) {
	cleanString := clean(string(text))
	if cleanString == "" {
		return
	}
	if m.normalTextMarker == out.Len() {
		out.WriteByte(' ')
	}
	out.WriteString(cleanString)
	m.normalTextMarker = out.Len()
}

// Header and footer.
func (_ *markdownRenderer) DocumentHeader(out *bytes.Buffer) {}
func (_ *markdownRenderer) DocumentFooter(out *bytes.Buffer) {}

// clean replaces each sequence of space, \n, \r, or \t characters
// with a single space and removes any trailing and leading spaces.
func clean(s string) string {
	var b []byte
	p := byte(' ')
	for i := 0; i < len(s); i++ {
		q := s[i]
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}
		if q != ' ' || p != ' ' {
			b = append(b, q)
			p = q
		}
	}
	// Remove trailing blank, if any.
	if n := len(b); n > 0 && p == ' ' {
		b = b[0 : n-1]
	}
	return string(b)
}

// NewRenderer returns a Markdown renderer.
func NewRenderer() blackfriday.Renderer {
	return &markdownRenderer{normalTextMarker: -1}
}

// Options specifies options for formatting.
type Options struct {
	// Currently none.
}

// Process formats Markdown.
// If opt is nil the defaults are used.
func Process(filename string, src []byte, opt *Options) ([]byte, error) {
	// Get source.
	text, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	output := blackfriday.Markdown(text, NewRenderer(), 0)
	return output, nil
}

// If src != nil, readSource returns src.
// If src == nil, readSource returns the result of reading the file specified by filename.
func readSource(filename string, src []byte) ([]byte, error) {
	if src != nil {
		return src, nil
	}
	return ioutil.ReadFile(filename)
}
