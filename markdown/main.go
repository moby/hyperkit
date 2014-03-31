// Package markdown provides a Markdown renderer.
package markdown

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/russross/blackfriday"
)

type markdownRenderer struct {
	normalTextMarker   map[*bytes.Buffer]int
	orderedListCounter map[int]int
	listDepth          int
}

func formatCode(lang string, text []byte) (formattedCode []byte, ok bool) {
	switch lang {
	case "Go", "go":
		gofmt, err := gofmt3b(string(text))
		if err != nil {
			return nil, false
		}
		return gofmt, true
	default:
		return nil, false
	}
}

// TODO: Unfinished implementation.
// Block-level callbacks.
func (_ *markdownRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	doubleSpace(out)

	// parse out the language name
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		out.WriteString("```")
		out.WriteString(elt)
		count++
		break
	}

	if count == 0 {
		out.WriteString("```")
	}
	out.WriteString("\n")

	if formattedCode, ok := formatCode(lang, text); ok {
		out.Write(formattedCode)
	} else {
		out.Write(text)
	}

	out.WriteString("```\n")
}
func (_ *markdownRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
	doubleSpace(out)
	lines := bytes.Split(text, []byte("\n"))
	for i, line := range lines {
		if i == len(lines)-1 {
			continue
		}
		out.WriteString(">")
		if len(line) != 0 {
			out.WriteString(" ")
			out.Write(line)
		}
		out.WriteString("\n")
	}
}
func (_ *markdownRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
	doubleSpace(out)
	out.Write(text)
	out.WriteByte('\n')
}
func (_ *markdownRenderer) Header(out *bytes.Buffer, text func() bool, level int) {
	marker := out.Len()
	doubleSpace(out)

	if level >= 3 {
		fmt.Fprint(out, strings.Repeat("#", level), " ")
	}

	textMarker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}

	switch level {
	case 1:
		fmt.Fprint(out, "\n", strings.Repeat("=", out.Len()-textMarker))
	case 2:
		fmt.Fprint(out, "\n", strings.Repeat("-", out.Len()-textMarker))
	}
	out.WriteString("\n")
}
func (_ *markdownRenderer) HRule(out *bytes.Buffer) {
	doubleSpace(out)
	out.WriteString("* * *")
}
func (m *markdownRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	doubleSpace(out)

	m.listDepth++
	defer func() { m.listDepth-- }()
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		m.orderedListCounter[m.listDepth] = 1
	}
	if !text() {
		out.Truncate(marker)
		return
	}
}
func (m *markdownRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	/*if flags&blackfriday.LIST_ITEM_CONTAINS_BLOCK != 0 {
		doubleSpace(out)
	}*/
	out.WriteString(strings.Repeat("\t", (m.listDepth - 1)))
	if flags&blackfriday.LIST_TYPE_ORDERED != 0 {
		fmt.Fprintf(out, "%d. %s", m.orderedListCounter[m.listDepth], string(text))
		m.orderedListCounter[m.listDepth]++
	} else {
		out.WriteString("- ")
		out.Write(text)
	}
	out.WriteString("\n")
}
func (_ *markdownRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}
func (_ *markdownRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.WriteString("<Table: Not implemented.>") // TODO
}
func (_ *markdownRenderer) TableRow(out *bytes.Buffer, text []byte) {
	out.WriteString("<TableRow: Not implemented.>") // TODO
}
func (_ *markdownRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("<TableHeaderCell: Not implemented.>") // TODO
}
func (_ *markdownRenderer) TableCell(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("<TableCell: Not implemented.>") // TODO
}
func (m *markdownRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	out.WriteString("<Footnotes: Not implemented.>") // TODO
}
func (_ *markdownRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	out.WriteString("<FootnoteItem: Not implemented.>") // TODO
}

