package set

type Set[T comparable] struct {
	set map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		set: make(map[T]bool),
	}
}

func (s Set[T]) Has(v T) bool {
	_, ok := s.set[v]
	return ok
}

func (s Set[T]) Put(v T) bool {
	if s.Has(v) {
		return false
	}
	s.set[v] = true
	return true
}

func (s Set[T]) Size() int {
	return len(s.set)
}

func (s Set[T]) Iterator() map[T]bool {
	return s.set
}
