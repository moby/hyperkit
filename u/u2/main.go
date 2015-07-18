// Package u2 offers a TrackingWriter that tracks the last byte written on every write.
package u2

import (
	"bufio"
	"io"
)

// TrackingWriter tracks the last byte written on every write.
type TrackingWriter struct {
	w    *bufio.Writer
	last byte
}

func NewTrackingWriter(w io.Writer) *TrackingWriter {
	return &TrackingWriter{
		w: bufio.NewWriter(w),
	}
}

func (t *TrackingWriter) Write(p []byte) (n int, err error) {
	n, err = t.w.Write(p)
	if n > 0 {
		t.last = p[n-1]
	}
	return
}

func (t *TrackingWriter) Flush() {
	t.w.Flush()
}

func (t *TrackingWriter) EndsWithSpace() bool {
	return t.last == ' '
}
