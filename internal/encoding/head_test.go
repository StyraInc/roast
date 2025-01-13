package encoding

import (
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

func TestRuleHeadEncoding(t *testing.T) {
	t.Parallel()

	head := ast.Head{
		Name: "omitted",
		Reference: ast.Ref{
			{
				Value: ast.String("foo"),
				Location: &ast.Location{
					Row:  1,
					Col:  1,
					Text: []byte("foo"),
				},
			},
			{
				Value: ast.String("bar"),
				Location: &ast.Location{
					Row:  1,
					Col:  5, // following "foo."
					Text: []byte("bar"),
				},
			},
		},

		Value: &ast.Term{
			Value: ast.Boolean(true),
			Location: &ast.Location{
				Row:  1,
				Col:  12, // following "foo.bar := "
				Text: []byte("true"),
			},
		},
		Assign: true,
		Location: &ast.Location{
			Row:  1,
			Col:  1,
			Text: []byte("foo.bar := true"),
		},
	}

	bs, err := jsoniter.ConfigFastest.MarshalIndent(head, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expect := `{
  "location": "1:1:1:16",
  "ref": [
    {
      "location": "1:1:1:4",
      "type": "string",
      "value": "foo"
    },
    {
      "location": "1:5:1:8",
      "type": "string",
      "value": "bar"
    }
  ],
  "assign": true,
  "value": {
    "location": "1:12:1:16",
    "type": "boolean",
    "value": true
  }
}`

	if string(bs) != expect {
		t.Fatalf("expected %s but got %s", expect, string(bs))
	}
}
