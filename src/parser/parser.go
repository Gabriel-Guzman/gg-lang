package parser

import (
	"gg-lang/src/iterator"
)

type Parser[T any] struct {
	iter *iterator.Iter[T]

	Curr    T
	HasCurr bool

	Next    T
	HasNext bool
}

func (p *Parser[T]) SetSeparator(sep string) {
	p.iter.Separator = sep
}

func (p *Parser[T]) SetStringer(s func(in T) string) {
	p.iter.Stringer = s
}

func (p *Parser[T]) String() string {
	return p.iter.String()
}

func New[T any](items []T) *Parser[T] {
	iter := iterator.New(items)
	curr, hasCurr := iter.Next()
	next, hasNext := iter.Peek()

	return &Parser[T]{
		iter:    iterator.New(items),
		Curr:    curr,
		HasCurr: hasCurr,
		Next:    next,
		HasNext: hasNext,
	}
}

func (p *Parser[T]) Index() int {
	return p.iter.Index()
}

func (p *Parser[T]) Advance() {
	curr, hasCurr := p.iter.Next()
	next, hasNext := p.iter.Peek()
	p.Curr = curr
	p.HasCurr = hasCurr
	p.Next = next
	p.HasNext = hasNext
}
