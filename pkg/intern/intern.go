package intern

import (
	"github.com/open-policy-agent/opa/v1/ast"

	"github.com/styrainc/roast/internal/intern"
)

// TODO: move to OPA.
var EmptyArray = ast.NewArray()

func StringTerm(s string) *ast.Term {
	if t, ok := intern.StringTerms[s]; ok {
		return t
	}

	return ast.StringTerm(s)
}

func StringValue(s string) ast.Value {
	if v, ok := intern.StringValues[s]; ok {
		return v
	}

	return ast.String(s)
}
