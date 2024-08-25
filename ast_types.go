package main

type expressionKind int

const (
	/*
		valueExpression implementing types
	*/
	ExprBinary expressionKind = iota
	ExprNumberLiteral
	ExprVariable
	ExprStringLiteral
	SentinelValueExpression

	/*
	   assignmentExpression implementing types
	*/
	ExprAssignment
)

type expression interface {
	kind() expressionKind
}

type valueExpression interface {
	name() string
	kind() expressionKind
}
