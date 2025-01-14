package util

import (
	"slices"
	"unsafe"
)

// Allocation free conversion from []byte to string (unsafe)
// Note that the byte slice must not be modified after conversion.
func ByteSliceToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

// Allocation free conversion from string to []byte (unsafe).
func StringToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// NewPtrSlice allocates a slice of pointers to T, with only 2
// allocations performed regardless of the size of n.
func NewPtrSlice[T any](n int) []*T {
	s := make([]*T, n)
	p := make([]T, n)

	for i := range s {
		s[i] = &p[i]
	}

	return s
}

// GrowPtrSlice grows a slice of pointers to T by n elements.
func GrowPtrSlice[T any](s []*T, n int) []*T {
	s = slices.Grow(s, n)
	p := make([]T, n)

	for i := range n {
		s = append(s, &p[i])
	}

	return s
}
