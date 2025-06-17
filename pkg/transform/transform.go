package transform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/open-policy-agent/opa/v1/ast"

	"github.com/styrainc/roast/internal/transforms"
	"github.com/styrainc/roast/pkg/encoding"

	_ "github.com/styrainc/roast/internal/encoding"
)

var (
	pathSeparatorTerm = ast.StringTerm(string(os.PathSeparator))

	environment [2]*ast.Term = ast.Item(ast.InternedStringTerm("environment"), ast.ObjectTerm(
		ast.Item(ast.InternedStringTerm("path_separator"), pathSeparatorTerm),
	))

	operationsLintItem = ast.Item(
		ast.InternedStringTerm("operations"),
		ast.ArrayTerm(ast.InternedStringTerm("lint")),
	)
	operationsLintCollectItem = ast.Item(ast.InternedStringTerm("operations"), ast.ArrayTerm(
		ast.InternedStringTerm("lint"),
		ast.InternedStringTerm("collect")),
	)
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

	if err = jsoniter.ConfigFastest.Unmarshal(bs, x); err != nil {
		return encoding.SafeNumberConfig.Unmarshal(bs, x)
	}

	return nil
}

func ToAST(name string, content string, module *ast.Module, collect bool) (ast.Value, error) {
	var preparedAST map[string]any

	if err := encoding.JSONRoundTrip(module, &preparedAST); err != nil {
		return nil, fmt.Errorf("JSON rountrip failed for module: %w", err)
	}

	astObj, err := ToOPAInputValue(preparedAST)
	if err != nil {
		return nil, fmt.Errorf("failed to convert prepared AST to OPA input value: %w", err)
	}

	if input, ok := astObj.(ast.Object); ok {
		abs, _ := filepath.Abs(name)

		var operations [2]*ast.Term
		if collect {
			operations = operationsLintCollectItem
		} else {
			operations = operationsLintItem
		}

		input.Insert(ast.InternedStringTerm("regal"), ast.ObjectTerm(
			ast.Item(ast.InternedStringTerm("file"), ast.ObjectTerm(
				ast.Item(ast.InternedStringTerm("name"), ast.StringTerm(name)),
				ast.Item(ast.InternedStringTerm("lines"), linesArrayTerm(content)),
				ast.Item(ast.InternedStringTerm("abs"), ast.StringTerm(abs)),
				ast.Item(ast.InternedStringTerm("rego_version"), ast.InternedStringTerm(module.RegoVersion().String())),
			)),
			environment,
			operations,
		))

		return input, nil
	}

	return nil, errors.New("prepared AST failed")
}

func linesArrayTerm(content string) *ast.Term {
	parts := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	terms := make([]*ast.Term, len(parts))

	for i := range parts {
		terms[i] = ast.InternedStringTerm(parts[i])
	}

	return ast.ArrayTerm(terms...)
}
