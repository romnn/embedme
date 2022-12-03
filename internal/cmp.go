package internal

import (
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a T, b T) T {
	if a > b {
		return b
	}
	return a
}
