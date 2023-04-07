package utils

import (
	"regexp"

	"common/convert"
)

func Map[T, K any](arr []T, f func(T) K) []K {
	res := make([]K, len(arr))
	for i, a := range arr {
		res[i] = f(a)
	}
	return res
}

func Filter[T any](arr []T, pred func(T) bool) []T {
	res := make([]T, 0, 10)
	for _, a := range arr {
		if pred(a) {
			res = append(res, a)
		}
	}
	return res
}

const regex = `<.*?>`

// StripHtmlRegex This method uses a regular expresion to remove HTML tags.
func StripHtmlRegex(s *string) *string {
	if s == nil {
		return nil
	}
	r := regexp.MustCompile(regex)
	return convert.ToPtr(r.ReplaceAllString(*s, ""))
}
