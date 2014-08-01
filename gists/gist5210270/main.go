// Package gist5210270 trims double quotes from a string.
package gist5210270

import (
	"strings"
)

func TrimQuotes(str string) string {
	return strings.Trim(str, "\"")
}

func main() {
	str := `"Quoted text becomes..."`
	print(str, " -> ", TrimQuotes(str))
}
