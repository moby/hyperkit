// Package trim contains helpers for trimming strings.
package trim

// LastNewline trims the last newline character from str, if any.
func LastNewline(str string) string {
	if len(str) >= 1 && str[len(str)-1] == '\n' {
		return str[:len(str)-1]
	} else {
		return str
	}
}

// FirstSpace trims the first space character from str, if any.
func FirstSpace(str string) string {
	if len(str) >= 1 && str[0] == ' ' {
		return str[1:]
	} else {
		return str
	}
}
