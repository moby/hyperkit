package gist6445065

import (
	//"fmt"
	"github.com/shurcooL/go-goon"
	"reflect"

	"go/ast"
	"go/token"
)

var _ = goon.Dump
var _ token.Pos

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

func (s *state) rdump(v reflect.Value, query func(i interface{}) bool) interface{} {
	if query(v.Interface()) {
		return v.Interface()
	}

	switch v.Kind() {
	/*case reflect.Int:
		println(v.Int())
	case reflect.String:
		fmt.Printf("%q\n", v.String())*/
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			//print(v.Field(i).Type().String() + ": ")
			q := s.rdump(unpackValue(v.Field(i)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			q := s.rdump(unpackValue(v.MapIndex(key)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			q := s.rdump(unpackValue(v.Index(i)), query)
			if q != nil {
				return q
			}
		}
	case reflect.Ptr:
		//println("reflect.Ptr", v.Type().String())
		if !v.IsNil() {
			if !s.Visited[v.Pointer()] {
				s.Visited[v.Pointer()] = true
				q := s.rdump(v.Elem(), query)
				if q != nil {
					return q
				}
			}
		} else {
			//println("nil")
		}
	case reflect.Interface:
		//println("reflect.Interface", v.Type().String())
		if !v.IsNil() {
			q := s.rdump(v.Elem(), query)
			if q != nil {
				return q
			}
		} else {
			//println("nil")
		}
	default:
		//println("<unsupported>")
	}

	return nil
}

func FindFirst(d interface{}, query func(i interface{}) bool) interface{} {
	s := state{Visited: make(map[uintptr]bool)}
	return s.rdump(reflect.ValueOf(d), query)
}

func Dump(a ...interface{}) {
	query := func(i interface{}) bool {
		if v, ok := i.(*ast.FuncLit); ok {
			println(">>>>>>>>>>>>>>>>>>>>>>", v.Type.Func)
			return true
		}

		return false
	}

	for _, arg := range a {
		s := state{Visited: make(map[uintptr]bool)}
		s.rdump(reflect.ValueOf(arg), query)
		println("\n---\n")
		goon.Dump(s.Visited)
	}
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
		Rhs   []ast.Expr
	}

	x := Lang{
		Name:  "Go",
		Year:  2009,
		URL:   "http",
		Inner: &Inner{},
		Rhs: ([]ast.Expr)([]ast.Expr{
			(*ast.FuncLit)(&ast.FuncLit{
				Type: (*ast.FuncType)(&ast.FuncType{
					Func: (token.Pos)(1112),
					Params: (*ast.FieldList)(&ast.FieldList{
						Opening: (token.Pos)(1116),
						List:    ([]*ast.Field)([]*ast.Field{}),
						Closing: (token.Pos)(1117),
					}),
					Results: (*ast.FieldList)(nil),
				}),
				Body: (*ast.BlockStmt)(&ast.BlockStmt{
					Lbrace: (token.Pos)(1119),
					List: ([]ast.Stmt)([]ast.Stmt{
						(*ast.ExprStmt)(&ast.ExprStmt{
							X: (*ast.CallExpr)(&ast.CallExpr{
								Fun: (*ast.Ident)(&ast.Ident{
									NamePos: (token.Pos)(1123),
									Name:    (string)("println"),
									Obj:     (*ast.Object)(nil),
								}),
								Lparen: (token.Pos)(1130),
								Args: ([]ast.Expr)([]ast.Expr{
									(*ast.BasicLit)(&ast.BasicLit{
										ValuePos: (token.Pos)(1131),
										Kind:     (token.Token)(9),
										Value:    (string)("\"Hello from anon func!\""),
									}),
								}),
								Ellipsis: (token.Pos)(0),
								Rparen:   (token.Pos)(1154),
							}),
						}),
					}),
					Rbrace: (token.Pos)(1157),
				}),
			}),
		}),
	}

	//x.Inner.Field3 = &Inner{}
	x.Inner.Field3 = x.Inner

	goon.Dump(x)
	println("\n---\n")
	Dump(x)
}
