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

	x.Wait() // Don't exit before x finishes its work to prevent a race.

	// Output:
}

func Example2() {
	x := u8.AfterSecond(func() { fmt.Println("hi") })

	time.Sleep(2 * time.Second)

	x.Cancel()

	x.Wait() // Don't exit before x finishes its work to prevent a race.

	// Output:
	//hi
}
