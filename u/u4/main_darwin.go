package u4

import (
	"log"
	"os/exec"
)

// Open opens a file (or a directory or Url), just as if the user had double-clicked the file's icon.
// It uses the default application, as determined by the OS.
func Open(path string) {
	cmd := exec.Command("open", path)
	err := cmd.Run()
	if err != nil {
		log.Println("u4.Open:", err)
	}
}
