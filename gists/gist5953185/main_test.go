package gist5953185

import "fmt"

func Example() {
	fmt.Print(Underline("Underline Test") + "\nstuff that goes here")

	// Output:
	//Underline Test
	//--------------
	//
	//stuff that goes here
}
