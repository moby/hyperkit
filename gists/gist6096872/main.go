// Package gist6096872 offers interchange between io.Writer/io.Reader types and channels of []byte.
package gist6096872

import (
	"bufio"
	"io"
)

type ChanWriter chan []byte

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	// TODO: Copy the slice contents rather than sending the original, as it may get modified after we return?
	cw <- p
	return len(p), nil
}

func LineReader(r io.Reader) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadBytes('\n')
			if err != nil {
				ch <- line
				close(ch)
				return
			}
			ch <- line[:len(line)-1] // Trim last newline.
		}
	}()
	return ch
}
