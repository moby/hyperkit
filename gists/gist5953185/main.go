// Package gist5953185 offers a func to underline a header string.
package gist5953185

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// Underline returns an underlined s.
func Underline(s string) string {
	return s + "\n" + strings.Repeat("-", runewidth.StringWidth(s)) + "\n"
}
