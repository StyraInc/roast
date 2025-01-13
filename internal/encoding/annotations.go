package encoding

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

type annotationsCodec struct{}

func (*annotationsCodec) IsEmpty(_ unsafe.Pointer) bool {
	return false
}

func (*annotationsCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	a := *((*ast.Annotations)(ptr))

	stream.WriteObjectStart()

	if a.Location != nil {
		stream.WriteObjectField(strLocation)
		stream.WriteVal(a.Location)
		stream.WriteMore()
	}

	stream.WriteObjectField(strScope)
	stream.WriteString(a.Scope)

	if a.Title != "" {
		stream.WriteMore()
		stream.WriteObjectField(strTitle)
		stream.WriteString(a.Title)
	}

	if a.Description != "" {
		stream.WriteMore()
		stream.WriteObjectField(strDescription)
		stream.WriteString(a.Description)
	}

	if a.Entrypoint {
		stream.WriteMore()
		stream.WriteObjectField(strEntrypoint)
		stream.WriteBool(a.Entrypoint)
	}

	if len(a.Organizations) > 0 {
		stream.WriteMore()
		stream.WriteObjectField(strOrganizations)
		stream.WriteArrayStart()

		for i, org := range a.Organizations {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteString(org)
		}

		stream.WriteArrayEnd()
	}

	if len(a.RelatedResources) > 0 {
		stream.WriteMore()
		stream.WriteObjectField(strRelatedResources)
		stream.WriteArrayStart()

		for i, res := range a.RelatedResources {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(res)
		}

		stream.WriteArrayEnd()
	}

	if len(a.Authors) > 0 {
		stream.WriteMore()
		stream.WriteObjectField(strAuthors)
		stream.WriteArrayStart()

		for i, author := range a.Authors {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(author)
		}

		stream.WriteArrayEnd()
	}

	if len(a.Schemas) > 0 {
		stream.WriteMore()
		stream.WriteObjectField(strSchemas)
		stream.WriteArrayStart()

		for i, schema := range a.Schemas {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteVal(schema)
		}

		stream.WriteArrayEnd()
	}

	if len(a.Custom) > 0 {
		stream.WriteMore()
		stream.WriteObjectField(strCustom)
		stream.WriteObjectStart()

		i := 0

		for key, value := range a.Custom {
			if i > 0 {
				stream.WriteMore()
			}

			stream.WriteObjectField(key)
			stream.WriteVal(value)

			i++
		}

		stream.WriteObjectEnd()
	}

	stream.WriteObjectEnd()
}
