package gist5892738

import ()

// Trims the last newline character from str.
//
// Pre-condition: str ends with '\n'.
func MustTrimLastNewline(str string) string {
	return str[:len(str)-1]
}

// Trims the last newline character from str, if it ends with '\n'.
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
