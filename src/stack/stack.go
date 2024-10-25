package stack

// a stack implementation using a singly-linked list.
type Stack[T any] struct {
	curr *elem[T]
}

type elem[T any] struct {
	parent *elem[T]
	value  T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(val T) {
	newElem := &elem[T]{parent: s.curr, value: val}
	s.curr = newElem
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.curr == nil {
		var ret T
		return ret, false
	}

	val := s.curr.value
	s.curr = s.curr.parent
	return val, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if s.curr == nil {
		var ret T
		return ret, false
	}

	return s.curr.value, true
}
