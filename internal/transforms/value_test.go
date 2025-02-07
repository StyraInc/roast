package transforms

import (
	"os"
	"testing"

	"github.com/anderseknert/roast/pkg/encoding"

	"github.com/open-policy-agent/opa/v1/ast"
)

func TestRoastAndOPAInterfaceToValueSameOutput(t *testing.T) {
	t.Parallel()

	inputMap := inputMap(t)

	roastValue, err := AnyToValue(inputMap)
	if err != nil {
		t.Fatal(err)
	}

	opaValue, err := ast.InterfaceToValue(inputMap)
	if err != nil {
		t.Fatal(err)
	}

	if roastValue.Compare(opaValue) != 0 {
		t.Fatal("values are not equal")
	}
}

// BenchmarkInterfaceToValue-10    	 741	   1615548 ns/op	 1376979 B/op	   24189 allocs/op
// ...
func BenchmarkInterfaceToValue(b *testing.B) {
	inputMap := inputMapB(b)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := AnyToValue(inputMap)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOPAInterfaceToValue-10    	616	   1942695 ns/op	 1566569 B/op	   45901 allocs/op
// BenchmarkOPAInterfaceToValue-10    	626	   1838247 ns/op	 1566848 B/op	   36037 allocs/op OPA 1.0
// ...
func BenchmarkOPAInterfaceToValue(b *testing.B) {
	inputMap := inputMapB(b)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := ast.InterfaceToValue(inputMap)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func inputMap(t *testing.T) map[string]any {
	t.Helper()

	bs, err := os.ReadFile("testdata/ast.rego")
	if err != nil {
		t.Fatal(err)
	}

	content := string(bs)

	module, err := ast.ParseModuleWithOpts("ast.rego", content, ast.ParserOptions{ProcessAnnotation: true})
	if err != nil {
		t.Fatal(err)
	}

	inputMap := make(map[string]any)

	if err := encoding.JSONRoundTrip(module, &inputMap); err != nil {
		t.Fatal(err)
	}

	return inputMap
}

func inputMapB(b *testing.B) map[string]any {
	b.Helper()

	bs, err := os.ReadFile("testdata/ast.rego")
	if err != nil {
		b.Fatal(err)
	}

	content := string(bs)

	module, err := ast.ParseModuleWithOpts("ast.rego", content, ast.ParserOptions{ProcessAnnotation: true})
	if err != nil {
		b.Fatal(err)
	}

	inputMap := make(map[string]any)

	if err := encoding.JSONRoundTrip(module, &inputMap); err != nil {
		b.Fatal(err)
	}

	return inputMap
}
