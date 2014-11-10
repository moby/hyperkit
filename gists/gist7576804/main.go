package gist7576804

import (
	"fmt"

	"golang.org/x/tools/go/types"
)

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
