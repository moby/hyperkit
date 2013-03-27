package main

import (
	"strings"
)

func GetLines(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return r == '\n' })
}

func GetLine(s string, n int) string {
	return GetLines(s)[n]
}

func main() {
	str := "First Line,\n2nd Line.\nThird!"

	print(GetLine(str, 1))
}