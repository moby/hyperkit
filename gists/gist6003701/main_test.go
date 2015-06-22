package gist6003701

import "fmt"

func ExampleUnderscoreSepToCamelCase() {
	fmt.Println(UnderscoreSepToCamelCase("string_URL_append"))

	// Output: StringUrlAppend
}

func ExampleCamelCaseToUnderscoreSep() {
	fmt.Println(CamelCaseToUnderscoreSep("StringUrlAppend"))

	// Output: string_URL_append
}

func ExampleUnderscoreSepToMixedCaps() {
	fmt.Println(UnderscoreSepToMixedCaps("string_URL_append"))

	// Output: StringURLAppend
}

/*func ExampleMixedCapsToUnderscoreSep() {
	fmt.Println(MixedCapsToUnderscoreSep("StringURLAppend"))

	// Output: string_URL_append
}*/
