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

type state2 struct {
	state
	Found map[interface{}]bool
}

func (s *state2) findAll(v reflect.Value, query func(i interface{}) bool) {
	//if !v.IsValid() { return }
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		// TODO: Instead of skipping nil values, maybe pass the info as a bool parameter to query?
		if v.IsNil() {
			return
		}
	}

	// TODO: Should I check v.CanInterface()? It seems like I might be able to get away without it...
	if query(v.Interface()) {
		s.Found[v.Interface()] = true
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			s.findAll(v.Field(i), query)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			s.findAll(v.MapIndex(key), query)
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			s.findAll(v.Index(i), query)
		}
	case reflect.Ptr:
		if !v.IsNil() {
			if !s.Visited[v.Pointer()] {
				s.Visited[v.Pointer()] = true
				s.findAll(v.Elem(), query)
			}
		}
	case reflect.Interface:
		if !v.IsNil() {
			s.findAll(v.Elem(), query)
		}
	}
}

func FindAll(d interface{}, query func(i interface{}) bool) map[interface{}]bool {
	s := state2{state: state{Visited: make(map[uintptr]bool)}, Found: make(map[interface{}]bool)}
	s.findAll(reflect.ValueOf(d), query)
	return s.Found
}
