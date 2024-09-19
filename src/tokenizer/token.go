package tokenizer

import (
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	uni "unicode"
)

type TokenType int

const (
	beginOperators TokenType = iota
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
	RQuote
	endContainers

	beginSeparators
	RTerm
	RComma
	RSpace
	endSeparators

	beginIdentifiers
	Var
	IntLiteral
	StringLiteral
	TrueLiteral
	FalseLiteral
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
func (t TokenType) IsMathOperator() bool {
	return t == RPlus || t == RMinus || t == RMul || t == RDiv
}

func (t TokenType) String() string {
	if s, ok := reservedTokens[t]; ok {
		return s
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

var reservedTokens = map[TokenType]string{
	// operators
	RPlus:   "+",
	RMinus:  "-",
	RMul:    "*",
	RDiv:    "/",
	RTerm:   ";",
	RAssign: "=",

	// containers
	ROpenParen:  "(",
	RCloseParen: ")",
	ROpenBrace:  "{",
	RCloseBrace: "}",
	RQuote:      "\"",

	// separators
	RComma: ",",
	RSpace: " ",

	// built-in literals
	TrueLiteral:  "true",
	FalseLiteral: "false",

	// keyword
	Function: "routine",
}

var reservedTokensMap = map[string]TokenType{}

func init() {
	for i, c := range reservedTokens {
		reservedTokensMap[c] = i
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
		// fully ignore spaces
		case uni.IsSpace(curr):
		case isReserved(string(curr)):
			// checking for reserved single rune
			opt := lookup(string(curr))
			// semicolon, instakill statement
			if opt == RTerm {
				stmts = append(stmts, stmt)
				stmt = nil
				break sw_stmt
			}

			// if non-semicolon but still reserved, add to current statement
			// TODO impl parseOperator here to handle ops longer than 1 rune
			stmt = append(stmt, newSingleRuneToken(iter, opt))

			// if what we added was { or }, end the statement
			if opt == ROpenBrace || opt == RCloseBrace {
				stmts = append(stmts, stmt)
				stmt = nil
			}
		case uni.IsLetter(curr):
			// track start, end for token struct
			start := iter.Index()
			// parse entire token expecting a variable
			vr := variable(iter)
			end := iter.Index()

			// checking for reserved tokens longer than 1 rune
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
			// NOTE: unary minus is not caught here and is instead parsed as -1 * x later
			tok, err := numLiteral(iter)
			if err != nil {
				return nil, err
			}
			stmt = append(stmt, tok)
		case string(curr) == RQuote.String():
			_, ok = iter.Next() // consume the opening quote
			if !ok {
				return nil, ggErrs.Runtime("unterminated string literal")
			}
			tok, err := parseStringLiteral(iter)
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

func newSingleRuneToken(iter *iterator.Iter[rune], tokenType TokenType) Token {
	start := iter.Index()
	ret := Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(iter.Current()),
		TokenType: tokenType,
	}
	return ret
}

func parseStringLiteral(iter *iterator.Iter[rune]) (Token, error) {
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

	return Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(str),
		TokenType: StringLiteral,
	}, nil
}

// parses a number literal.
// currently, only number runes are supported (no decimal, scientific notation, etc.)
// iter points to the first rune of the number
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
		TokenType: IntLiteral,
	}, nil
}

// this checks runes with index in identifier > 0,
// the first rune is always a letter at this point
func idRune(r rune) bool {
	return uni.IsLetter(r) || uni.IsDigit(r) || r == '_'
}

// parse an entire identifier token
func variable(iter *iterator.Iter[rune]) Token {
	start := iter.Index()
	id := []rune{iter.Current()}

	next, ok := iter.Peek()
	for ok {
		if !idRune(next) {
			break
		}
		id = append(id, next)

		iter.Next() // consume the next rune we just checked
		next, ok = iter.Peek()
	}

	return Token{
		Start:     start,
		End:       iter.Index(),
		Str:       string(id),
		TokenType: Var,
	}
}
