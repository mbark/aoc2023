package fns

func Every[T any](s []T, fn func(i int) bool) bool {
	for i := 0; i < len(s); i++ {
		if !fn(i) {
			return false
		}
	}

	return true
}

func Some[T any](s []T, fn func(i int) bool) bool {
	for i := 0; i < len(s); i++ {
		if fn(i) {
			return true
		}
	}

	return false
}

func Map[T, R any](s []T, fn func(t T) R) []R {
	arr := make([]R, len(s))
	for i, t := range s {
		arr[i] = fn(t)
	}
	return arr
}

func Associate[V any, K comparable](s []V, fn func(t V) K) map[K]V {
	m := make(map[K]V, len(s))
	for _, t := range s {
		m[fn(t)] = t
	}
	return m
}

func Find[T any](s []T, fn func(t T) bool) (T, bool) {
	for _, t := range s {
		if fn(t) {
			return t, true
		}
	}

	var d T
	return d, false
}

func Filter[T any](s []T, fn func(t T) bool) []T {
	var arr []T
	for _, t := range s {
		if fn(t) {
			arr = append(arr, t)
		}
	}
	return arr
}

func Flatten[T any](ts [][]T) []T {
	length := 0
	for i := range ts {
		length += len(ts[i])
	}

	arr := make([]T, 0, length)
	for i := range ts {
		arr = append(arr, ts[i]...)
	}

	return arr
}