package gist6724654

import (
	"reflect"
	"unsafe"
)

// reflectValue mirrors the struct layout of the reflect package Value type.
var reflectValue struct {
	typ  unsafe.Pointer
	val  unsafe.Pointer
	flag uintptr
}

// flagIndir indicates whether the value field of a reflect.Value is the actual
// data or a pointer to the data.
const flagIndir = 1 << 1

// unsafeReflectValue converts the passed reflect.Value into a one that bypasses
// the typical safety restrictions preventing access to unaddressable and
// unexported data.  It works by digging the raw pointer to the underlying
// value out of the protected value and generating a new unprotected (unsafe)
// reflect.Value to it.
//
// This allows us to check for implementations of the Stringer and error
// interfaces to be used for pretty printing ordinarily unaddressable and
// inaccessible values such as unexported struct fields.
func UnsafeReflectValue(v reflect.Value) (rv reflect.Value) {
	indirects := 1
	vt := v.Type()
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + unsafe.Offsetof(reflectValue.val))
	rvf := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + unsafe.Offsetof(reflectValue.flag)))
	if rvf&flagIndir != 0 {
		vt = reflect.PtrTo(v.Type())
		indirects++
	}

	pv := reflect.NewAt(vt, upv)
	rv = pv
	for i := 0; i < indirects; i++ {
		rv = rv.Elem()
	}
	return rv
}
