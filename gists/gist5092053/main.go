// Package gist5092053 offers a func to turn a map[string]int into a PairList, then sort and return it.
package gist5092053

import (
	"sort"

	"github.com/shurcooL/go/gists/gist5408860"
)

// Pair is a data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Pair.Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// SortMapByValue turns a map into a PairList, then sorts and returns it.
func SortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(gist5408860.Reverse{p})
	return p
}

// RuneIntPair is a data structure to hold a key/value pair.
type RuneIntPair struct {
	Key   rune
	Value int
}

// RuneIntPairList is a slice of RuneIntPair that implements sort.Interface to sort by RuneIntPair.Value.
type RuneIntPairList []RuneIntPair

func (p RuneIntPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RuneIntPairList) Len() int           { return len(p) }
func (p RuneIntPairList) Less(i, j int) bool { return p[i].Key < p[j].Key }

// SortMapByKey sorts map by key.
func SortMapByKey(m map[rune]int, reverse bool) RuneIntPairList {
	sm := make(RuneIntPairList, len(m))
	i := 0
	for k, v := range m {
		sm[i] = RuneIntPair{k, v}
		i++
	}
	if !reverse {
		sort.Sort(sm)
	} else {
		sort.Sort(gist5408860.Reverse{sm})
	}
	return sm
}
