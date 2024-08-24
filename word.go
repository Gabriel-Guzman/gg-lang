package main

type token int

const (
	VAR token = iota
	OPERATOR
	NUMBER_LITERAL
)

type word struct {
	str  string
	role token
}

type wordIter struct {
	words []word
	curr  int
}

func newWordIter(words []word) *wordIter {
	return &wordIter{words: words, curr: -1}
}

func (wi *wordIter) Copy() *wordIter {
	newWords := make([]word, len(wi.words))
	copy(newWords, wi.words)
	return &wordIter{words: newWords, curr: wi.curr}
}

func (wi *wordIter) Current() (word, bool) {
	if wi.curr < 0 || wi.curr >= len(wi.words) {
		return word{}, false
	}
	return wi.words[wi.curr], true
}

func (wi *wordIter) Next() (word, bool) {
	wi.curr++ // Move to the next word in the slice.
	if wi.curr >= len(wi.words) {
		return word{}, false
	}

	w := wi.words[wi.curr]
	return w, true
}

func (wi *wordIter) Reset() {
	wi.curr = -1
}

func (wi *wordIter) Peek() (word, bool) {
	if (wi.curr + 1) >= len(wi.words) {
		return word{}, false
	}
	return wi.words[wi.curr+1], true
}

func (wi *wordIter) Prev() (word, bool) {
	if wi.curr <= 0 {
		return word{}, false
	}
	return (wi.words)[wi.curr-1], true
}
