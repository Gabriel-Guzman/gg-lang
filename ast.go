package main

import (
	"errors"
	"fmt"
	"strings"
)

type ast struct {
	body []expression
}

func (a *ast) String() string {
	var sb strings.Builder
	for _, expr := range a.body {
		sb.WriteString(exprString(expr))
		sb.WriteString("\n")
	}
	return sb.String()
}

func newAST(toks [][]token) (*ast, error) {
	a := &ast{}
	for _, stmt := range toks {
		tokIter := newIter(stmt)
		expr, err := parseStmt(tokIter)
		if err != nil {
			return nil, err
		}

		a.body = append(a.body, expr)
	}

	return a, nil
}

func parseStmt(tokIter *iter[token]) (expression, error) {
	curr, exists := tokIter.Next()
	if !exists {
		return nil, errors.New("expected a statement")
	}
	switch curr.tokenType {
	case VAR:
		// var decl
		_, exists := tokIter.Next()
		if !exists {
			return nil, fmt.Errorf("solo expressions are not allowed: %s", curr.str)
		}

		expr, err := parseValueExpr(tokIter)
		if err != nil {
			return nil, err
		}

		id, err := newIdentifier(curr)
		if err != nil {
			return nil, err
		}
		at := newAssignmentExpression(*id, expr)
		return at, nil
	}

	return nil, errors.New("invalid statement")
}

func parseValueExpr(tokIter *iter[token]) (valueExpression, error) {
	curr, exists := tokIter.Current()
	if !exists {
		return nil, fmt.Errorf("expected a value expression")
	}

	switch curr.tokenType {
	case OPERATOR:
		return nil, fmt.Errorf("unexpected op %s at %d", curr.str, curr.start)
	}

	// if this value expression is only one token
	_, exists = tokIter.Peek()
	if !exists {
		return newIdentifier(curr)
	}

	// if its more than one token, it should be a binary expression
	lhs, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	expectOp, _ := tokIter.Next()
	if expectOp.tokenType != OPERATOR {
		return nil, fmt.Errorf("expected an op after name expression %s, got %s at %d", curr.str, expectOp.str, expectOp.start)
	}

	tokIter.Next()
	rhs, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	return newBinaryExpression(lhs, expectOp.str, rhs), nil
}

func newIdentifier(t token) (*identifier, error) {
	var ik idKind
	switch t.tokenType {
	case NUMBER_LITERAL:
		ik = IdExprNumber
	case VAR:
		ik = IdVariable
	case STRING_LITERAL:
		ik = IdExprString
	default:
		return nil, fmt.Errorf("invalid value expression %s at %d: ", t.str, t.start)
	}
	return &identifier{raw: t.str, idKind: ik}, nil
}
