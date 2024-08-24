package main

type expressionKind int

const (
	ExprBinary expressionKind = iota
	ExprAssignment
	ExprNumberLiteral
	ExprIdentifier
)

type expression interface {
	kind() expressionKind
}

type valueExpression interface {
	name() string
	kind() expressionKind
}
