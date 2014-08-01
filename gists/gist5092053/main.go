package gist5092053

import (
	. "gist.github.com/5408860.git"
	"sort"
)

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func SortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(Reverse{p})
	return p
}

// A data structure to hold a key/value pair.
type RuneIntPair struct {
	Key   rune
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type RuneIntPairList []RuneIntPair

func (p RuneIntPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RuneIntPairList) Len() int           { return len(p) }
func (p RuneIntPairList) Less(i, j int) bool { return p[i].Key < p[j].Key }

func SortMapByKey(m map[rune]int, rev bool) RuneIntPairList {
	sm := make(RuneIntPairList, len(m))
	i := 0
	for k, v := range m {
		sm[i] = RuneIntPair{k, v}
		i++
	}
	if !rev {
		sort.Sort(sm)
	} else {
		sort.Sort(Reverse{sm})
	}
	return sm
}

func main() {
	m := map[string]int{
		"blah": 5,
		"boo":  9,
		"yah":  1,
	}

	sm := SortMapByValue(m)

	for _, v := range sm {
		println(v.Value, v.Key)
	}
}
