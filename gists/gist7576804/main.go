package gist7576804

import (
	"code.google.com/p/go.tools/go/types"
	"fmt"
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
