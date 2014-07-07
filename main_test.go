package gist6418462

import "fmt"

var f2 = func() { panic(1337) }

func Example1() {
	f := func() {
		println("Hello from anon func!") // Comments are currently not preserved
	}
	if 5*5 > 26 {
		f = f2
	}

	fmt.Println(GetSourceAsString(f))

	// Output:
	//func() {
	//	println("Hello from anon func!")
	//}
}

func Example2() {
	f2 := func(a int, b int) int {
		c := a + b
		return c
	}

	fmt.Println(GetSourceAsString(f2))

	// Output:
	//func(a int, b int) int {
	//	c := a + b
	//	return c
	//}
}

func ExampleNil() {
	var f func()

	fmt.Println(GetSourceAsString(f))

	// Output:
	//nil
}
