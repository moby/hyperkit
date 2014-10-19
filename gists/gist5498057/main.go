package gist5498057

import (
	"io/ioutil"
	"os"
)

func ReadAllStdinB() []byte {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return b
}

func ReadAllStdin() string {
	return string(ReadAllStdinB())
}

func main() {
	print(ReadAllStdin())
}
