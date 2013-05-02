package main

import (
	"io/ioutil"
	"os"
	. "gist.github.com/5286084.git"
)

func ReadAllStdinB() []byte {
	b, err := ioutil.ReadAll(os.Stdin)
	CheckError(err)

	return b
}

func ReadAllStdin() string {
	return string(ReadAllStdinB())
}

func main() {
	print(ReadAllStdin())
}