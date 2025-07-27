//go:generate stringer -type=ExpressionKind

package gg_ast

type ExpressionKind int

const (
	/*
		ValueExpression implementing kinds
	*/
	ExprBinary ExpressionKind = iota
	ExprUnary
	ExprIntLiteral
	ExprBoolLiteral
	ExprVariable
	ExprStringLiteral
	ExprFunctionCall
	ExprObject
	ExprArrayDecl
	ExprArrayIndex
	ExprArrayIndexAssignment
	ExprDotAccess
	ExprParenthesized
	SentinelValueExpression

	/*
	   Expression implementing kinds
	*/
	ExprAssignment
	ExprDotAccessAssignment
	ExprFuncDecl
	ExprForLoop
	ExprIfElse
	ExprBlock
	ExprReturn
	ExprTryCatch
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
