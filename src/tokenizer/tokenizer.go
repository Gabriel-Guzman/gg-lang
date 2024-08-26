//go:generate stringer -type=Token
package tokenizer

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/iterator"
	"unicode"
)

type tokenType int

const (
	VAR tokenType = iota
	Operator
	NumberLiteral
	StringLiteral
)

type Token struct {
	Start     int
	End       int
	Str       string
	TokenType tokenType
}

type tokenizer struct {
	stmts [][]Token

	_iter    *iterator.Iter[rune]
	_line    []Token
	_curWord []rune

	_curWordRole tokenType
}

func BuildStmts(ins []rune) ([][]Token, error) {
	toke := &tokenizer{}
	err := toke.fromRunes(ins)

	return toke.stmts, err
}

func (t *tokenizer) fromRunes(ins []rune) error {
	t._curWord = make([]rune, 0)
	t._iter = iterator.New[rune](ins)

	for {
		r, ok := t._iter.Next()
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
			t._curWordRole = Operator
			t.addToWord(r)
			t.wordDone()
		case unicode.IsLetter(r):
			if len(t._curWord) != 0 && t._curWordRole == NumberLiteral {
				return fmt.Errorf("invalid character %s in number literal %s", string(r), string(t._curWord))
			}
			t._curWordRole = VAR
			t.addToWord(r)
		case unicode.IsDigit(r):
			if len(t._curWord) == 0 {
				t._curWordRole = NumberLiteral
			}
			t.addToWord(r)
		case r == '"':
			if len(t._curWord) > 0 {
				return fmt.Errorf("unexpected character %c", r)
			}
			strMember, ok := t._iter.Next()

		stringSearch:
			for ok {
				if strMember == '"' {
					t._curWordRole = StringLiteral
					t.wordDone()
					break stringSearch
				}

				t.addToWord(strMember)
				strMember, ok = t._iter.Next()
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

	t._line = append(t._line, Token{
		Start:     t._iter.Index() - len(t._curWord),
		End:       t._iter.Index(),
		Str:       string(t._curWord),
		TokenType: t._curWordRole,
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
	RPlus  opType = "+"
	RMinus        = "-"
	RMul          = "*"
	RDiv          = "/"
)

func isReservedRune(word rune) bool {
	strWord := opType(word)
	return strWord == RPlus ||
		strWord == RMinus ||
		strWord == RMul ||
		strWord == RDiv
}
