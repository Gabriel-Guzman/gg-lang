package tokenizer

import (
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	uni "unicode"
)

type TokenType int

const (
	beginOperators TokenType = 1 << iota
	RPlus
	RMinus
	RMul
	RDiv
	RAssign
	endOperators

	beginContainers
	ROpenParen
	RCloseParen
	ROpenBrace
	RCloseBrace
	endContainers

	beginSeparators
	RTerm
	RComma
	RSpace
	endSeparators

	beginIdentifiers
	Var
	NumberLiteral
	StringLiteral
	endIdentifiers

	beginKeywords
	Function
	endKeywords
)

func (t TokenType) IsOperator() bool {
	return t > beginOperators && t < endOperators
}
func (t TokenType) IsContainer() bool {
	return t > beginContainers && t < endContainers
}
func (t TokenType) IsSeparator() bool {
	return t > beginSeparators && t < endSeparators
}
func (t TokenType) IsIdentifier() bool {
	return t > beginIdentifiers && t < endIdentifiers
}

type token string

var reservedTokens = map[TokenType]token{
	RPlus:   "+",
	RMinus:  "-",
	RMul:    "*",
	RDiv:    "/",
	RTerm:   ";",
	RAssign: "=",

	ROpenParen:  "(",
	RCloseParen: ")",
	ROpenBrace:  "{",
	RCloseBrace: "}",
	RComma:      ",",
	RSpace:      " ",
	Function:    "routine",
}

var reservedTokensMap = map[string]TokenType{}

func init() {
	for i, c := range reservedTokens {
		reservedTokensMap[string(c)] = i
	}
}

func isReserved(in string) bool {
	_, ok := reservedTokensMap[in]
	return ok
}

func lookup(in string) TokenType {
	return reservedTokensMap[in]
}

type Token struct {
	Start     int
	End       int
	Str       string
	TokenType TokenType
}

func (t Token) String() string {
	return fmt.Sprintf("(%d-%d) %s", t.Start, t.End, t.Str)
}

func TokenizeRunes(ins []rune) ([][]Token, error) {
	iter := iterator.New(ins)
	iter.Stringer = func(in rune) string {
		return string(in)
	}
	iter.Separator = ""

	var stmts [][]Token
	var stmt []Token
	curr, ok := iter.Next()

	for ok {
	sw_stmt:
		switch {
		case uni.IsSpace(curr):
		case isReserved(string(curr)):
			opt := lookup(string(curr))
			if opt == RTerm {
				stmts = append(stmts, stmt)
				stmt = nil
				break sw_stmt
			}

			stmt = append(stmt, newToken(iter, opt))
			if opt == ROpenBrace {
				stmts = append(stmts, stmt)
				stmt = nil
				break sw_stmt
			}
		case uni.IsLetter(curr):
			start := iter.Index()
			vr := variable(iter)
			end := iter.Index()
			if isReserved(vr.Str) {
				opt := lookup(vr.Str)
				stmt = append(stmt, Token{
					Start:     start,
					End:       end,
					Str:       vr.Str,
					TokenType: opt,
				})
			} else {
				stmt = append(stmt, vr)
			}
		case uni.IsDigit(curr):
			tok, err := numLiteral(iter)
			if err != nil {
				return nil, err
			}
			stmt = append(stmt, tok)
		case curr == '"':
			_, ok = iter.Next() // consume the opening quote
			if !ok {
				return nil, ggErrs.Runtime("unterminated string literal")
			}
			tok, err := stringLiteral(iter)
			if err != nil {
				return nil, ggErrs.Runtime("invalid character %c in string literal", curr)
			}
			stmt = append(stmt, tok)
		default:
			return nil, ggErrs.Runtime("unexpected character %c, %d", curr, curr)
		}
		curr, ok = iter.Next()
	}

	if len(stmt) > 0 {
		return nil, ggErrs.Runtime("unterminated statement\n%s\n%s", stmt, iter.String())
	}
	return stmts, nil
}

func newToken(iter *iterator.Iter[rune], tokenType TokenType) Token {
	start := iter.Index()
	ret := Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(iter.Current()),
		TokenType: tokenType,
	}
	return ret
}

func stringLiteral(iter *iterator.Iter[rune]) (Token, error) {
	start := iter.Index()
	str := []rune{iter.Current()}

	curr, ok := iter.Next()
	for ok && curr != '"' {
		str = append(str, curr)
		curr, ok = iter.Next()
	}

	if !ok {
		return Token{}, ggErrs.Runtime("unexpected end of input after string literal")
	}

	ret := Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(str),
		TokenType: StringLiteral,
	}

	return ret, nil
}

func numLiteral(iter *iterator.Iter[rune]) (Token, error) {
	start := iter.Index()
	num := []rune{iter.Current()}
	next, ok := iter.Peek()

loop:
	for ok {
		switch {
		case uni.IsDigit(next):
			num = append(num, next)
			_, ok = iter.Next() // consume the next rune
			next, ok = iter.Peek()
		case uni.IsLetter(next):
			return Token{}, ggErrs.Runtime("unexpected character %c after number literal", next)
		default:
			break loop
		}
	}

	return Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(num),
		TokenType: NumberLiteral,
	}, nil
}

// assume rune is not the first rune of the identifier
func idRune(r rune) bool {
	return uni.IsLetter(r) || uni.IsDigit(r) || r == '_'
}

func variable(iter *iterator.Iter[rune]) Token {
	start := iter.Index()
	id := []rune{iter.Current()}

	next, ok := iter.Peek()
	for ok {
		if !idRune(next) {
			break
		}
		id = append(id, next)

		_, _ = iter.Next() // consume the next rune we just checked
		next, ok = iter.Peek()
	}

	return Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(id),
		TokenType: Var,
	}
}
