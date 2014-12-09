// Package trim contains helpers for trimming strings.
package trim

func willTrimLastNewline(str string) string {
	return str[:len(str)-1]
}

// LastNewline trims the last newline character from str, if any.
func LastNewline(str string) string {
	if len(str) >= 1 && '\n' == str[len(str)-1] {
		return willTrimLastNewline(str)
	} else {
		return str
	}
}

func willTrimFirstSpace(str string) string {
	return str[1:]
}

// FirstSpace trims the first space character from str, if any.
func FirstSpace(str string) string {
	if len(str) >= 1 && ' ' == str[0] {
		return willTrimFirstSpace(str)
	} else {
		return str
	}
}
