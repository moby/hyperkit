package u8_test

import (
	"fmt"
	"time"

	"github.com/shurcooL/go/u/u8"
)

func Example1() {
	x := u8.AfterSecond(func() { fmt.Println("hi") })

	time.Sleep(500 * time.Millisecond)

	x.Cancel()

	time.Sleep(1500 * time.Millisecond)

	// Output:
}

func Example2() {
	x := u8.AfterSecond(func() { fmt.Println("hi") })

	time.Sleep(2 * time.Second)

	x.Cancel()

	// Output:
	//hi
}
