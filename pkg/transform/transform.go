package transform

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"

	"github.com/anderseknert/roast/internal/transforms"

	"github.com/open-policy-agent/opa/v1/ast"

	_ "github.com/anderseknert/roast/internal/encoding"
)

// InterfaceToValue converts a native Go value x to a Value.
// This is an optimized version of the same function in the OPA codebase,
// and optimized in a way that makes it useful only for a map[string]any
// unmarshaled from RoAST JSON. Don't use it for anything else.
func AnyToValue(x any) (ast.Value, error) {
	return transforms.AnyToValue(x)
}

// ToOPAInputValue converts provided x to an ast.Value suitable for use as
// parsed input to OPA (`rego.EvalParsedInput`). This will have the value
// pass through the same kind of roundtrip as OPA would otherwise have to
// do when provided unparsed input, but much more efficiently as both JSON
// marshalling and the custom InterfaceToValue function provided here are
// optimized for performance.
func ToOPAInputValue(x any) (ast.Value, error) {
	ptr := reference(x)
	if err := anyPtrRoundTrip(ptr); err != nil {
		return nil, err
	}

	value, err := AnyToValue(*ptr)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// From OPA's util package
//
// Reference returns a pointer to its argument unless the argument already is
// a pointer. If the argument is **t, or ***t, etc, it will return *t.
//
// Used for preparing Go types (including pointers to structs) into values to be
// put through util.RoundTrip().
func reference(x any) *any {
	var y any

	rv := reflect.ValueOf(x)
	if rv.Kind() == reflect.Ptr {
		return reference(rv.Elem().Interface())
	}

	if rv.Kind() != reflect.Invalid {
		y = rv.Interface()

		return &y
	}

	return &x
}

func anyPtrRoundTrip(x *any) error {
	bs, err := jsoniter.ConfigFastest.Marshal(x)
	if err != nil {
		return err
	}

	return jsoniter.ConfigFastest.Unmarshal(bs, x)
}
