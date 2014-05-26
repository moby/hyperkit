package gist5423254

import (
	"testing"

	"fmt"
)

func Test(t *testing.T) {
	if ".olleH" != Reverse("Hello.") {
		t.Fail()
	}
}

func Example() {
	fmt.Println(Reverse("Hello."))
	fmt.Print("`", Reverse(""), "`\n")
	fmt.Print("`", Reverse("1"), "`\n")
	fmt.Print("`", Reverse("12"), "`\n")
	fmt.Print("`", Reverse("123"), "`")

	// Output:
	//.olleH
	//``
	//`1`
	//`21`
	//`321`
}
