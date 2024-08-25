package main

import (
	"fmt"
	"unicode"
)

type tokenizer struct {
	stmts [][]token

	_ctr     int
	_line    []token
	_curWord []rune

	_curWordRole tokenType
}

func BuildTokens(ins []rune) ([][]token, error) {
	toke := &tokenizer{}
	err := toke.fromRunes(ins)

	return toke.stmts, err
}

func (t *tokenizer) fromRunes(ins []rune) error {

	t._ctr = -1

	t._curWord = make([]rune, 0)

	runeIter := newIter[rune](ins)
	for {
		r, ok := runeIter.Next()
		if !ok {
			break
		}

		switch {
		case unicode.IsSpace(r):
			t.wordDone()
		case r == ';':
			t.wordDone()
			t.lineDone()
		case isReservedRune(r):
			t.wordDone()
			t._curWordRole = OPERATOR
			t.addToWord(r)
			t.wordDone()
		case unicode.IsLetter(r):
			if len(t._curWord) != 0 && t._curWordRole == NUMBER_LITERAL {
				return fmt.Errorf("invalid character %s in number literal %s", string(r), string(t._curWord))
			}
			t._curWordRole = VAR
			t.addToWord(r)
		case unicode.IsDigit(r):
			if len(t._curWord) == 0 {
				t._curWordRole = NUMBER_LITERAL
			}
			t.addToWord(r)
		case r == '"':
			if len(t._curWord) > 0 {
				return fmt.Errorf("unexpected character %c", r)
			}
			strMember, ok := runeIter.Next()

		stringSearch:
			for ok {
				if strMember == '"' {
					t._curWordRole = STRING_LITERAL
					t.wordDone()
					break stringSearch
				}

				t.addToWord(strMember)
				strMember, ok = runeIter.Next()
			}
			if !ok {
				return fmt.Errorf("unterminated string %s", string(t._curWord))
			}
		}
	}

	return nil
}

func (t *tokenizer) addToWord(r rune) {
	t._curWord = append(t._curWord, r)
}

func (t *tokenizer) wordDone() {
	if len(t._curWord) == 0 {
		return
	}

	t._line = append(t._line, token{
		start:     t._ctr - len(t._curWord),
		end:       t._ctr,
		str:       string(t._curWord),
		tokenType: t._curWordRole,
	})
	t._curWord = nil
}

func (t *tokenizer) lineDone() {
	if len(t._line) == 0 {
		return
	}
	t.stmts = append(t.stmts, t._line)
	t._line = nil
}

type opType string

const (
	R_PLUS  opType = "+"
	R_MINUS        = "-"
	R_MUL          = "*"
	R_DIV          = "/"
)

func isReservedRune(word rune) bool {
	strWord := opType(word)
	return strWord == R_PLUS ||
		strWord == R_MINUS ||
		strWord == R_MUL ||
		strWord == R_DIV
}

//func handleSpace(_curWord []rune, _line *[]string) {
//	wordDone(_curWord, _line)
//}
//
//func (T *tokenizer) handleSemicolon(_curWord []rune, _line *[]string) {
//	wordDone(_curWord, _line)
//	_curWord = make([]rune, 0)
//}
