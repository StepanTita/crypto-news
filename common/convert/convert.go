package convert

func FromPtr[T any](v *T) T {
	if v == nil {
		var res T
		return res
	}
	return *v
}

func ToPtr[T any](v T) *T {
	return &v
}
