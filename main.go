package main

import (
)

// Panics on error
func CheckError(err error) {
	if nil != err {
		panic(err)
	}
}

func main() {
}