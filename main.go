package gist5892738

// Trims the last newline character from str.
//
// Pre-condition: str ends with '\n'.
func MustTrimLastNewline(str string) string {
	return str[:len(str)-1]
}

// Trims the last newline character from str, if any.
func TrimLastNewline(str string) string {
	if len(str) >= 1 && '\n' == str[len(str)-1] {
		return MustTrimLastNewline(str)
	}
	return str
}

// Trims the first space character from str.
//
// Pre-condition: str starts with ' '.
func MustTrimFirstSpace(str string) string {
	return str[1:]
}

// Trims the first space character from str, if any.
func TrimFirstSpace(str string) string {
	if len(str) >= 1 && ' ' == str[0] {
		return MustTrimFirstSpace(str)
	}
	return str
}

func main() {
	println("'" + TrimLastNewline("String\n") + "'")
	println("'" + TrimLastNewline("String") + "'")
	println("'" + TrimLastNewline("") + "'")
	println("'" + TrimLastNewline("\n") + "'")
}
