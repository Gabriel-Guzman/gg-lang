package main

import (
	"fmt"
	"strconv"
)

type session struct {
	variables map[string]variable
}

func (s *session) run(ast *ast) error {
	for i, expr := range ast.body {
		switch expr.kind() {
		case ExprAssignment:
			if err := s.evaluateAssignment(expr.(*assignmentExpression)); err != nil {
				return fmt.Errorf("expr %d: %v", i, err)
			}
		}
	}

	return nil
}

func (s *session) evaluateAssignment(expr *assignmentExpression) error {
	switch expr.target.kind() {
	case ExprIdentifier:
	default:
		return fmt.Errorf("invalid assignment target: %s", expr.target.raw)
	}

	newVar := variable{name: expr.target.raw}

	switch expr.value.kind() {
	case ExprAssignment:
		return fmt.Errorf("one assignment per statement")
	}

	val, typ, err := s.toValueExpectExpr(expr.value)
	if err != nil {
		return err
	}
	newVar.value = val
	newVar.typ = typ

	s.variables[newVar.name] = newVar

	return nil
}

func (s *session) evaluateBinary(expr *binaryExpression) (int, varType, error) {
	left, ltyp, err := s.toValueExpectIdent(expr.lhs.(*identifier))
	if err != nil {
		return 0, 0, err
	}

	right, rtyp, err := s.toValueExpectExpr(expr.rhs)
	if err != nil {
		return 0, 0, err
	}

	if ltyp != INTEGER || rtyp != INTEGER {
		return 0, 0, fmt.Errorf("invalid types for binary operation: %v %v", left, right)
	}

	if len(expr.operator.str) != 1 {
		return 0, 0, fmt.Errorf("invalid operator length: %s", expr.operator.str)
	}

	asRune := []rune(expr.operator.str)[0]

	switch asRune {
	case R_PLUS:
		return left.(int) + right.(int), INTEGER, nil
	case R_MINUS:
		return left.(int) - right.(int), INTEGER, nil
	case R_MUL:
		return left.(int) * right.(int), INTEGER, nil
	case R_DIV:
		return left.(int) / right.(int), INTEGER, nil
	default:
		return 0, 0, fmt.Errorf("invalid operator: %s", string(asRune))
	}
}

func (s *session) toValueExpectExpr(expr expression) (interface{}, varType, error) {
	switch expr.kind() {
	case ExprIdentifier:
		return s.toValueExpectIdent(expr.(*identifier))
	case ExprNumberLiteral:
		return s.toValueExpectIdent(expr.(*numberLiteral))
	case ExprBinary:
		return s.evaluateBinary(expr.(*binaryExpression))
	default:
		return nil, 0, fmt.Errorf("invalid expression type: %v", expr)
	}
}

func (s *session) toValueExpectIdent(ident valueExpression) (interface{}, varType, error) {
	name := ident.name()

	switch ident.kind() {
	case ExprNumberLiteral:
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return err, 0, err
		}

		return intVal, INTEGER, nil
	case ExprIdentifier:
		if val, ok := s.variables[name]; ok {
			return val.value, val.typ, nil
		} else {
			return nil, 0, fmt.Errorf("undefined variable: %s", name)
		}
	default:
		return nil, 0, fmt.Errorf("invalid value type: %v", ident)
	}
}
