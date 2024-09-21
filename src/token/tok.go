package token

import "fmt"

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
	endSeparators

	beginIdentifiers
	Ident
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

func isStr(str string, typ TokenType) bool {
	result, ok := reservedTokensMap[str]
	return ok && result == typ
}

func isRuneReserved(r rune, typ TokenType) bool {
	result, ok := reservedTokensMap[string(r)]
	return ok && result == typ
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
