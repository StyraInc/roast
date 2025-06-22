// Package providing tools for working with Rego's AST library (not Roast)
package rast

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/open-policy-agent/opa/v1/ast"
)

// UnquotedPath returns a slice of strings from a path without quotes.
// e.g. data.foo["bar"] -> ["foo", "bar"], note that the data is not included.
func UnquotedPath(path ast.Ref) []string {
	ret := make([]string, 0, len(path)-1)
	for _, ref := range path[1:] {
		ret = append(ret, strings.Trim(ref.String(), `"`))
	}

	return ret
}

// RefStringToBody converts a simple dot-delimited string path to an ast.Body.
// This is a lightweight alternative to ast.ParseBody that avoids the overhead of parsing,
// and benefits from using interned terms when possible. It is also nowhere near as competent,
// and can only handle simple string paths without vars, numbers, etc. Suitable for use with
// e.g. rego.ParsedQuery and other places where a simple ref is needed. Do *NOT* use the returned
// ast.Body anywhere it might be mutated (like having location data added), as that modifies the
// globally interned terms.
//
// Implementations tested:
// -----------------------
// 333.6 ns/op	     472 B/op	      19 allocs/op - SplitSeq
// 330.7 ns/op	     496 B/op	      16 allocs/op - Split
// 269.1 ns/op	     400 B/op	      15 allocs/op - IndexOf for loop (current)
func RefStringToBody(path string) ast.Body {
	var i int
	if i = strings.Index(path, "."); i == -1 {
		return ast.NewBody(ast.NewExpr(ast.RefTerm(refHeadTerm(path))))
	}

	terms := append(make([]*ast.Term, 0, strings.Count(path, ".")+1), refHeadTerm(path[:i]))

	for {
		path = path[i+1:]
		if i = strings.Index(path, "."); i == -1 {
			if len(path) > 0 {
				terms = append(terms, ast.InternedStringTerm(path))
			}
			break
		}
		terms = append(terms, ast.InternedStringTerm(path[:i]))
	}

	return ast.NewBody(ast.NewExpr(ast.RefTerm(terms...)))
}

func RefStringToRef(path string) ast.Ref {
	var i int
	if i = strings.Index(path, "."); i == -1 {
		return ast.Ref([]*ast.Term{refHeadTerm(path)})
	}

	terms := append(make([]*ast.Term, 0, strings.Count(path, ".")+1), refHeadTerm(path[:i]))

	for {
		path = path[i+1:]
		if i = strings.Index(path, "."); i == -1 {
			if len(path) > 0 {
				terms = append(terms, ast.InternedStringTerm(path))
			}
			break
		}
		terms = append(terms, ast.InternedStringTerm(path[:i]))
	}

	return ast.Ref(terms)
}

func refHeadTerm(name string) *ast.Term {
	switch name {
	case "data":
		return ast.DefaultRootDocument
	case "input":
		return ast.InputRootDocument
	default:
		return ast.VarTerm(name)
	}
}

// StructToValue converts a struct to ast.Value using 'json' struct tags (e.g., `json:"field,omitempty"`)
// but without an expensive JSON roundtrip.
func StructToValue(input any) ast.Value {
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	kvs := make([][2]*ast.Term, 0, t.NumField())
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		value := v.Field(i)
		if strings.Contains(tag, ",") {
			parts := strings.Split(tag, ",")
			tag = parts[0]
			omitempty := slices.Contains(parts[1:], "omitempty")
			if omitempty && isZeroValue(value) {
				continue
			}
		}
		kvs = append(kvs, ast.Item(ast.InternedStringTerm(tag), ast.NewTerm(toAstValue(value.Interface()))))
	}
	return ast.NewObject(kvs...)
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Array:
		return v.IsNil() || v.Len() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		for i := range v.NumField() {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

func toAstValue(v any) ast.Value {
	if v == nil {
		return ast.NullValue
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ast.NullValue
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		return StructToValue(rv.Interface())
	case reflect.Slice, reflect.Array:
		l := rv.Len()
		if l == 0 {
			return ast.InternedEmptyArrayValue
		}
		arr := make([]*ast.Term, 0, l)
		for i := range l {
			arr = append(arr, internedAny(rv.Index(i).Interface()))
		}
		return ast.NewArray(arr...)
	case reflect.Map:
		kvs := make([][2]*ast.Term, 0, rv.Len())
		for _, key := range rv.MapKeys() {
			var k *ast.Term
			ki := key.Interface()
			if s, ok := ki.(string); ok {
				k = ast.InternedStringTerm(s)
			} else {
				k = ast.InternedStringTerm(fmt.Sprintf("%v", ki))
			}
			kvs = append(kvs, [2]*ast.Term{k, internedAny(rv.MapIndex(key).Interface())})
		}
		return ast.NewObject(kvs...)
	case reflect.String:
		return ast.String(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ast.Number(fmt.Sprintf("%d", rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ast.Number(fmt.Sprintf("%d", rv.Uint()))
	case reflect.Float32, reflect.Float64:
		return ast.Number(fmt.Sprintf("%v", rv.Float()))
	case reflect.Bool:
		return ast.InternedBooleanTerm(rv.Bool()).Value
	}
	// Fallback: string representation
	fmt.Println("WARNING: Unsupported type for conversion to ast.Value:", rv.Kind())
	return ast.String(fmt.Sprintf("%v", v))
}

func internedAny(v any) *ast.Term {
	switch value := any(v).(type) {
	case bool:
		return ast.InternedBooleanTerm(value)
	case string:
		return ast.InternedStringTerm(value)
	case int:
		return ast.InternedIntNumberTerm(value)
	case uint:
		return ast.InternedIntNumberTerm(int(value))
	case int64:
		return ast.InternedIntNumberTerm(int(value))
	case float64:
		return ast.FloatNumberTerm(value)
	default:
		return ast.NewTerm(toAstValue(v))
	}
}
