package main

import ()

func Reverse(s string) (string) {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	println(Reverse("Hello."))
	print("`", Reverse(""), "`\n")
	print("`", Reverse("1"), "`\n")
	print("`", Reverse("12"), "`\n")
	print("`", Reverse("123"), "`")
}

/*
Test cases are in three places:

1. main() in this file
2. Test() in main_test.go of this package:
	`go test gist.github.com/5423254.git`
3. Another package:
	./GoLand/src/gist.github.com/5423515.git/main.go

What's the best way?
*/