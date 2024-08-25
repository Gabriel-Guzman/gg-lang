package src

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/errors"
	"github.com/gabriel-guzman/gg-lang/src/godTree"
	"github.com/gabriel-guzman/gg-lang/src/operators"
	"github.com/gabriel-guzman/gg-lang/src/variables"
	"strconv"
)

type Session struct {
	Variables map[string]variables.Variable
	Omap      *operators.Opmap
}

func (s *Session) Run(ast *godTree.Ast) error {
	for i, expr := range ast.Body {
		switch expr.Kind() {
		case godTree.ExprAssignment:
			if err := s.evaluateAssignment(expr.(*godTree.AssignmentExpression)); err != nil {
				return &errors.RuntimeError{Message: fmt.Sprintf("expr %d: %v", i, err)}
			}
		}
	}

	return nil
}

func (s *Session) evaluateAssignment(expr *godTree.AssignmentExpression) error {
	switch expr.Target.Kind() {
	case godTree.ExprVariable:
	default:
		return &errors.RuntimeError{Message: fmt.Sprintf("invalid assignment target: %s", expr.Target.Raw)}
	}

	newVar := variables.Variable{Name: expr.Target.Raw}

	if expr.Value.Kind() > godTree.SentinelValueExpression {
		return &errors.RuntimeError{Message: fmt.Sprintf("cannot make value for %s", godTree.ExprString(expr, 0))}
	}

	val, typ, err := s.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}
	newVar.Value = val
	newVar.Typ = typ

	s.Variables[newVar.Name] = newVar

	return nil
}

func (s *Session) evaluateValueExpr(expr godTree.ValueExpression) (interface{}, variables.VarType, error) {
	switch expr.Kind() {
	case godTree.ExprVariable:
		name := expr.(*godTree.Identifier).Name()
		if val, ok := s.Variables[name]; ok {
			return val.Value, val.Typ, nil
		}
		return nil, 0, fmt.Errorf("undefined variable: %s", name)
	case godTree.ExprNumberLiteral:
		name := expr.(*godTree.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return err, 0, err
		}
		return intVal, variables.INTEGER, nil
	case godTree.ExprStringLiteral:
		return expr.(*godTree.Identifier).Name(), variables.STRING, nil
	case godTree.ExprBinary:
		return s.evaluateBinaryExpr(expr.(*godTree.BinaryExpression))
	default:
		return nil, 0, fmt.Errorf("invalid expression type: %v", expr)
	}
}

func (s *Session) evaluateBinaryExpr(expr *godTree.BinaryExpression) (interface{}, variables.VarType, error) {
	left, ltyp, err := s.evaluateValueExpr(expr.Lhs)
	if err != nil {
		return nil, 0, err
	}

	right, rtyp, err := s.evaluateValueExpr(expr.Rhs)
	if err != nil {
		return nil, 0, err
	}

	op, exists := s.Omap.Get(expr.Op, ltyp, rtyp)
	if !exists {
		return nil, 0, fmt.Errorf("op %s not supported between types %v and %v", expr.Op, ltyp, rtyp)
	}

	value := op.Evaluate(left, right)
	resultType := op.ResultType()

	return value, resultType, nil
}
