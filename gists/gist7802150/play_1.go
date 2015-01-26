// +build ignore

package main

import (
	. "github.com/shurcooL/go/gists/gist7802150"

	"fmt"

	"github.com/shurcooL/go/gists/gist6418290"
)

type file struct {
	ViewGroup
}

func (*file) SetSelf(s string) {
	fmt.Println(gist6418290.GetParentFuncAsString())
}

type websocket struct {
	ViewGroup
}

func (*websocket) SetSelf(s string) {
	fmt.Println(gist6418290.GetParentFuncAsString())
}

type memory struct {
	ViewGroup
}

func (*memory) SetSelf(s string) {
	fmt.Println(gist6418290.GetParentFuncAsString())
}

func main() {
	f := file{}
	f.InitViewGroup(&f, "memory://???")
	w := websocket{}
	w.InitViewGroup(&w, "memory://???")
	m := memory{}
	m.InitViewGroup(&m, "memory://???")

	f.AddAndSetViewGroup(&w, "")
	f.AddAndSetViewGroup(&m, "")

	fmt.Println("---")

	SetViewGroupOther(&m, "hey")

	fmt.Println("---")

	SetViewGroupOther(&w, "hey from websocket")

	fmt.Println("---")

	m.RemoveView(&w)

	SetViewGroupOther(&m, "hey 2")
}
