package intern

import (
	"github.com/open-policy-agent/opa/v1/ast"

	"github.com/styrainc/roast/internal/intern"
)

func StringTerm(s string) *ast.Term {
	if t, ok := intern.StringTerms[s]; ok {
		return t
	}

	return ast.InternedStringTerm(s)
}

func StringValue(s string) ast.Value {
	if v, ok := intern.StringValues[s]; ok {
		return v
	}

	return ast.String(s)
}
