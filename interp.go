package main

import (
	"fmt"
	"strconv"
)

type session struct {
	variables map[string]variable
	omap      *opmap
}

func (s *session) run(ast *ast) error {
	for i, expr := range ast.body {
		switch expr.kind() {
		case ExprAssignment:
			if err := s.evaluateAssignment(expr.(*assignmentExpression)); err != nil {
				return &runtimeError{fmt.Sprintf("expr %d: %v", i, err)}
			}
		}
	}

	return nil
}

func (s *session) evaluateAssignment(expr *assignmentExpression) error {
	switch expr.target.kind() {
	case ExprVariable:
	default:
		return &runtimeError{fmt.Sprintf("invalid assignment target: %s", expr.target.raw)}
	}

	newVar := variable{name: expr.target.raw}

	if expr.value.kind() > SentinelValueExpression {

		return fmt.Errorf("cannot make value for %s", exprString(expr))
	}

	val, typ, err := s.evaluateValueExpr(expr.value)
	if err != nil {
		return err
	}
	newVar.value = val
	newVar.typ = typ

	s.variables[newVar.name] = newVar

	return nil
}

func (s *session) evaluateValueExpr(expr valueExpression) (interface{}, varType, error) {
	switch expr.kind() {
	case ExprVariable:
		name := expr.(*identifier).name()
		if val, ok := s.variables[name]; ok {
			return val.value, val.typ, nil
		}
		return nil, 0, fmt.Errorf("undefined variable: %s", name)
	case ExprNumberLiteral:
		name := expr.(*identifier).name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return err, 0, err
		}
		return intVal, INTEGER, nil
	case ExprStringLiteral:
		return expr.(*identifier).name(), STRING, nil
	case ExprBinary:
		return s.evaluateBinaryExpr(expr.(*binaryExpression))
	default:
		return nil, 0, fmt.Errorf("invalid expression type: %v", expr)
	}
}

func (s *session) evaluateBinaryExpr(expr *binaryExpression) (interface{}, varType, error) {
	left, ltyp, err := s.evaluateValueExpr(expr.lhs)
	if err != nil {
		return nil, 0, err
	}

	right, rtyp, err := s.evaluateValueExpr(expr.rhs)
	if err != nil {
		return nil, 0, err
	}

	op, exists := s.omap.get(expr.op, ltyp, rtyp)
	if !exists {
		return nil, 0, fmt.Errorf("op %s not supported between types %v and %v", expr.op, ltyp, rtyp)
	}

	value := op.evaluate(left, right)
	resultType := op.resultType()

	return value, resultType, nil
}
