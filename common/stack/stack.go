package stack

type Stack[T any] struct {
	container []T
	last      int
}

func New[T any]() *Stack[T] {
	return &Stack[T]{
		container: make([]T, 10),
		last:      -1,
	}
}

func (s *Stack[T]) Empty() bool {
	return s.last < 0
}

func (s *Stack[T]) Push(val T) {
	if len(s.container) > s.last {
		s.last++
		s.container[s.last] = val
		return
	}

	s.container = append(s.container, val)
}

func (s *Stack[T]) Top() T {
	return s.container[s.last]
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.Empty() {
		var zero T
		return zero, false
	}
	val := s.container[s.last]
	s.last--
	return val, true
}
