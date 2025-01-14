package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

type moduleCodec struct{}

func (*moduleCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*moduleCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	mod := *((*ast.Module)(ptr))

	stream.WriteObjectStart()

	hasWritten := false

	if mod.Package != nil {
		stream.WriteObjectField(strPackage)

		if len(mod.Annotations) > 0 {
			stream.Attachment = mod.Annotations
		}

		stream.WriteVal(mod.Package)

		stream.Attachment = nil
		hasWritten = true
	}

	if len(mod.Imports) > 0 {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strImports)
		stream.WriteArrayStart()

		for i, imp := range mod.Imports {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(imp)
		}

		stream.WriteArrayEnd()

		hasWritten = true
	}

	if len(mod.Rules) > 0 {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strRules)
		stream.WriteArrayStart()

		for i, rule := range mod.Rules {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(rule)
		}

		stream.WriteArrayEnd()

		hasWritten = true
	}

	if len(mod.Comments) > 0 {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strComments)
		stream.WriteArrayStart()

		for i, comment := range mod.Comments {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(comment)
		}

		stream.WriteArrayEnd()
	}

	stream.WriteObjectEnd()
}
