// Changes cwd to the src dir of Go package.
package main

import (
	"flag"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()

	args := []string{"list", "-f", "{{.Dir}}"}
	args = append(args, flag.Args()...)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	// TODO: Finish.
}
