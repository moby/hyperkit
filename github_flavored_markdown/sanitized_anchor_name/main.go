// Package sanitized_anchor_name provides a func to create sanitized anchor names.
package sanitized_anchor_name

import "unicode"

// Create returns a sanitized anchor name for the given text.
func Create(text string) string {
	var anchorName []rune
	var futureDash = false
	for _, r := range []rune(text) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}
