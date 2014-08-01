// +build ignore

package main

import . "gist.github.com/7802150.git"

import (
	"fmt"
	. "gist.github.com/6418290.git"
)

type file struct {
	ViewGroup
}

func (*file) SetSelf(s string) {
	fmt.Println(GetParentFuncAsString())
}

type websocket struct {
	ViewGroup
}

func (*websocket) SetSelf(s string) {
	fmt.Println(GetParentFuncAsString())
}

type memory struct {
	ViewGroup
}

func (*memory) SetSelf(s string) {
	fmt.Println(GetParentFuncAsString())
}

func main() {
	f := file{}
	f.InitViewGroup(&f)
	w := websocket{}
	w.InitViewGroup(&w)
	m := memory{}
	m.InitViewGroup(&m)

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
