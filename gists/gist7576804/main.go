// Package gist7576804 implements TypeChainString to get a full type chain string.
package gist7576804

import (
	"fmt"

	"golang.org/x/tools/go/types"
)

// TypeChainString returns the full type chain as a string.
func TypeChainString(t types.Type) string {
	out := fmt.Sprintf("%s", t)
	for {
		if t == t.Underlying() {
			break
		} else {
			t = t.Underlying()
		}
		out += fmt.Sprintf(" -> %s", t)
	}
	return out
}
