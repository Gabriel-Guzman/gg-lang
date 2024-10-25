package token

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/parser"
	uni "unicode"
)

func tokenize(par *parser.Parser[rune]) ([]Token, error) {
	tk := &tkzr{
		Par: par,
	}
	var toks []Token
	a := func(tok Token) {
		toks = append(toks, tok)
	}
	for tk.Par.HasCurr {
		switch {
		case shouldIgnore(par.Curr):
			tk.Par.Advance()
		case isRuneReserved(tk.Par.Curr, Term):
			a(tk.parseReservedSingleRuneTok(Term))
		case isRuneReserved(tk.Par.Curr, Quote):
			strTok, err := parseStringLiteral(par)
			if err != nil {
				return nil, err
			}
			a(strTok)
		case isReserved(string(par.Curr)) && lookup(string(par.Curr)).IsOperator():
			tok, err := parseOperator(par)
			if err != nil {
				return nil, err
			}
			a(tok)
		case isRuneReserved(tk.Par.Curr, OpenBrace):
			a(tk.parseReservedSingleRuneTok(OpenBrace))
		case isRuneReserved(tk.Par.Curr, CloseBrace):
			a(tk.parseReservedSingleRuneTok(CloseBrace))
		case isRuneReserved(tk.Par.Curr, Comma):
			a(tk.parseReservedSingleRuneTok(Comma))
		case isRuneReserved(tk.Par.Curr, OpenParen):
			a(tk.parseReservedSingleRuneTok(OpenParen))
		case isRuneReserved(tk.Par.Curr, CloseParen):
			a(tk.parseReservedSingleRuneTok(CloseParen))
		case uni.IsDigit(par.Curr):
			numTok, err := parseNumLiteral(par)
			if err != nil {
				return nil, err
			}
			a(numTok)
		case uni.IsLetter(par.Curr):
			idTok, err := parseIdentifier(par)
			if err != nil {
				return nil, err
			}
			a(idTok)
		default:
			return nil, ggErrs.Crit("unexpected character\n%s", par.String())
		}
	}

	if len(toks) == 0 {
		return nil, nil
	}
	return toks, nil
}

type tkzr struct {
	Par *parser.Parser[rune]
}

func TokenizeRunes(ins []rune) ([]Token, error) {
	par := parser.New(ins)
	par.SetStringer(func(in rune) string {
		return string(in)
	})
	par.SetSeparator("")

	return tokenize(par)
}

func (t *tkzr) parseReservedSingleRuneTok(tokType Type) Token {
	curr := t.Par.Curr
	ret := Token{
		Start:     t.Par.Index(),
		End:       t.Par.Index() + 1,
		Str:       string(curr),
		TokenType: tokType,
	}

	t.Par.Advance()
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
		TokenType: Ident,
	}, nil
}

func parseNumLiteral(p *parser.Parser[rune]) (Token, error) {
	start := p.Index()
	if !p.HasCurr {
		return Token{}, ggErrs.Crit("number parser called with nothing in parser\n%s", p.String())
	}

	num := ""

	for {
		if p.HasCurr && uni.IsDigit(p.Curr) {
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

func parseOperator(p *parser.Parser[rune]) (Token, error) {
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
	if string(p.Curr) != reservedTokens[Quote] {
		return Token{}, ggErrs.Crit("string literal parser called on non-quote\n%s", p.String())
	}

	p.Advance() // consume opening quote
	start := p.Index()

	str := ""

	for {
		if !p.HasCurr {
			return Token{}, ggErrs.Runtime("unterminated string literal\n%s", p.String())
		}
		if string(p.Curr) == reservedTokens[Quote] {
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
