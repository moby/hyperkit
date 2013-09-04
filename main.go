package main

import (
	"bytes"
)

// Gets the starting and ending caret indicies of line with specified lineIndex.
// Does not include newline character.
// First line has index 0.
// Returns (-1, -1) if line is not found.
func GetLineStartEndIndicies(b []byte, lineIndex int) (startIndex, endIndex int) {
	n := 0
	line := 0
	for {
		o := bytes.IndexByte(b[n:], '\n')
		if line == lineIndex {
			if o == -1 {
				return n, len(b)
			} else {
				return n, n + o
			}
		}
		if o == -1 {
			break
		}
		n += o + 1
		line++
	}

	return -1, -1
}

func main() {
	b := []byte(`this

this is a longer line
and
stuff
last`)

	for lineIndex := 0; ; lineIndex++ {
		s, e := GetLineStartEndIndicies(b, lineIndex)
		print(lineIndex, ": [", s, ", ", e, "]\n")
		if s == -1 {
			break
		}

	}
}
