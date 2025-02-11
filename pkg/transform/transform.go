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
	"github.com/styrainc/roast/pkg/intern"

	_ "github.com/styrainc/roast/internal/encoding"
)

var (
	pathSeparatorTerm = ast.StringTerm(string(os.PathSeparator))

	environment [2]*ast.Term = ast.Item(intern.StringTerm("environment"), ast.ObjectTerm(
		ast.Item(intern.StringTerm("path_separator"), pathSeparatorTerm),
	))

	operationsLintItem        = ast.Item(intern.StringTerm("operations"), ast.ArrayTerm(intern.StringTerm("lint")))
	operationsLintCollectItem = ast.Item(intern.StringTerm("operations"), ast.ArrayTerm(
		intern.StringTerm("lint"),
		intern.StringTerm("collect")),
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

	return jsoniter.ConfigFastest.Unmarshal(bs, x)
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

		input.Insert(intern.StringTerm("regal"), ast.ObjectTerm(
			ast.Item(intern.StringTerm("file"), ast.ObjectTerm(
				ast.Item(intern.StringTerm("name"), ast.StringTerm(name)),
				ast.Item(intern.StringTerm("lines"), linesArrayTerm(content)),
				ast.Item(intern.StringTerm("abs"), ast.StringTerm(abs)),
				ast.Item(intern.StringTerm("rego_version"), intern.StringTerm(module.RegoVersion().String())),
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
		terms[i] = intern.StringTerm(parts[i])
	}

	return ast.ArrayTerm(terms...)
}
