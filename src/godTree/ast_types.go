//go:generate stringer -type=ExpressionKind

package godTree

type ExpressionKind int

const (
	/*
		valueExpression implementing types
	*/
	ExprBinary ExpressionKind = iota
	ExprNumberLiteral
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

type IValExpr interface {
	Expression

	Name() string
}
