package gist5092053

import "sort"

// reverseAdapter is a reverse sort.Interface adapter.
type reverseAdapter struct {
	sort.Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r reverseAdapter) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
