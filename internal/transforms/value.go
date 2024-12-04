//nolint:gochecknoglobals
package transforms

import (
	"fmt"
	"strconv"

	"github.com/anderseknert/roast/pkg/util"
	"github.com/open-policy-agent/opa/ast"
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

	minusOneValue = [1]ast.Value{ast.Number("-1")}

	intToValue = [...]ast.Value{
		ast.Number("0"),
		ast.Number("1"),
		ast.Number("2"),
		ast.Number("3"),
		ast.Number("4"),
		ast.Number("5"),
		ast.Number("6"),
		ast.Number("7"),
		ast.Number("8"),
		ast.Number("9"),
		ast.Number("10"),
		ast.Number("11"),
		ast.Number("12"),
		ast.Number("13"),
		ast.Number("14"),
		ast.Number("15"),
		ast.Number("16"),
		ast.Number("17"),
		ast.Number("18"),
		ast.Number("19"),
		ast.Number("20"),
		ast.Number("21"),
		ast.Number("22"),
		ast.Number("23"),
		ast.Number("24"),
		ast.Number("25"),
		ast.Number("26"),
		ast.Number("27"),
		ast.Number("28"),
		ast.Number("29"),
		ast.Number("30"),
		ast.Number("31"),
		ast.Number("32"),
		ast.Number("33"),
		ast.Number("34"),
		ast.Number("35"),
		ast.Number("36"),
		ast.Number("37"),
		ast.Number("38"),
		ast.Number("39"),
		ast.Number("40"),
		ast.Number("41"),
		ast.Number("42"),
		ast.Number("43"),
		ast.Number("44"),
		ast.Number("45"),
		ast.Number("46"),
		ast.Number("47"),
		ast.Number("48"),
		ast.Number("49"),
		ast.Number("50"),
		ast.Number("51"),
		ast.Number("52"),
		ast.Number("53"),
		ast.Number("54"),
		ast.Number("55"),
		ast.Number("56"),
		ast.Number("57"),
		ast.Number("58"),
		ast.Number("59"),
		ast.Number("60"),
		ast.Number("61"),
		ast.Number("62"),
		ast.Number("63"),
		ast.Number("64"),
		ast.Number("65"),
		ast.Number("66"),
		ast.Number("67"),
		ast.Number("68"),
		ast.Number("69"),
		ast.Number("70"),
		ast.Number("71"),
		ast.Number("72"),
		ast.Number("73"),
		ast.Number("74"),
		ast.Number("75"),
		ast.Number("76"),
		ast.Number("77"),
		ast.Number("78"),
		ast.Number("79"),
		ast.Number("80"),
		ast.Number("81"),
		ast.Number("82"),
		ast.Number("83"),
		ast.Number("84"),
		ast.Number("85"),
		ast.Number("86"),
		ast.Number("87"),
		ast.Number("88"),
		ast.Number("89"),
		ast.Number("90"),
		ast.Number("91"),
		ast.Number("92"),
		ast.Number("93"),
		ast.Number("94"),
		ast.Number("95"),
		ast.Number("96"),
		ast.Number("97"),
		ast.Number("98"),
		ast.Number("99"),
		ast.Number("100"),
		ast.Number("101"),
		ast.Number("102"),
		ast.Number("103"),
		ast.Number("104"),
		ast.Number("105"),
		ast.Number("106"),
		ast.Number("107"),
		ast.Number("108"),
		ast.Number("109"),
		ast.Number("110"),
		ast.Number("111"),
		ast.Number("112"),
		ast.Number("113"),
		ast.Number("114"),
		ast.Number("115"),
		ast.Number("116"),
		ast.Number("117"),
		ast.Number("118"),
		ast.Number("119"),
		ast.Number("120"),
		ast.Number("121"),
		ast.Number("122"),
		ast.Number("123"),
		ast.Number("124"),
		ast.Number("125"),
		ast.Number("126"),
		ast.Number("127"),
	}

	booleanValues = [2]ast.Value{
		ast.Boolean(false),
		ast.Boolean(true),
	}

	nullValue = [1]ast.Value{
		ast.Null{},
	}

	roastKeyValues map[string]ast.Value

	roastKeyTerms map[string]*ast.Term

	commonStringValues map[string]ast.Value

	commonStringTerms map[string]*ast.Term
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
		if x {
			return booleanValues[1], nil
		}

		return booleanValues[0], nil
	case float64:
		ix := int(x)
		if x == float64(ix) {
			if ix == -1 {
				return minusOneValue[0], nil
			}

			if ix >= 0 && ix < len(intToValue) {
				return intToValue[ix], nil
			}
		}

		return ast.Number(strconv.FormatFloat(x, 'g', -1, 64)), nil
	case string:
		if s, ok := commonStringValues[x]; ok {
			return s, nil
		}

		return ast.String(x), nil
	case []string:
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

		r := ast.NewObject(tuples...)

		return r, nil
	default:
		panic(fmt.Sprintf("%v: unsupported type: %T", x, x))
	}
}
