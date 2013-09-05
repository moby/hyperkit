package gist6445065

import (
	"reflect"
)

type state struct {
	Visited map[uintptr]bool
}

func (s *state) findFirst(v reflect.Value, query func(i interface{}) bool) interface{} {
	// TODO: Should I check v.CanInterface()? It seems like I might be able to get away without it...
	if query(v.Interface()) {
		return v.Interface()
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if q := s.findFirst(v.Field(i), query); q != nil {
				return q
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if q := s.findFirst(v.MapIndex(key), query); q != nil {
				return q
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if q := s.findFirst(v.Index(i), query); q != nil {
				return q
			}
		}
	case reflect.Ptr:
		if !v.IsNil() {
			if !s.Visited[v.Pointer()] {
				s.Visited[v.Pointer()] = true
				if q := s.findFirst(v.Elem(), query); q != nil {
					return q
				}
			}
		}
	case reflect.Interface:
		if !v.IsNil() {
			if q := s.findFirst(v.Elem(), query); q != nil {
				return q
			}
		}
	}

	return nil
}

func FindFirst(d interface{}, query func(i interface{}) bool) interface{} {
	s := state{Visited: make(map[uintptr]bool)}
	return s.findFirst(reflect.ValueOf(d), query)
}

func main() {
	type Inner struct {
		Field1 string
		Field2 int
		Field3 *Inner
	}
	type Lang struct {
		Name  string
		Year  int
		URL   string
		Inner *Inner
	}

	x := Lang{
		Name:  "Go",
		Year:  2009,
		URL:   "http",
		Inner: &Inner{},
	}

	//x.Inner.Field3 = &Inner{}
	x.Inner.Field3 = x.Inner

	//goon.Dump(x)
	println("\n---\n")
	//Dump(x)
	// ... unfinished test
}
