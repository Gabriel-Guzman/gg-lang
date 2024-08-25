package main

type iter[T any] struct {
	members []T
	curr    int
}

func newIter[T any](words []T) *iter[T] {
	return &iter[T]{members: words, curr: -1}
}

func (wi *iter[T]) Copy() *iter[T] {
	newWords := make([]T, len(wi.members))
	copy(newWords, wi.members)
	return &iter[T]{members: newWords, curr: wi.curr}
}

func (wi *iter[T]) Current() (T, bool) {
	if wi.curr < 0 || wi.curr >= len(wi.members) {
		var ret T
		return ret, false
	}
	return wi.members[wi.curr], true
}

func (wi *iter[T]) Next() (T, bool) {
	wi.curr++ // Move to the next token in the slice.
	if wi.curr >= len(wi.members) {
		var ret T
		return ret, false
	}

	w := wi.members[wi.curr]
	return w, true
}

func (wi *iter[T]) Reset() {
	wi.curr = -1
}

func (wi *iter[T]) Peek() (T, bool) {
	if (wi.curr + 1) >= len(wi.members) {
		var ret T
		return ret, false
	}
	return wi.members[wi.curr+1], true
}

func (wi *iter[T]) Prev() (T, bool) {
	if wi.curr <= 0 {
		var ret T
		return ret, false
	}
	return (wi.members)[wi.curr-1], true
}
