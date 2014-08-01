package gist5408736

import (
	"fmt"
	"io/ioutil"

	. "gist.github.com/5092053.git"
	. "gist.github.com/5286084.git"
	. "gist.github.com/5571468.git"
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
		fmt.Printf("%v   \t%q \t%v\n", v.Key, v.Key, v.Value)
	}
}

func main() {
	//s := "abc   Z"
	s := MustReadFile("/Users/Dmitri/Desktop/1.txt")

	PrintRuneStats(s)
}
