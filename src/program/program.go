package program

import (
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/operators"
	"gg-lang/src/variables"
	"strconv"
	"strings"
)

type Program struct {
	variables map[string]variables.Variable
	opMap     *operators.OpMap
}

func (p *Program) String() string {
	var sb strings.Builder
	sb.WriteString("Variables:\n")
	for k, v := range p.variables {
		sb.WriteString(fmt.Sprintf("%s: %+v\n", k, v))
	}
	sb.WriteString("\nOperators:\n")
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
				return ggErrs.Runtime("expr %d: %v", i, err)
			}
		default:
			return ggErrs.Crit("invalid expression kind: %d", expr.Kind())
		}
	}

	return nil
}

func (p *Program) evaluateAssignment(expr *godTree.AssignmentExpression) error {
	switch expr.Target.Kind() {
	case godTree.ExprVariable:
	default:
		return ggErrs.Runtime("invalid assignment target: %s", expr.Target.Raw)
	}

	newVar := variables.Variable{Name: expr.Target.Raw}

	if expr.Value.Kind() > godTree.SentinelValueExpression {
		return ggErrs.Runtime("cannot make value for %v", expr)
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
		return nil, 0, ggErrs.Runtime("undefined variable: %s", name)
	case godTree.ExprNumberLiteral:
		name := expr.(*godTree.Identifier).Name()
		intVal, err := strconv.Atoi(name)
		if err != nil {
			return nil, 0, ggErrs.Crit("unable to evaluate number literal: %s", err.Error())
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
			return nil, 0, ggErrs.Runtime("evaluateValueExpr: op %s not supported between types %v and %v", binExp.Op, ltyp, rtyp)
		}

		value := op.Evaluate(left, right)
		resultType := op.ResultType()

		return value, resultType, nil
	default:
		return nil, 0, ggErrs.Crit("evaluateValueExpr: invalid expression type: %v", expr)
	}
}
