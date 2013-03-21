package main

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