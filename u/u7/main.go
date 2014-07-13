package u7

import (
	"bufio"
	"bytes"
	"io"
	"text/template"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sourcegraph/annotate"
	"github.com/sourcegraph/syntaxhighlight"
)

var gfmDiff = HTMLConfig{
	"",
	"gi",
	"gd",
	"gu",
	"gh",
}

func Print(s *Scanner, w io.Writer) error {
	var p syntaxhighlight.Printer = HTMLPrinter(gfmDiff)

	for s.Scan() {
		tok, kind := s.Token()
		err := p.Print(w, tok, kind)
		if err != nil {
			return err
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

type HTMLConfig []string

type HTMLPrinter HTMLConfig

func (p HTMLPrinter) Print(w io.Writer, tok []byte, kind int) error {
	class := HTMLConfig(p)[kind]
	if class != "" {
		_, err := w.Write([]byte(`<span class="`))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, class)
		if err != nil {
			return err
		}
		io.WriteString(w, " input-block") // For "display: block;" style.
		_, err = w.Write([]byte(`">`))
		if err != nil {
			return err
		}
	}
	template.HTMLEscape(w, tok)
	if class != "" {
		_, err := w.Write([]byte(`</span>`))
		if err != nil {
			return err
		}
	}
	return nil
}

type Scanner struct {
	br   *bufio.Reader
	line []byte
}

func NewScanner(src []byte) *Scanner {
	r := bytes.NewReader(src)
	return &Scanner{br: bufio.NewReader(r)}
}

func (s *Scanner) Scan() bool {
	var err error
	s.line, err = s.br.ReadBytes('\n')
	return err == nil
}

func (s *Scanner) Token() ([]byte, int) {
	var kind int
	switch {
	case len(s.line) == 0 || s.line[0] == ' ':
		kind = 0
	case s.line[0] == '+':
		kind = 1
	case s.line[0] == '-':
		kind = 2
	case s.line[0] == '@':
		kind = 3
	default:
		kind = 4
	}
	return s.line, kind
}

func (s *Scanner) Err() error {
	return nil
}

// ---

type HTMLAnnotator HTMLConfig

func (a HTMLAnnotator) Annotate(start int, tok []byte, kind int) (*annotate.Annotation, error) {
	class := HTMLConfig(a)[kind]
	if class != "" {
		left := []byte(`<span class="`)
		left = append(left, []byte(class)...)
		left = append(left, []byte(" input-block")...) // For "display: block;" style.
		left = append(left, []byte(`">`)...)
		return &annotate.Annotation{
			Start: start, End: start + len(tok),
			Left: left, Right: []byte("</span>"),
		}, nil
	}
	return nil, nil
}

func Annotate(src []byte) (annotate.Annotations, error) {
	var a syntaxhighlight.Annotator = HTMLAnnotator(gfmDiff)
	s := NewScanner(src)

	var anns annotate.Annotations
	read := 0
	for s.Scan() {
		tok, kind := s.Token()
		ann, err := a.Annotate(read, tok, kind)
		if err != nil {
			return nil, err
		}
		read += len(tok)
		if ann != nil {
			anns = append(anns, ann)
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return anns, nil
}

// ---

func HighlightedDiffFunc(leftContent, rightContent string, segments *[2][]*annotate.Annotation, offsets [2]int) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(leftContent, rightContent, true)

	for side := range *segments {
		offset := offsets[side]

		for _, diff := range diffs {
			if side == 0 && diff.Type == -1 {
				(*segments)[side] = append((*segments)[side], &annotate.Annotation{Start: offset, End: offset + len(diff.Text), Left: []byte(`<span class="x">`), Right: []byte(`</span>`), WantInner: 1})
				offset += len(diff.Text)
			}
			if side == 1 && diff.Type == +1 {
				(*segments)[side] = append((*segments)[side], &annotate.Annotation{Start: offset, End: offset + len(diff.Text), Left: []byte(`<span class="x">`), Right: []byte(`</span>`), WantInner: 1})
				offset += len(diff.Text)
			}
			if diff.Type == 0 {
				offset += len(diff.Text)
			}
		}
	}
}

/*func highlightedDiffFunc2(leftContent, rightContent string) (leftOut, rightOut string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(leftContent, rightContent, true)

	for _, diff := range diffs {
		switch diff.Type {
		case -1:
			leftOut += `<span class="x">` + diff.Text + `</span>`
		case +1:
			rightOut += `<span class="x">` + diff.Text + `</span>`
		case 0:
			leftOut += diff.Text
			rightOut += diff.Text
		}
	}

	return
}*/
