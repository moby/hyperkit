package gist6003701

import "fmt"

func ExampleUnderscoreSepToCamelCase() {
	fmt.Println(UnderscoreSepToCamelCase("string_URL_append"))

	// Output:
	//StringUrlAppend
}

func ExampleCamelCaseToUnderscoreSep() {
	fmt.Println(CamelCaseToUnderscoreSep("StringUrlAppend"))

	// Output:
	//string_url_append
}
