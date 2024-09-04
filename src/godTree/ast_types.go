//go:generate stringer -type=ExpressionKind

package godTree

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
	ExprFuncDecl
)

type Expression interface {
	Kind() ExpressionKind
}

type ValueExpression interface {
	Expression

	Name() string
}
