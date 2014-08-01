package gist5498057

import (
	"io/ioutil"
	"os"

	. "github.com/shurcooL/go/gists/gist5286084"
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
