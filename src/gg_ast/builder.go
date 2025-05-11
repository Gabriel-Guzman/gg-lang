package gg_ast

import (
	"gg-lang/src/gg"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

type builder struct {
	par *parser.Parser[token.Token]
}

func parseBlockStatement(par *parser.Parser[token.Token]) (BlockStatement, error) {
	if !advanceIfCurrIs(par, token.OpenBrace) {
		return nil, gg.Syntax("expected opening brace for block statement\n%s", par.String())
	}
	var expressions []Expression
	for par.HasCurr {
		if advanceIfCurrIs(par, token.CloseBrace) {
			return expressions, nil
		}

		stmt, err := parseExpression(par)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, stmt)
	}

	return nil, gg.Syntax("no closing brace for block statement\n%s", par.String())
}

func newAstBuilder(ins []token.Token) *builder {
	par := parser.New(ins)

	return &builder{
		par: par,
	}
}

func BuildFromString(ins string) (*Ast, error) {
	tokens, err := token.TokenizeRunes([]rune(ins))
	if err != nil {
		return nil, err
	}

	return BuildFromTokens(tokens)
}

func BuildFromTokens(ins []token.Token) (*Ast, error) {
	a := newAstBuilder(ins)

	var expressions []Expression
	for a.par.HasCurr {
		expr, err := parseExpression(a.par)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return &Ast{Body: expressions}, nil
}
