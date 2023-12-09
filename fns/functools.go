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

func Filter[T any](s []T, fn func(t T) bool) []T {
	var arr []T
	for _, t := range s {
		if fn(t) {
			arr = append(arr, t)
		}
	}
	return arr
}
