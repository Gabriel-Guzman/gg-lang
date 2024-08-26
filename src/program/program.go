package program

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/ggErrs"
	"github.com/gabriel-guzman/gg-lang/src/godTree"
	"github.com/gabriel-guzman/gg-lang/src/operators"
	"github.com/gabriel-guzman/gg-lang/src/variables"
	"strconv"
	"strings"
)

type Program struct {
	variables map[string]variables.Variable
	opMap     *operators.OpMap
}

func (p *Program) String() string {
	var sb strings.Builder
	sb.WriteString("variables:\n")
	for k, v := range p.variables {
		sb.WriteString(fmt.Sprintf("%s: %+v\n", k, v))
	}
	sb.WriteString("Operators:\n")
	sb.WriteString(p.opMap.String())
	return sb.String()
}

func New() *Program {
	return &Program{
		variables: make(map[string]variables.Variable),
		opMap:     operators.Default(),
	}
}

func (p *Program) Run(ast *godTree.Ast) error {
	for i, expr := range ast.Body {
		switch expr.Kind() {
		case godTree.ExprAssignment:
			if err := p.evaluateAssignment(expr.(*godTree.AssignmentExpression)); err != nil {
				return ggErrs.NewRuntime("expr %d: %v", i, err)
			}
		}
	}

	return nil
}

func (p *Program) evaluateAssignment(expr *godTree.AssignmentExpression) error {
	switch expr.Target.Kind() {
	case godTree.ExprVariable:
	default:
		return ggErrs.NewRuntime(fmt.Sprintf("invalid assignment target: %s", expr.Target.Raw))
	}

	newVar := variables.Variable{Name: expr.Target.Raw}

	if expr.Value.Kind() > godTree.SentinelValueExpression {
		return ggErrs.NewRuntime("cannot make value for %v", expr)
	}

	val, typ, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}

	newVar.Value = val
	newVar.Typ = typ

	p.variables[newVar.Name] = newVar

	return nil
}

func (p *Program) evaluateValueExpr(expr godTree.ValueExpression) (interface{}, variables.VarType, error) {
	switch expr.Kind() {
	case godTree.ExprVariable:
		name := expr.(*godTree.Identifier).Name()
		if val, ok := p.variables[name]; ok {
			return val.Value, val.Typ, nil
		}
		return nil, 0, ggErrs.NewRuntime("undefined variable: %s", name)
	case godTree.ExprNumberLiteral:
		name := expr.(*godTree.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, 0, ggErrs.NewInternal(err, "unable to evaluate number literal: %s", name)
		}
		return intVal, variables.Integer, nil
	case godTree.ExprStringLiteral:
		return expr.(*godTree.Identifier).Name(), variables.String, nil
	case godTree.ExprBinary:
		binExp := expr.(*godTree.BinaryExpression)
		left, ltyp, err := p.evaluateValueExpr(binExp.Lhs)
		if err != nil {
			return nil, 0, err
		}

		right, rtyp, err := p.evaluateValueExpr(binExp.Rhs)
		if err != nil {
			return nil, 0, err
		}

		op, exists := p.opMap.Get(binExp.Op, ltyp, rtyp)
		if !exists {
			return nil, 0, ggErrs.NewRuntime("op %s not supported between types %v and %v", binExp.Op, ltyp, rtyp)
		}

		value := op.Evaluate(left, right)
		resultType := op.ResultType()

		return value, resultType, nil
	default:
		return nil, 0, ggErrs.NewInternal(nil, "invalid expression type: %v", expr)
	}
}
