package gist5423254

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	if got, want := Reverse("Hello."), ".olleH"; got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

func Example() {
	fmt.Println(Reverse("Hello."))
	fmt.Printf("%q\n", Reverse(""))
	fmt.Printf("%q\n", Reverse("1"))
	fmt.Printf("%q\n", Reverse("12"))
	fmt.Printf("%q\n", Reverse("123"))
	fmt.Printf("%q\n", Reverse("Hello, 世界"))

	// Output:
	// .olleH
	// ""
	// "1"
	// "21"
	// "321"
	// "界世 ,olleH"
}
