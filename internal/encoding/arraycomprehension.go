package encoding

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/open-policy-agent/opa/ast"
	"unsafe"
)

type arrayComprehensionCodec struct{}

func (*arrayComprehensionCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*arrayComprehensionCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ac := *((*ast.ArrayComprehension)(ptr))

	stream.WriteObjectStart()

	stream.WriteObjectField("term")
	stream.WriteVal(ac.Term)
	stream.WriteMore()
	stream.WriteObjectField("body")
	stream.WriteVal(ac.Body)

	stream.WriteObjectEnd()
}
