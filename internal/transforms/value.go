//nolint:gochecknoglobals
package transforms

import (
	"fmt"
	"strconv"

	"github.com/anderseknert/roast/pkg/util"

	"github.com/open-policy-agent/opa/v1/ast"
)

var (
	// These are strings commonly found in the AST of all Rego policies,
	// like the name of built-in functions, keywords, etc.
	regoStrings = [...]string{
		"",
		" ",
		",",
		"/",
		"array",
		"assign",
		"data",
		"description",
		"equal",
		"file",
		"input",
		"internal",
		"member_2",
		"number",
		"object",
		"policy",
		"rego",
		"set",
		"type",
		"var",
		"string",
		"text",
		"v1",
		"union",
		"IE1FVEFEQVRB", // METADATA as enconded in the comment nodes
	}

	// These are strings commonly found in linter policies, but
	// not necessarily anywhere else.
	regalStrings = [...]string{
		"ast",
		"boolean",
		"bugs",
		"call",
		"category",
		"col",
		"config",
		"error",
		"idiomatic",
		"level",
		"location",
		"module",
		"violation",
		"title",
		"term",
		"r",
		"ref",
		"regal",
		"report",
		"result",
		"row",
		"rule",
		"rules",
		"style",
		"value",
		"end",
	}

	roastKeys = [...]string{
		// OPA / RoAST keys
		"alias",
		"assign",
		"authors",
		"body",
		"custom",
		"default",
		"description",
		"else",
		"entrypoint",
		"head",
		"imports",
		"rules",
		"package",
		"annotations",
		"comments",
		"related_resources",
		"scope",
		"symbols",
		"negated",
		"key",
		"term",
		"domain",
		"location",
		"type",
		"value",
		"path",
		"args",
		"name",
		"schema",
		"schemas",
		"terms",
		"text",
		"title",
		"ref",
		"with",
		"target",
		// Regal specific keys
		"file",
		"abs",
		"environment",
		"path_separator",
		"lines",
		"operations",
		"regal",
		"severity",
	}

	nullValue = [1]ast.Value{
		ast.Null{},
	}

	roastKeyValues map[string]ast.Value

	roastKeyTerms map[string]*ast.Term

	commonStringValues map[string]ast.Value

	commonStringTerms map[string]*ast.Term

	emptyObject = ast.NewObject()
	emptyArray  = ast.NewArray()
)

func init() {
	roastKeyValues = make(map[string]ast.Value, len(roastKeys))
	roastKeyTerms = make(map[string]*ast.Term, len(roastKeys))

	for _, k := range roastKeys {
		roastKeyValues[k] = ast.String(k)
		roastKeyTerms[k] = ast.NewTerm(roastKeyValues[k])
	}

	commonStringValues = make(map[string]ast.Value, len(regoStrings)+len(regalStrings))
	commonStringTerms = make(map[string]*ast.Term, len(regoStrings)+len(regalStrings))

	for _, s := range regoStrings {
		commonStringValues[s] = ast.String(s)
		commonStringTerms[s] = ast.NewTerm(commonStringValues[s])
	}

	for _, s := range regalStrings {
		commonStringValues[s] = ast.String(s)
		commonStringTerms[s] = ast.NewTerm(commonStringValues[s])
	}
}

// AnyToValue converts a native Go value x to a Value.
// This is an optimized version of the same function in the OPA codebase,
// and optimized in a way that makes it useful only for a map[string]any
// unmarshaled from RoAST JSON. Don't use it for anything else.
func AnyToValue(x any) (ast.Value, error) {
	switch x := x.(type) {
	case nil:
		return nullValue[0], nil
	case bool:
		return ast.InternedBooleanTerm(x).Value, nil
	case float64:
		ix := int(x)
		if x == float64(ix) {
			return ast.InternedIntNumberTerm(ix).Value, nil
		}

		return ast.Number(strconv.FormatFloat(x, 'g', -1, 64)), nil
	case string:
		if s, ok := commonStringValues[x]; ok {
			return s, nil
		}

		return ast.String(x), nil
	case []string:
		if len(x) == 0 {
			return emptyArray, nil
		}

		r := util.NewPtrSlice[ast.Term](len(x))

		for i, e := range x {
			if s, ok := commonStringValues[e]; ok {
				r[i].Value = s
			} else {
				r[i].Value = ast.String(e)
			}
		}

		return ast.NewArray(r...), nil
	case []any:
		if len(x) == 0 {
			return emptyArray, nil
		}

		r := util.NewPtrSlice[ast.Term](len(x))

		for i, e := range x {
			e, err := AnyToValue(e)
			if err != nil {
				return nil, err
			}

			r[i].Value = e
		}

		return ast.NewArray(r...), nil
	case map[string]any:
		if len(x) == 0 {
			return emptyObject, nil
		}

		kvs := util.NewPtrSlice[ast.Term](len(x) * 2)
		idx := 0

		for k, v := range x {
			if t, ok := roastKeyValues[k]; ok {
				kvs[idx].Value = t
			} else {
				t, err := AnyToValue(k)
				if err != nil {
					return nil, err
				}

				kvs[idx].Value = t
			}

			v, err := AnyToValue(v)
			if err != nil {
				return nil, err
			}

			kvs[idx+1].Value = v

			idx += 2
		}

		tuples := make([][2]*ast.Term, len(kvs)/2)
		for i := 0; i < len(kvs); i += 2 {
			tuples[i/2] = *(*[2]*ast.Term)(kvs[i : i+2])
		}

		return ast.NewObject(tuples...), nil
	default:
		panic(fmt.Sprintf("%v: unsupported type: %T", x, x))
	}
}
