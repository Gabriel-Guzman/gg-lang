package iterator

import (
	"fmt"
	"strings"
)

type Iter[T any] struct {
	members []T
	curr    int
}

func New[T any](words []T) *Iter[T] {
	return &Iter[T]{members: words, curr: -1}
}

func (wi *Iter[T]) Index() int {
	return wi.curr
}

func (wi *Iter[T]) Copy() *Iter[T] {
	newWords := make([]T, len(wi.members))
	copy(newWords, wi.members)
	return &Iter[T]{members: newWords, curr: wi.curr}
}

func (wi *Iter[T]) String() string {
	var out []string
	for i, w := range wi.members {
		if i == wi.curr {
			out = append(out, fmt.Sprintf("[[%+v]]", w))
			continue
		}
		out = append(out, fmt.Sprintf("%+v", w))
	}

	done := strings.Join(out, ", ")
	return done
}

func (wi *Iter[T]) Current() T {
	return wi.members[wi.curr]
}
func (wi *Iter[T]) HasCurrent() bool {
	return wi.curr >= 0 && wi.curr < len(wi.members)
}

func (wi *Iter[T]) Next() (T, bool) {
	wi.curr++
	if wi.curr >= len(wi.members) {
		var ret T
		return ret, false
	}

	w := wi.members[wi.curr]
	return w, true
}

func (wi *Iter[T]) HasNext() bool {
	return (wi.curr + 1) < len(wi.members)
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
