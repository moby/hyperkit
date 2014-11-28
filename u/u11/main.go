package u11

import (
	"io"
	"os"
)

// WriteFile copies from r to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions 0644;
// otherwise WriteFile truncates it before writing.
func WriteFile(r io.Reader, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
