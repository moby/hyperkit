package gist5892738

import (
)

// TODO: Use runes, unicode rather than bytes
func TrimLastNewline(str string) string {
	if len(str) > 0 && '\n' == str[len(str)-1] {
		return str[:len(str)-1]
	}
	return str
}

func main() {
	println("'" + TrimLastNewline("String\n") + "'")
	println("'" + TrimLastNewline("String") + "'")
	println("'" + TrimLastNewline("") + "'")
	println("'" + TrimLastNewline("\n") + "'")
}