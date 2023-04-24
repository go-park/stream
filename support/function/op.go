package function

import "golang.org/x/exp/constraints"

func Sum[T constraints.Ordered](i, j T) T {
	return i + j
}

func Max[T constraints.Ordered](i, j T) T {
	if i >= j {
		return i
	}
	return j
}

func Min[T constraints.Ordered](i, j T) T {
	if i <= j {
		return i
	}
	return j
}
