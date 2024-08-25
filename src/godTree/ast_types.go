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
	SentinelValueExpression

	/*
	   assignmentExpression implementing types
	*/
	ExprAssignment
)

type Expression interface {
	Kind() ExpressionKind
}

type ValueExpression interface {
	Name() string
	Kind() ExpressionKind
}
