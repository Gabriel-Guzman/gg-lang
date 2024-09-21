//go:generate stringer -type=ExpressionKind

package gg_ast

type ExpressionKind int

const (
	/*
		ValueExpression implementing kinds
	*/
	ExprBinary ExpressionKind = iota
	ExprIntLiteral
	ExprBoolLiteral
	ExprVariable
	ExprStringLiteral
	ExprFunctionCall
	SentinelValueExpression

	/*
	   Expression implementing kinds
	*/
	ExprAssignment
	ExprFuncDecl
	ExprForLoop
)

type Expression interface {
	Kind() ExpressionKind
}

type ValueExpression interface {
	Expression

	Name() string
}
