package utils

func Unique[T comparable](values []T) []T {
	set := NewSet[T]()

	for _, v := range values {
		set.Put(v)
	}

	res := make([]T, 0, set.Size())

	for k := range set.Iterator() {
		res = append(res, k)
	}
	return res
}
