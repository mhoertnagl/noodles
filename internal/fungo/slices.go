package fungo

type T = interface{}

func Apply(xs []T, f func(x T) T) []T {
	rs := make([]T, len(xs))
	for i, x := range xs {
		rs[i] = f(x)
	}
	return rs
}

func IApply(xs []T, f func(i int, x T) T) []T {
	rs := make([]T, len(xs))
	for i, x := range xs {
		rs[i] = f(i, x)
	}
	return rs
}
