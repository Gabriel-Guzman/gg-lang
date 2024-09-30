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
	ExprIfElse
	ExprBlock
	ExprReturn
)

type Expression interface {
	Kind() ExpressionKind
}

type ValueExpression interface {
	Expression

	Name() string
}

// a BlockExpression is an Expression followed by a block statement
// such as token.For, token.Function
type BlockExpression interface {
	Expression

	SetStatements(statements []Expression)
}
