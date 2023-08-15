package iteration

import "common/containers/set"

func Unique[T comparable](values []T) []T {
	s := set.NewSet[T]()

	for _, v := range values {
		s.Put(v)
	}

	res := make([]T, 0, s.Size())

	for k := range s.Iterator() {
		res = append(res, k)
	}
	return res
}
