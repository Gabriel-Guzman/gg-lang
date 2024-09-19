package iterator

import (
	"fmt"
	"strings"
)

type Iter[T any] struct {
	members   []T
	curr      int
	reverse   bool
	Stringer  func(T) string
	Separator string
}

func New[T any](words []T) *Iter[T] {
	return &Iter[T]{members: words, curr: -1, Separator: " "}
}

func (wi *Iter[T]) SetIndex(index int) {
	wi.curr = index
}

func (wi *Iter[T]) Index() int {
	return wi.curr
}

func (wi *Iter[T]) Reverse() *Iter[T] {
	reversed := wi.Copy()
	reversed.reverse = !wi.reverse
	return reversed
}

// modifying the members in the copy will affect the original iterator members
func (wi *Iter[T]) Copy() *Iter[T] {
	return &Iter[T]{
		members:   wi.members,
		curr:      wi.curr,
		Stringer:  wi.Stringer,
		Separator: wi.Separator,
		reverse:   wi.reverse,
	}
}

func (wi *Iter[T]) nextIndex() int {
	if wi.reverse {
		return wi.curr - 1
	} else {
		return wi.curr + 1
	}
}

func (wi *Iter[T]) String() string {
	var out []string
	if wi.curr == -1 {
		out = append(out, ">><<")
	}
	for i, w := range wi.members {
		var str string
		if wi.Stringer != nil {
			str = wi.Stringer(w)
		} else {
			str = fmt.Sprintf("%+v", w)
		}
		if i == wi.curr {
			out = append(out, fmt.Sprintf(">>%s<<", str))
			continue
		}
		out = append(out, str)
	}

	if wi.curr != -1 && !wi.HasCurrent() {
		out = append(out, ">><<")
	}
	done := strings.Join(out, wi.Separator)
	return done
}

func (wi *Iter[T]) Current() T {
	return wi.members[wi.curr]
}
func (wi *Iter[T]) HasCurrent() bool {
	return wi.hasIndex(wi.curr)
}
func (wi *Iter[T]) hasIndex(index int) bool {
	return index >= 0 && index < len(wi.members)
}

func (wi *Iter[T]) Next() (T, bool) {
	wi.curr = wi.nextIndex()
	if !wi.hasIndex(wi.curr) {
		var ret T
		return ret, false
	}

	w := wi.members[wi.curr]
	return w, true
}

func (wi *Iter[T]) HasNext() bool {
	return wi.hasIndex(wi.nextIndex())
}

func (wi *Iter[T]) End() {
	if wi.reverse {
		wi.curr = -1
		return
	}
	wi.curr = len(wi.members)
}

func (wi *Iter[T]) Reset() {
	if wi.reverse {
		wi.curr = len(wi.members)
		return
	}
	wi.curr = -1
}

func (wi *Iter[T]) Peek() (T, bool) {
	if !wi.HasNext() {
		var ret T
		return ret, false
	}
	return wi.members[wi.nextIndex()], true
}

func (wi *Iter[T]) Prev() (T, bool) {
	var prevIndex int
	if wi.reverse {
		prevIndex = wi.curr + 1
	} else {
		prevIndex = wi.curr - 1
	}
	if !wi.hasIndex(prevIndex) {
		var ret T
		return ret, false
	}
	return wi.members[prevIndex], true
}

func (wi *Iter[T]) Len() int {
	return len(wi.members)
}
