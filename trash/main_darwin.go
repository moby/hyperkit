package trash

import (
	"os"
	"path/filepath"
)

// TODO: Use better type than string.
func MoveToTrash(name string) error {
	name = filepath.Clean(name)
	home := os.Getenv("HOME")
	_, file := filepath.Split(name)
	target := filepath.Join(home, ".Trash", file)

	// TODO: If target name exists in Trash, come up with a unique one (perhaps append a timestamp) instead of overwriting.
	// TODO: Support OS X "Put Back". Figure out how it's done and do it.

	return os.Rename(name, target)
}
