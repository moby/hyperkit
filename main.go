package main

import (
	"fmt"
	. "gist.github.com/5092053.git"

	"io/ioutil"
	. "gist.github.com/5286084.git"
)

var _ = ioutil.ReadFile
var _ = CheckError

func PrintRuneStats(s string) {
	r := []rune(s)
	fmt.Printf("Total runes: %v\n", len(r))

	m := map[rune]int{}
	for _, v := range r {
		m[v]++
	}
	fmt.Printf("Total unique runes: %v\n\n", len(m))

	sm := SortMapByKey(m, true)

	//for i := len(sm) - 1; i >= 0; i-- { v := sm[i]
	for _, v := range sm {
		fmt.Printf("%v\t\t%q \t%v\n", v.Key, v.Key, v.Value)
	}
}

func main() {
	x := "abc   Z"
	//b, err := ioutil.ReadFile("/Users/Dmitri/Desktop/out.md"); CheckError(err)
	//x := string(b)

	PrintRuneStats(x)
}