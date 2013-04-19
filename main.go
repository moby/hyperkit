package main

import ()

func Reverse(s string) string {
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