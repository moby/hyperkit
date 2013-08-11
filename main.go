package gist6096872

import (
	"bufio"
	"io"
)

// Credit to Tarmigan
func ByteReader(r io.Reader) (<-chan []byte, <-chan error) {
	ch := make(chan []byte)
	errCh := make(chan error)

	go func() {
		for {
			buf := make([]byte, 2048)
			s := 0
		inner:
			for {
				n, err := r.Read(buf[s:])
				if n > 0 {
					ch <- buf[s : s+n]
					s += n
				}
				if err != nil {
					errCh <- err
					return
				}
				if s >= len(buf) {
					break inner
				}
			}
		}
	}()

	return ch, errCh
}

func LineReader(r io.Reader) (<-chan []byte, <-chan error) {
	ch := make(chan []byte)
	errCh := make(chan error)

	go func() {
		br := bufio.NewReader(r)
		for {
			line, err := br.ReadBytes('\n')
			if err == nil {
				ch <- line[:len(line)-1] // Trim newline
			} else {
				ch <- line
				errCh <- err
				return
			}
		}
	}()

	return ch, errCh
}
