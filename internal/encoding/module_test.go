package encoding

import (
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"
)

func TestAnnotationsOnPackage(t *testing.T) {
	t.Parallel()

	module := ast.Module{
		Package: &ast.Package{
			Location: &ast.Location{
				Row:  3,
				Col:  1,
				Text: []byte("foo"),
			},
			Path: ast.Ref{
				ast.DefaultRootDocument,
				ast.StringTerm("foo"),
			},
		},
		Annotations: []*ast.Annotations{
			{
				Location: &ast.Location{
					Row: 1,
					Col: 1,
				},
				Scope: "package",
				Title: "foo",
			},
		},
	}

	json := jsoniter.ConfigFastest

	roast, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal annotations: %v", err)
	}

	// package annotations should end up on the package object
	// and *not* on the module object, contrary to how OPA
	// currently does it

	expected := `{
  "package": {
    "location": "3:1:3:4",
    "path": [
      {
        "type": "var",
        "value": "data"
      },
      {
        "type": "string",
        "value": "foo"
      }
    ],
    "annotations": [
      {
        "location": "1:1:1:1",
        "scope": "package",
        "title": "foo"
      }
    ]
  }
}`

	if string(roast) != expected {
		t.Fatalf("expected %s but got %s", expected, roast)
	}
}

func TestAnnotationsOnPackageBothPackageAndSubpackagesScope(t *testing.T) {
	t.Parallel()

	module := ast.Module{
		Package: &ast.Package{
			Location: &ast.Location{
				Row:  6,
				Col:  1,
				Text: []byte("foo"),
			},
			Path: ast.Ref{
				ast.DefaultRootDocument,
				ast.StringTerm("foo"),
			},
		},
		Annotations: []*ast.Annotations{
			{
				Location: &ast.Location{
					Row: 1,
					Col: 1,
				},
				Scope: "package",
				Title: "foo",
			},
			{
				Location: &ast.Location{
					Row: 3,
					Col: 1,
				},
				Scope: "subpackages",
				Title: "bar",
			},
		},
	}

	json := jsoniter.ConfigFastest

	roast, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal annotations: %v", err)
	}

	expected := `{
  "package": {
    "location": "6:1:6:4",
    "path": [
      {
        "type": "var",
        "value": "data"
      },
      {
        "type": "string",
        "value": "foo"
      }
    ],
    "annotations": [
      {
        "location": "1:1:1:1",
        "scope": "package",
        "title": "foo"
      },
      {
        "location": "3:1:3:1",
        "scope": "subpackages",
        "title": "bar"
      }
    ]
  }
}`

	if string(roast) != expected {
		t.Fatalf("expected %s but got %s", expected, roast)
	}
}
