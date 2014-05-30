package u7

import (
	"bufio"
	"bytes"
	"io"
	"text/template"
)

func Print(s *Scanner, w io.Writer) error {
	var p Printer = HTMLPrinter(HTMLConfig{
		"",
		"gi",
		"gd",
		"gu",
		"gh",
	})

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

type Printer interface {
	Print(w io.Writer, tok []byte, kind int) error
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
	s := &Scanner{br: bufio.NewReader(r)}
	return s
}

func (s *Scanner) Scan() bool {
	var err error
	s.line, err = s.br.ReadBytes('\n')
	return err == nil
}

func (s *Scanner) Token() ([]byte, int) {
	var kind int
	switch {
	case len(s.line) == 0 ||
		s.line[0] == ' ':

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
