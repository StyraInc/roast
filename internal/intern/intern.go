package intern

import (
	"github.com/open-policy-agent/opa/v1/ast"
)

var (
	// These are strings commonly found in the AST of all Rego policies,
	// like the name of built-in functions, keywords, etc.
	strings = [...]string{
		// Rego
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
		"v0",
		"v1",
		"v0v1",
		"unknown",

		// These are strings commonly found in linter policies, but
		// not necessarily anywhere else.
		"}",
		"# METADATA",
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

		// OPA / Roast keys
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
		"col",
		"row",
		"end",
		"package_path",
		"rule",
		"category",
		"aggregate_source",
		"aggregate_data",
		"rego_version",
		"negated_refs",
		"ast",
		"refs",
		"config",
		"lint",
		"collect",
	}

	StringTerms  = stringTermsMap(strings[:])
	StringValues = stringValuesMap(strings[:])
)

func stringTermsMap([]string) map[string]*ast.Term {
	m := make(map[string]*ast.Term, len(strings))

	for _, s := range strings {
		m[s] = ast.NewTerm(StringValues[s])
	}

	return m
}

func stringValuesMap([]string) map[string]ast.Value {
	m := make(map[string]ast.Value, len(strings))

	for _, s := range strings {
		m[s] = ast.String(s)
	}

	return m
}
