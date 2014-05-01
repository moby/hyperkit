package gist6003701

import (
	"strings"
	"unicode"
)

// Converts "string_URL_append" to "StringUrlAppend" form.
func UnderscoreSepToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

func addSegment(inout, seg []rune) []rune {
	if len(seg) == 0 {
		return inout
	}
	if len(inout) != 0 {
		inout = append(inout, '_')
	}
	inout = append(inout, seg...)
	return inout
}

// Converts "StringUrlAppend" to "string_url_append" form.
func CamelCaseToUnderscoreSep(s string) string {
	var out []rune
	var seg []rune
	for _, r := range s {
		if !unicode.IsLower(r) {
			out = addSegment(out, seg)
			seg = nil
		}
		seg = append(seg, unicode.ToLower(r))
	}
	out = addSegment(out, seg)
	return string(out)
}

func main() {
	in := "string_URL_append"
	out := UnderscoreSepToCamelCase(in)
	println(in, "->", UnderscoreSepToCamelCase(in))
	println(out, "->", CamelCaseToUnderscoreSep(out))
}
