package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

type termCodec struct{}

func (*termCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*termCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	term := *((*ast.Term)(ptr))

	stream.WriteObjectStart()

	if term.Location != nil {
		stream.WriteObjectField(strLocation)
		stream.WriteVal(term.Location)
	}

	if term.Value != nil {
		if term.Location != nil {
			stream.WriteMore()
		}

		stream.WriteObjectField(strType)
		stream.WriteString(valueName(term.Value))
		stream.WriteMore()
		stream.WriteObjectField(strValue)
		stream.WriteVal(term.Value)
	}

	stream.WriteObjectEnd()
}

// TODO: remove once this is in OPA and we can use that instead.
func valueName(x ast.Value) string {
	switch x.(type) {
	case ast.String:
		return "string"
	case ast.Boolean:
		return "boolean"
	case ast.Number:
		return "number"
	case ast.Null:
		return "null"
	case ast.Var:
		return "var"
	case ast.Object:
		return "object"
	case ast.Set:
		return "set"
	case ast.Ref:
		return "ref"
	case ast.Call:
		return "call"
	case *ast.Array:
		return "array"
	case *ast.ArrayComprehension:
		return "arraycomprehension"
	case *ast.ObjectComprehension:
		return "objectcomprehension"
	case *ast.SetComprehension:
		return "setcomprehension"
	}

	return ast.TypeName(x)
}
