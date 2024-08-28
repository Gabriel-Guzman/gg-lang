package godTree

import "github.com/gabriel-guzman/gg-lang/src/tokenizer"

type ExpressionKind int

const (
	/*
		valueExpression implementing types
	*/
	ExprBinary ExpressionKind = iota
	ExprNumberLiteral
	ExprVariable
	ExprStringLiteral
	ExprFunctionCall
	SentinelValueExpression

	/*
	   assignmentExpression implementing types
	*/
	ExprAssignment
)

type Expression interface {
	Kind() ExpressionKind
	MinShape() []tokenizer.TokenType
}

type ValueExpression interface {
	Name() string
	// Expression members
	Kind() ExpressionKind
	MinShape() []tokenizer.TokenType
}
