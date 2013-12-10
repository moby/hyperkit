package gist6003701

import (
	"strings"
)

func UnderscoreSepToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

func main() {
	println(UnderscoreSepToCamelCase("g_string_append"))
}
