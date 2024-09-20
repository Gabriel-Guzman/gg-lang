package tokenizer

import (
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/parser"
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
	par := parser.New(ins)
	par.SetStringer(func(in rune) string {
		return string(in)
	})
	par.SetSeparator("")

	var stmts [][]Token

	var currStmt []Token

	endStmt := func() {
		stmts = append(stmts, currStmt)
		currStmt = nil
	}

	addToCurr := func(tok Token) {
		currStmt = append(currStmt, tok)
	}

	isRTok := func(t TokenType) bool {
		return isReserved(string(par.Curr)) && lookup(string(par.Curr)) == t
	}

	// par starts at -1
	par.Advance()
	for par.HasCurr {
		switch {
		case shouldIgnore(par.Curr):
			par.Advance()
		case isRTok(RTerm): // begin reserved characters
			par.Advance()
			endStmt()
		case isRTok(RQuote):
			strTok, err := parseStringLiteral(par)
			if err != nil {
				return nil, err
			}
			addToCurr(strTok)
		case isReserved(string(par.Curr)) && lookup(string(par.Curr)).IsOperator():
			tok, err := _parseOperator(par)
			if err != nil {
				return nil, err
			}
			addToCurr(tok)
		case isRTok(ROpenBrace): // begin containers
			fallthrough
		case isRTok(RCloseBrace):
			addToCurr(singleRuneTok(par, lookup(string(par.Curr))))
			endStmt()
		case isRTok(RComma):
			fallthrough
		case isRTok(ROpenParen):
			fallthrough
		case isRTok(RCloseParen):
			addToCurr(singleRuneTok(par, lookup(string(par.Curr))))
		case uni.IsDigit(par.Curr):
			numTok, err := parseNumLiteral(par)
			if err != nil {
				return nil, err
			}
			addToCurr(numTok)
		case uni.IsLetter(par.Curr):
			idTok, err := parseIdentifier(par)
			if err != nil {
				return nil, err
			}
			addToCurr(idTok)
		}
	}

	return stmts, nil
}

func singleRuneTok(p *parser.Parser[rune], tokType TokenType) Token {
	curr := p.Curr
	ret := Token{
		Start:     p.Index(),
		End:       p.Index() + 1,
		Str:       string(curr),
		TokenType: tokType,
	}

	p.Advance()
	return ret
}

// parsers must consume every rune that they add to a token
func parseIdentifier(p *parser.Parser[rune]) (Token, error) {
	start := p.Index()
	if !p.HasCurr {
		return Token{}, ggErrs.Crit("identifier parser called with nothing in parser\n%s", p.String())
	}

	if !uni.IsLetter(p.Curr) {
		return Token{}, ggErrs.Crit("expected letter, got %s\n%s", string(p.Curr), p.String())
	}

	id := ""
	for {
		if p.HasCurr && idRune(p.Curr) {
			id += string(p.Curr)
			p.Advance()
			continue
		}
		// not a letter or digit or underscore
		break
	}

	if id == "" {
		return Token{}, ggErrs.Crit("could not parse identifier\n%s", p.String())
	}

	if isReserved(id) {
		tt := lookup(id)
		return Token{
			Start:     start,
			End:       p.Index() + 1,
			Str:       id,
			TokenType: tt,
		}, nil
	}

	return Token{
		Start:     start,
		End:       p.Index() + 1,
		Str:       id,
		TokenType: Var,
	}, nil
}

func parseNumLiteral(p *parser.Parser[rune]) (Token, error) {
	start := p.Index()
	if !p.HasCurr {
		return Token{}, ggErrs.Crit("number parser called with nothing in parser\n%s", p.String())
	}

	num := ""

	for {
		if p.HasCurr && (uni.IsDigit(p.Curr)) {
			num += string(p.Curr)
			p.Advance()
			continue
		}
		// not a digit
		break
	}

	if num == "" {
		return Token{}, ggErrs.Crit("could not parse number\n%s", p.String())
	}

	return Token{
		Start:     start,
		End:       p.Index() + 1,
		Str:       num,
		TokenType: IntLiteral,
	}, nil
}

func _parseOperator(p *parser.Parser[rune]) (Token, error) {
	start := p.Index()
	if !p.HasCurr {
		return Token{}, ggErrs.Crit("operator parser called with nothing in parser\n%s", p.String())
	}
	op := ""

	for {
		if p.HasCurr {
			if isReserved(string(p.Curr)) && lookup(string(p.Curr)).IsOperator() {
				// token is an operator, add it to the operator string
				op += string(p.Curr)
				p.Advance()
				continue
			}
			// not reserved or not an operator, done parsing
			break
		} else {
			// nothing left in parser
			break
		}
	}

	if op == "" {
		return Token{}, ggErrs.Crit("could not parse operator\n%s", p.String())
	}

	if isReserved(op) {
		realOp := lookup(op)
		if !realOp.IsOperator() {
			return Token{}, ggErrs.Runtime("unknown operator \n%s", p.String())
		}

		return Token{
			Start:     start,
			End:       p.Index() + 1,
			Str:       op,
			TokenType: realOp,
		}, nil
	}

	return Token{}, ggErrs.Runtime("unknown operator \n%s", p.String())
}

func parseStringLiteral(p *parser.Parser[rune]) (Token, error) {
	if !p.HasCurr {
		return Token{}, ggErrs.Crit("string literal parser called with nothing in parser\n%s", p.String())
	}
	if string(p.Curr) != reservedTokens[RQuote] {
		return Token{}, ggErrs.Crit("string literal parser called on non-quote\n%s", p.String())
	}

	p.Advance() // consume opening quote
	start := p.Index()

	str := ""

	for {
		if !p.HasCurr {
			return Token{}, ggErrs.Runtime("unterminated string literal\n%s", p.String())
		}
		if string(p.Curr) == reservedTokens[RQuote] {
			p.Advance() // consume closing quote
			break
		}

		str += string(p.Curr)
		p.Advance()
	}

	return Token{
		Start:     start,
		End:       p.Index() + 1,
		Str:       str,
		TokenType: StringLiteral,
	}, nil
}

func shouldIgnore(curr rune) bool {
	return uni.IsSpace(curr)
}

// this checks runes with index in identifier > 0,
// the first rune is always a letter at this point
func idRune(r rune) bool {
	return uni.IsLetter(r) || uni.IsDigit(r) || r == '_'
}
