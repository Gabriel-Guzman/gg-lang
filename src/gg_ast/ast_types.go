//go:generate stringer -type=ExpressionKind

package gg_ast

type ExpressionKind int

const (
	/*
		valueExpression implementing types
	*/
	ExprBinary ExpressionKind = iota
	ExprIntLiteral
	ExprBoolLiteral
	ExprVariable
	ExprStringLiteral
	ExprFunctionCall
	SentinelValueExpression

	/*
	   assignmentExpression implementing types
	*/
	ExprAssignment
	ExprFuncDecl
)

type Expression interface {
	Kind() ExpressionKind
}

type ValExpression interface {
	Expression

	Name() string
}