// Span-level callbacks.
func (_ *markdownRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}
func (m *markdownRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	m.spaceIfNeeded(out)
	out.WriteByte('`')
	out.Write(text)
	out.WriteByte('`')
	m.normalTextMarker[out] = out.Len()
}
func (m *markdownRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	m.spaceIfNeeded(out)
	out.WriteString("**")
	out.Write(text)
	out.WriteString("**")
	m.normalTextMarker[out] = out.Len()
}
func (m *markdownRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	m.spaceIfNeeded(out)
	out.WriteByte('*')
	out.Write(text)
	out.WriteByte('*')
	m.normalTextMarker[out] = out.Len()
}
func (_ *markdownRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	out.WriteString("![")
	out.Write(alt)
	out.WriteString("](")
	out.Write(link)
	out.WriteString(")")
}
func (_ *markdownRenderer) LineBreak(out *bytes.Buffer) {
	out.WriteString("<br>")
}
func (m *markdownRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	m.spaceIfNeeded(out)
	out.WriteString("[")
	out.Write(content)
	out.WriteString("](")
	out.Write(link)
	out.WriteString(")")
	m.normalTextMarker[out] = out.Len()
}
func (_ *markdownRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	out.Write(tag)
}
func (m *markdownRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	m.spaceIfNeeded(out)
	out.WriteString("***")
	out.Write(text)
	out.WriteString("***")
	m.normalTextMarker[out] = out.Len()
}
func (_ *markdownRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString("~~")
	out.Write(text)
	out.WriteString("~~")
}
func (_ *markdownRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	out.WriteString("<FootnoteRef: Not implemented.>") // TODO
}

// Low-level callbacks
func (_ *markdownRenderer) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}
func (m *markdownRenderer) NormalText(out *bytes.Buffer, text []byte) {
	cleanString := clean(string(text))
	if cleanString == "" {
		return
	}
	m.spaceIfNeededNormalText(out, cleanString)
	out.WriteString(cleanString)
	m.normalTextMarker[out] = out.Len()
}

// Header and footer.
func (_ *markdownRenderer) DocumentHeader(out *bytes.Buffer) {}
func (_ *markdownRenderer) DocumentFooter(out *bytes.Buffer) {}

func (_ *markdownRenderer) GetFlags() int { return 0 }

func (m *markdownRenderer) spaceIfNeeded(out *bytes.Buffer) {
	if _, ok := m.normalTextMarker[out]; !ok {
		m.normalTextMarker[out] = -1
	}
	if m.normalTextMarker[out] == out.Len() {
		out.WriteByte(' ')
	}
}

func (m *markdownRenderer) spaceIfNeededNormalText(out *bytes.Buffer, cleanString string) {
	if _, ok := m.normalTextMarker[out]; !ok {
		m.normalTextMarker[out] = -1
	}
	if m.normalTextMarker[out] == out.Len() && !isPunctuation(cleanString[0]) {
		out.WriteByte(' ')
	}
}

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

func isPunctuation(b byte) bool {
	switch b {
	case ',', '.', ':', ';', '_':
		return true
	default:
		return false
	}
}

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

// TODO: Replace with go1.1's go/format
// Actually executes gofmt binary as a new process
// TODO: Can't use it until go/format is fixed to be consistent with gofmt, currently it strips comments out of partial Go programs
// See: https://code.google.com/p/go/issues/detail?id=5551
func gofmt3b(str string) ([]byte, error) {
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "gofmt"))

	// TODO: Error checking and other niceness
	// http://stackoverflow.com/questions/13432947/exec-external-program-script-and-detect-if-it-requests-user-input
	in, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		_, err = in.Write([]byte(str))
		if err != nil {
			panic(err)
		}
		err = in.Close()
		if err != nil {
			panic(err)
		}
	}()

	data, err := cmd.Output()
	if nil != err {
		return []byte("gofmt error!\n" + str), err
	}
	return data, nil
}

// NewRenderer returns a Markdown renderer.
func NewRenderer() blackfriday.Renderer {
	return &markdownRenderer{
		normalTextMarker:   make(map[*bytes.Buffer]int),
		orderedListCounter: make(map[int]int),
	}
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

	// GitHub Flavored Markdown-like extensions.
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	//extensions |= blackfriday.EXTENSION_TABLES // TODO: Implement.
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	output := blackfriday.Markdown(text, NewRenderer(), extensions)
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
