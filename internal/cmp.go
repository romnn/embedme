package internal

import (
	"golang.org/x/exp/constraints"
)

// Max returns the maximum of two values
func Max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two values
func Min[T constraints.Ordered](a T, b T) T {
	if a > b {
		return b
	}
	return a
}
