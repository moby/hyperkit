// Package indentwriter implements an io.Writer wrapper that indents
// every non-empty line with specified number of tabs.
package indentwriter

import (
	"io"

	"github.com/bradfitz/iter"
)

type indentWriter struct {
	w      io.Writer
	indent int

	wroteIndent bool
}

func New(w io.Writer, indent int) *indentWriter {
	return &indentWriter{w: w, indent: indent}
}

func (iw *indentWriter) Write(p []byte) (n int, err error) {
	//strings.Repeat("\t", mr.listDepth)
	//return iw.w.Write(bytes.Replace(p, []byte("\n"), []byte("\n\t"), -1))
	for _, b := range p {
		err = iw.WriteByte(b)
		if err != nil {
			return
		}
		n++
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}
	return len(p), nil
}

func (iw *indentWriter) WriteString(s string) (n int, err error) {
	return iw.Write([]byte(s))
}

func (iw *indentWriter) WriteByte(c byte) error {
	//return iw.Write([]byte{b})

	if c == '\n' {
		iw.wroteIndent = false
	} else {
		if !iw.wroteIndent {
			iw.wroteIndent = true
			for _, _ = range iter.N(iw.indent) {
				iw.w.Write([]byte{'\t'})
			}
		}
	}
	_, err := iw.w.Write([]byte{c})
	return err
}
