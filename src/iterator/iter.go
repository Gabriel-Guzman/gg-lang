package iterator

type Iter[T any] struct {
	members []T
	curr    int
}

func NewIter[T any](words []T) *Iter[T] {
	return &Iter[T]{members: words, curr: -1}
}

func (wi *Iter[T]) Copy() *Iter[T] {
	newWords := make([]T, len(wi.members))
	copy(newWords, wi.members)
	return &Iter[T]{members: newWords, curr: wi.curr}
}

func (wi *Iter[T]) Current() (T, bool) {
	if wi.curr < 0 || wi.curr >= len(wi.members) {
		var ret T
		return ret, false
	}
	return wi.members[wi.curr], true
}

func (wi *Iter[T]) Next() (T, bool) {
	wi.curr++ // Move to the next token in the slice.
	if wi.curr >= len(wi.members) {
		var ret T
		return ret, false
	}

	w := wi.members[wi.curr]
	return w, true
}

func (wi *Iter[T]) Reset() {
	wi.curr = -1
}

func (wi *Iter[T]) Peek() (T, bool) {
	if (wi.curr + 1) >= len(wi.members) {
		var ret T
		return ret, false
	}
	return wi.members[wi.curr+1], true
}

func (wi *Iter[T]) Prev() (T, bool) {
	if wi.curr <= 0 {
		var ret T
		return ret, false
	}
	return (wi.members)[wi.curr-1], true
}
