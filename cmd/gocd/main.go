// Changes cwd to the src dir of Go package.
package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()

	args := []string{"list", "-f", "{{.Dir}}"}
	args = append(args, flag.Args()...)

	var buf = new(bytes.Buffer)

	cmd := exec.Command("go", args...)
	cmd.Stdout = buf
	_ = cmd.Run()

	lines := bytes.Split(buf.Bytes(), []byte("\n"))

	os.Stdout.Write(lines[0])

	// function gocd { cd `go list -f '{{.Dir}}' $1`; }
}
