package main

import (
	"testing"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/go-goon"
)

var _= fmt.Printf
var _ = spew.Dump
var _ = goon.Dump

func Test(t *testing.T) {
	if ".olleH" != Reverse("Hello.") {
		t.Fail()
	}
}