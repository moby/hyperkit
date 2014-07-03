// Pretty-prints JSON from stdin.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, in, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	out.WriteTo(os.Stdout)
}
