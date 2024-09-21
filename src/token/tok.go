package token

import "fmt"

type Type int

const (
	beginOperators Type = iota
	Plus
	Minus
	Mul
	Div
	BitwiseAnd
	BitwiseOr
	LogicalNot
	LogicalAnd
	LogicalOr
	Equal
	NotEqual
	Assign
	endOperators

	beginContainers
	OpenParen
	CloseParen
	OpenBrace
	CloseBrace
	Quote
	endContainers

	beginSeparators
	Term
	Comma
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

func (t Type) IsOperator() bool {
	return t > beginOperators && t < endOperators
}
func (t Type) IsContainer() bool {
	return t > beginContainers && t < endContainers
}
func (t Type) IsSeparator() bool {
	return t > beginSeparators && t < endSeparators
}
func (t Type) IsIdentifier() bool {
	return t > beginIdentifiers && t < endIdentifiers
}
func (t Type) IsMathOperator() bool {
	return t == Plus || t == Minus || t == Mul || t == Div
}

func (t Type) String() string {
	if s, ok := reservedTokens[t]; ok {
		return s
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

var reservedTokens = map[Type]string{
	// operators
	Plus:       "+",
	Minus:      "-",
	Mul:        "*",
	Div:        "/",
	BitwiseAnd: "&",
	BitwiseOr:  "|",
	LogicalNot: "!",
	LogicalAnd: "&&",
	LogicalOr:  "||",
	Equal:      "==",
	NotEqual:   "!=",
	Assign:     "=",

	// terminators
	Term: ";",

	// containers
	OpenParen:  "(",
	CloseParen: ")",
	OpenBrace:  "{",
	CloseBrace: "}",
	Quote:      "\"",

	// separators
	Comma: ",",

	// built-in literals
	TrueLiteral:  "true",
	FalseLiteral: "false",

	// keyword
	Function: "routine",
}

var reservedTokensMap = map[string]Type{}

func init() {
	for i, c := range reservedTokens {
		reservedTokensMap[c] = i
	}
}

func isReserved(in string) bool {
	_, ok := reservedTokensMap[in]
	return ok
}

func isRuneReserved(r rune, typ Type) bool {
	result, ok := reservedTokensMap[string(r)]
	return ok && result == typ
}

func lookup(in string) Type {
	return reservedTokensMap[in]
}

type Token struct {
	Start     int `json:"-"`
	End       int `json:"-"`
	Str       string
	TokenType Type
}

func (t Token) String() string {
	return fmt.Sprintf("(%d-%d) %s", t.Start, t.End, t.Str)
}
