package gist6433744

import (
	"bytes"
)

// Gets the starting and ending caret indicies of line with specified lineIndex.
// Does not include newline character.
// First line has index 0.
// Returns (-1, -1) if line is not found.
func GetLineStartEndIndicies(b []byte, lineIndex int) (startIndex, endIndex int) {
	index := 0
	for line := 0; ; line++ {
		lineLength := bytes.IndexByte(b[index:], '\n')
		if line == lineIndex {
			if lineLength == -1 {
				return index, len(b)
			} else {
				return index, index + lineLength
			}
		}
		if lineLength == -1 {
			break
		}
		index += lineLength + 1
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
