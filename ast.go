package main

import (
	"errors"
	"fmt"
	"strings"
)

type ast struct {
	body []expression
}

func exprString(e expression) string {
	switch e.kind() {
	case ExprAssignment:
		val := e.(*assignmentExpression)

		return fmt.Sprintf(`assign of (%v) to (%v)`, exprString(val.value), val.target)
	case ExprBinary:
		val := e.(*binaryExpression)
		return fmt.Sprintf(`operate (%v) (%v) (%v)`, exprString(val.lhs), val.operator.str, exprString(val.rhs))

	case ExprNumberLiteral:
		return e.(*numberLiteral).raw
	case ExprIdentifier:
		return e.(*identifier).raw
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}
}

func (a *ast) String() string {
	var sb strings.Builder
	for _, expr := range a.body {
		sb.WriteString(exprString(expr))
		sb.WriteString("\n")
	}
	return sb.String()
}

func fromTokens(toks [][]word) (*ast, error) {
	a := &ast{}
	for _, stmt := range toks {
		witer := newWordIter(stmt)
		expr, err := parseStmt(witer)
		if err != nil {
			return nil, err
		}

		a.body = append(a.body, expr)
	}

	return a, nil
}

func parseStmt(witer *wordIter) (expression, error) {
	curr, exists := witer.Next()
	if !exists {
		return nil, errors.New("expected a statement")
	}
	switch curr.role {
	case VAR:
		// var decl
		_, exists := witer.Next()
		if !exists {
			return nil, fmt.Errorf("solo expressions are not allowed: %s", curr.str)
		}

		expr, err := parseValueExpr(witer)
		if err != nil {
			return nil, err
		}

		id := newIdentifier(curr)
		at := newAssignmentExpression(*id, expr)
		return at, nil
	}

	return nil, errors.New("invalid statement")
}

func parseValueExpr(witer *wordIter) (expression, error) {
	curr, exists := witer.Current()
	if !exists {
		return nil, fmt.Errorf("expected a name expression")
	}

	_, exists = witer.Peek()
	if exists {
		expectValue := curr
		lhs, err := valueExprFromWord(expectValue)
		if err != nil {
			return nil, err
		}
		// if operator, binary expression
		expectOp, _ := witer.Next()
		if expectOp.role != OPERATOR {
			return nil, fmt.Errorf("expected an operator after a name expression, got %s", expectOp.str)
		}

		witer.Next()
		rhs, err := parseValueExpr(witer)
		if err != nil {
			return nil, err
		}

		return newBinaryExpression(lhs, expectOp, rhs), nil
	}

	return valueExprFromWord(curr)
}

func valueExprFromWord(w word) (valueExpression, error) {
	switch w.role {
	case NUMBER_LITERAL:
		return &numberLiteral{raw: w.str}, nil
	case VAR:
		return &identifier{raw: w.str}, nil
	default:
		return nil, fmt.Errorf("invalid name expression: %s", w.str)
	}
}
