// Package sanitized_anchor_name provides a func to create sanitized anchor names.
package sanitized_anchor_name

import "unicode"

// Create returns a sanitized anchor name for the given text.
func Create(text string) string {
	var anchorName []rune
	for _, r := range []rune(text) {
		switch {
		case r == ' ':
			anchorName = append(anchorName, '-')
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			anchorName = append(anchorName, unicode.ToLower(r))
		}
	}
	return string(anchorName)
}
