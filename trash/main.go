// +build !darwin

package trash

import "errors"

// TODO: Use better type than string.
// MoveToTrash moves name to trash.
func MoveToTrash(name string) error {
	return errors.New("MoveToTrash: not yet implemented on non-darwin")
}
