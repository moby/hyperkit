// +build js

// Package jsutil provides utility functions for interacting with
// native JavaScript APIs via syscall/js API.
// It has support for common types in honnef.co/go/js/dom/v2.
package jsutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"syscall/js"

	"honnef.co/go/js/dom/v2"
)

// Wrap returns a wrapper func that handles the conversion from native JavaScript js.Value parameters
// to the following types.
//
// It supports js.Value (left unmodified), dom.Document, dom.Element, dom.Event, dom.HTMLElement, dom.Node.
// It has to be one of those types exactly; it can't be another type that implements the interface like *dom.BasicElement.
//
// For other types, the input is assumed to be a JSON string which is then unmarshalled into that type.
//
// If the number of arguments provided to the wrapped func doesn't match
// the number of arguments for original func, it panics.
//
// Here is example usage:
//
// 	<span onclick="Handler(event, this, {{.SomeStruct | json}});">Example</span>
//
// 	func Handler(event dom.Event, htmlElement dom.HTMLElement, data someStruct) {
// 		data.Foo = ... // Use event, htmlElement, data.
// 	}
//
// 	func main() {
// 		js.Global().Set("Handler", jsutil.Wrap(Handler))
// 	}
func Wrap(fn interface{}) js.Func {
	v := reflect.ValueOf(fn)
	return js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		if len(args) != v.Type().NumIn() {
			panic(fmt.Errorf("wrapped %v got %v arguments, want %v", v.Type().String(), len(args), v.Type().NumIn()))
		}
		in := make([]reflect.Value, v.Type().NumIn())
		for i := range in {
			switch t := v.Type().In(i); t {
			// js.Value is passed through.
			case typeOf((*js.Value)(nil)):
				in[i] = reflect.ValueOf(args[i])

			// dom types are wrapped.
			case typeOf((*dom.Document)(nil)):
				in[i] = reflect.ValueOf(dom.WrapDocument(args[i]))
			case typeOf((*dom.Element)(nil)):
				in[i] = reflect.ValueOf(dom.WrapElement(args[i]))
			case typeOf((*dom.Event)(nil)):
				in[i] = reflect.ValueOf(dom.WrapEvent(args[i]))
			case typeOf((*dom.HTMLElement)(nil)):
				in[i] = reflect.ValueOf(dom.WrapHTMLElement(args[i]))
			case typeOf((*dom.Node)(nil)):
				in[i] = reflect.ValueOf(dom.WrapNode(args[i]))

			// Unmarshal incoming encoded JSON into the Go type.
			default:
				if args[i].Type() != js.TypeString {
					panic(fmt.Errorf("jsutil: incoming value type is %s; want a string with JSON content", args[i].Type()))
				}
				p := reflect.New(t)
				err := json.Unmarshal([]byte(args[i].String()), p.Interface())
				if err != nil {
					panic(fmt.Errorf("jsutil: unmarshaling JSON %q into type %s failed: %v", args[i], t, err))
				}
				in[i] = reflect.Indirect(p)
			}
		}
		v.Call(in)
		return nil
	})
}

// typeOf returns the reflect.Type of what the pointer points to.
func typeOf(pointer interface{}) reflect.Type {
	return reflect.TypeOf(pointer).Elem()
}
