package gist5953185

import (
	"strings"
)

// TODO: Use runes, unicode for getting string length
func Underline(s string) string {
	return s + "\n" + strings.Repeat("-", len(s)) + "\n"
}

func main() {
	println(Underline("Underline Test") + "stuff that goes here")
}