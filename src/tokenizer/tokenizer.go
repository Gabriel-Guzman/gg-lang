//go:generate stringer -type=tokenType
package tokenizer

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/iterator"
	"unicode"
)

type tokenType int

const (
	Var tokenType = iota
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

	iter    *iterator.Iter[rune]
	stmt    []Token
	curWord []rune

	curWordRole tokenType
}

func BuildStmts(ins []rune) ([][]Token, error) {
	toke := &tokenizer{}
	err := toke.fromRunes(ins)

	return toke.stmts, err
}

func (t *tokenizer) fromRunes(ins []rune) error {
	t.curWord = make([]rune, 0)
	t.iter = iterator.New[rune](ins)

	for {
		r, ok := t.iter.Next()
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
			t.curWordRole = Operator
			t.addToWord(r)
			t.wordDone()
		case unicode.IsLetter(r):
			if len(t.curWord) != 0 && t.curWordRole == NumberLiteral {
				return fmt.Errorf("invalid character %s in number literal %s", string(r), string(t.curWord))
			}
			t.curWordRole = Var
			t.addToWord(r)
		case unicode.IsDigit(r):
			if len(t.curWord) == 0 {
				t.curWordRole = NumberLiteral
			}
			t.addToWord(r)
		case r == '"':
			if len(t.curWord) > 0 {
				return fmt.Errorf("unexpected character %c", r)
			}
			strMember, ok := t.iter.Next()

		stringSearch:
			for ok {
				if strMember == '"' {
					t.curWordRole = StringLiteral
					t.wordDone()
					break stringSearch
				}

				t.addToWord(strMember)
				strMember, ok = t.iter.Next()
			}
			if !ok {
				return fmt.Errorf("unterminated string %s", string(t.curWord))
			}
		}
	}

	return nil
}

func (t *tokenizer) addToWord(r rune) {
	t.curWord = append(t.curWord, r)
}

func (t *tokenizer) wordDone() {
	if len(t.curWord) == 0 {
		return
	}

	t.stmt = append(t.stmt, Token{
		Start:     t.iter.Index() - len(t.curWord),
		End:       t.iter.Index(),
		Str:       string(t.curWord),
		TokenType: t.curWordRole,
	})
	t.curWord = nil
}

func (t *tokenizer) lineDone() {
	if len(t.stmt) == 0 {
		return
	}
	t.stmts = append(t.stmts, t.stmt)
	t.stmt = nil
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
