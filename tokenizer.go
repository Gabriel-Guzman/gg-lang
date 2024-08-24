package main

import (
	"errors"
	"unicode"
)

type tokenizer struct {
	stmts [][]word

	_ctr     int
	_line    []word
	_curWord []rune

	_curWordRole token
}

func TokenizeRunes(ins []rune) ([][]word, error) {
	toke := &tokenizer{}
	err := toke.fromRunes(ins)

	return toke.stmts, err
}

func (t *tokenizer) fromRunes(ins []rune) error {

	t._ctr = -1

	t._curWord = make([]rune, 0)

	for {
		t._ctr++
		if t._ctr >= len(ins) {
			break
		}

		r := ins[t._ctr]

		if unicode.IsSpace(r) {
			t.wordDone()
			continue
		}

		if r == ';' {
			t.wordDone()
			t.lineDone()
			continue
		}

		if isReservedRune(r) {
			t.wordDone()
			t._curWordRole = OPERATOR
			t.addToWord(r)
			t.wordDone()
			continue
		}

		if unicode.IsLetter(r) {
			// can't have a number then a letter
			if len(t._curWord) != 0 && t._curWordRole == NUMBER_LITERAL {
				return errors.New("invalid character in number literal")
			}
			t._curWordRole = VAR
		}

		if unicode.IsDigit(r) {
			if len(t._curWord) == 0 {
				t._curWordRole = NUMBER_LITERAL
			}
		}

		t.addToWord(r)
	}

	t.finished()
	return nil
}

func (t *tokenizer) finished() {
	t._line = nil
	t._curWord = nil
	t._ctr = 0
}

func (t *tokenizer) addToWord(r rune) {
	t._curWord = append(t._curWord, r)
}

func (t *tokenizer) wordDone() {
	if len(t._curWord) == 0 {
		return
	}

	t._line = append(t._line, word{string(t._curWord), t._curWordRole})
	t._curWord = nil
}

func (t *tokenizer) lineDone() {
	if len(t._line) == 0 {
		return
	}
	t.stmts = append(t.stmts, t._line)
	t._line = nil
}

const (
	R_PLUS  rune = '+'
	R_MINUS      = '-'
	R_MUL        = '*'
	R_DIV        = '/'
)

func isReservedRune(word rune) bool {
	return word == R_PLUS ||
		word == R_MINUS ||
		word == R_MUL ||
		word == R_DIV
}

//func handleSpace(_curWord []rune, _line *[]string) {
//	wordDone(_curWord, _line)
//}
//
//func (t *tokenizer) handleSemicolon(_curWord []rune, _line *[]string) {
//	wordDone(_curWord, _line)
//	_curWord = make([]rune, 0)
//}
