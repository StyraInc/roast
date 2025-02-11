package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"

	"github.com/styrainc/roast/internal/encoding/util"
)

type ruleCodec struct{}

func (*ruleCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*ruleCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	rule := *((*ast.Rule)(ptr))

	stream.WriteObjectStart()

	hasWritten := false

	if rule.Location != nil {
		stream.WriteObjectField(strLocation)
		stream.WriteVal(rule.Location)

		hasWritten = true
	}

	if len(rule.Annotations) > 0 {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strAnnotations)
		util.WriteValsArray(stream, rule.Annotations)

		hasWritten = true
	}

	if rule.Default {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strDefault)
		stream.WriteBool(rule.Default)

		hasWritten = true
	}

	if rule.Head != nil {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strHead)
		stream.WriteVal(rule.Head)
	}

	if !isBodyGenerated(&rule) {
		if hasWritten {
			stream.WriteMore()
		}

		stream.WriteObjectField(strBody)
		stream.WriteVal(rule.Body)
	}

	if rule.Else != nil {
		stream.WriteMore()
		stream.WriteObjectField(strElse)
		stream.WriteVal(rule.Else)
	}

	stream.WriteObjectEnd()
}

func isBodyGenerated(rule *ast.Rule) bool {
	if rule.Default {
		return true
	}

	if len(rule.Body) == 0 {
		return true
	}

	if rule.Head == nil {
		return false
	}

	if rule.Body[0] != nil && rule.Body[0].Location == rule.Location {
		return true
	}

	if rule.Body[0] != nil && rule.Head.Value != nil && rule.Body[0].Location == rule.Head.Value.Location {
		return true
	}

	if rule.Head.Key != nil &&
		rule.Body[0].Location.Row == rule.Head.Key.Location.Row &&
		rule.Body[0].Location.Col < rule.Head.Key.Location.Col {
		// This is a quirk in the original AST — the generated body will have a location
		// set before the key, i.e. "message"
		return true
	}

	return false
}
