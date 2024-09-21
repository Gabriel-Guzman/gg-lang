package gg_ast

import (
	"gg-lang/src/token"
	"strings"
)

type Ast struct {
	Body []Expression
}

func (a *Ast) String() string {
	var sb strings.Builder
	for _, expr := range a.Body {
		ExprString(expr, 0, &sb)
	}
	return sb.String()
}

func tokStringer(t token.Token) string {
	return t.Str
}
