package gist6445065

import (
	"reflect"
)

type state struct {
	Visited map[uintptr]bool
}

func unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface && !v.IsNil() {
		return v.Elem()
	} else {
		return v
	}
}

func (s *state) findFirst(v reflect.Value, query func(i interface{}) bool) interface{} {
	// TODO: Should I check v.CanInterface()? Maybe I can get away without it...
	if query(v.Interface()) {
		return v.Interface()
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			q := s.findFirst(unpackValue(v.Field(i)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			q := s.findFirst(unpackValue(v.MapIndex(key)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			q := s.findFirst(unpackValue(v.Index(i)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Ptr:
		if !v.IsNil() {
			if !s.Visited[v.Pointer()] {
				s.Visited[v.Pointer()] = true
				q := s.findFirst(v.Elem(), query)
				if q != nil {
					return q
				}
			}
		}
	}

	return nil
}

func FindFirst(d interface{}, query func(i interface{}) bool) interface{} {
	s := state{Visited: make(map[uintptr]bool)}
	return s.findFirst(unpackValue(reflect.ValueOf(d)), query)
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
