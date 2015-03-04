// Package gist5408860 provides a reverse adapter for sort.
package gist5408860

import "sort"

// Reverse is a reverse sort.Interface adatper.
type Reverse struct {
	sort.Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
