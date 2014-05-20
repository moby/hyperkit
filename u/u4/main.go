package u4

import (
	"log"
	"os/exec"
	"runtime"
)

// Open opens a file (or a directory or Url), just as if the user had double-clicked the file's icon.
// It uses the default application, as determined by the OS.
func Open(path string) {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	args = append(args, path)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Println("u4.Open:", err)
	}
}
