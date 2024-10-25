package parser

import (
	"fmt"
	"strings"
)

type ItemPrinter[T any] func(in T) string
type Parser[T any] struct {
	items []T
	curr  int

	stringer  ItemPrinter[T]
	separator string

	TruncBefore  int
	TruncAfter   int
	WrapSelected bool

	Curr    T
	HasCurr bool

	Next    T
	HasNext bool
}

func (p *Parser[T]) SetSeparator(sep string) {
	p.separator = sep
}

func (p *Parser[T]) SetStringer(s ItemPrinter[T]) {
	p.stringer = s
}

func (p *Parser[T]) IsDone() bool {
	return !p.HasCurr
}

// returns the range of items from TruncBefore inclusive to TruncAfter exclusive
func (p *Parser[T]) truncate() []T {
	lower := max(p.curr-p.TruncBefore, 0)
	upper := min(p.curr+p.TruncAfter, len(p.items)-1)
	return p.items[lower:upper]
}

func (p *Parser[T]) String() string {
	sb := &strings.Builder{}
	for i, item := range p.truncate() {
		sb.WriteString(p.stringer(item))
		if i < len(p.items)-1 {
			sb.WriteString(p.separator)
		}
	}
	return sb.String()
}

func New[T any](items []T) *Parser[T] {
	ret := &Parser[T]{
		curr:         -1,
		items:        items,
		WrapSelected: true,
		TruncBefore:  5,
		TruncAfter:   2,
		stringer:     func(in T) string { return fmt.Sprintf("%v", in) },
		separator:    ",",
	}

	ret.Advance()
	return ret
}

func (p *Parser[T]) Index() int {
	return p.curr
}

func (p *Parser[T]) Advance() {
	p.curr++
	p.HasCurr = p.curr < len(p.items) && p.curr >= 0
	p.HasNext = p.curr+1 < len(p.items) && p.curr+1 >= 0
	if p.HasCurr {
		p.Curr = p.items[p.curr]
	}
	if p.HasNext {
		p.Next = p.items[p.curr+1]
	}
}
