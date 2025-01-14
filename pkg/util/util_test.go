package util

import (
	"testing"

	"github.com/open-policy-agent/opa/v1/ast"
)

func BenchmarkStringRepeatNewPtrSlice(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = stringRepeatNewPtrSlice("test", 1000)
	}
}

func BenchmarkStringRepeatMake(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = stringRepeatMake("test", 1000)
	}
}

func stringRepeatNewPtrSlice(s string, n int) []*ast.Term {
	sl := NewPtrSlice[ast.Term](n)
	for i := range s {
		sl[i].Value = ast.String("test")
	}

	return sl
}

func stringRepeatMake(s string, n int) []*ast.Term {
	sl := make([]*ast.Term, n)
	for i := range s {
		sl[i] = &ast.Term{Value: ast.String("test")}
	}

	return sl
}
