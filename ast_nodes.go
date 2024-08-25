package main

import "fmt"

// utility to prevent wrong assignments
type idKind int

const (
	IdExprNumber = idKind(ExprNumberLiteral)
	IdExprString = idKind(ExprStringLiteral)
	IdVariable   = idKind(ExprVariable)
)

type identifier struct {
	raw    string
	idKind idKind
}

func (id *identifier) name() string {
	return id.raw
}

func (id *identifier) kind() expressionKind { return expressionKind(id.idKind) }

// a + b
type binaryExpression struct {
	lhs valueExpression
	op  string
	rhs valueExpression
}

func (be *binaryExpression) name() string         { return be.op }
func (be *binaryExpression) kind() expressionKind { return ExprBinary }
func newBinaryExpression(lhs valueExpression, operator string, rhs valueExpression) *binaryExpression {
	return &binaryExpression{
		lhs: lhs,
		op:  operator,
		rhs: rhs,
	}
}

// a 32
type assignmentExpression struct {
	target identifier
	value  valueExpression
}

func (ae *assignmentExpression) kind() expressionKind { return ExprAssignment }
func newAssignmentExpression(target identifier, value valueExpression) *assignmentExpression {
	return &assignmentExpression{
		target: target,
		value:  value,
	}
}

func exprString(e expression) string {
	switch e.kind() {
	case ExprAssignment:
		val := e.(*assignmentExpression)
		return fmt.Sprintf(`assign of (%v) to (%v)`, exprString(val.value), val.target)
	case ExprBinary:
		val := e.(*binaryExpression)
		return fmt.Sprintf(`operate (%v) (%v) (%v)`, exprString(val.lhs), val.op, exprString(val.rhs))
	case ExprNumberLiteral:
		goto IdentifierStr
	case ExprVariable:
		goto IdentifierStr
	case ExprStringLiteral:
		goto IdentifierStr
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}

IdentifierStr:
	return e.(*identifier).raw
}
