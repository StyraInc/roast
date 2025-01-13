package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

type bodyCodec struct{}

func (*bodyCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*bodyCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	body := *((*ast.Body)(ptr))

	stream.WriteArrayStart()

	for i, expr := range body {
		if i > 0 {
			stream.WriteMore()
		}

		stream.WriteVal(expr)
	}

	stream.WriteArrayEnd()
}
