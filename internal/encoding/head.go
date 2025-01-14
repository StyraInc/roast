package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

type headCodec struct{}

func (*headCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*headCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	head := *((*ast.Head)(ptr))

	stream.WriteObjectStart()

	var hasWritten bool

	if head.Location != nil {
		stream.WriteObjectField(strLocation)
		stream.WriteVal(head.Location)

		hasWritten = true
	}

	if head.Reference != nil {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strRef)
		stream.WriteVal(head.Reference)

		hasWritten = true
	}

	if len(head.Args) > 0 {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strArgs)
		writeTermsArray(stream, head.Args)

		hasWritten = true
	}

	if head.Assign {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strAssign)
		stream.WriteBool(head.Assign)

		hasWritten = true
	}

	if head.Key != nil {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strKey)
		stream.WriteVal(head.Key)

		hasWritten = true
	}

	if head.Value != nil {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strValue)
		stream.WriteVal(head.Value)
	}

	stream.WriteObjectEnd()
}
