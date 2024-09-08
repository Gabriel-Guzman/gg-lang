package program

import (
	"fmt"
	"gg-lang/src/builtin"
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/operators"
	"gg-lang/src/variables"
	"os"
	"strings"
)

type Scope struct {
	Parent    *Scope
	variables map[string]*variables.Variable
}

type Program struct {
	top     *Scope
	current *Scope
	opMap   *operators.OpMap
}

func (p *Program) String() string {
	var sb strings.Builder
	sb.WriteString("Variables:\n")
	for k, v := range p.top.variables {
		sb.WriteString(fmt.Sprintf("\t%s: %+v\n", k, v))
	}
	sb.WriteString("\nOperators:\n")
	sb.WriteString(p.opMap.String())
	return sb.String()
}

func New() *Program {
	top := &Scope{variables: make(map[string]*variables.Variable)}

	for _, fn := range builtin.Defaults() {
		//top.variables[fn.Name()]] = newV
		top.variables[fn.Name()] = &variables.Variable{
			Name: fn.Name(),
			Value: &variables.RuntimeValue{
				Val: fn,
				Typ: variables.BuiltinFunction,
			},
		}
	}

	return &Program{
		top:     top,
		current: top,
		opMap:   operators.Default(),
	}
}

func (p *Program) Run(ast *godTree.Ast) error {
	for i, expr := range ast.Body {
		switch expr.Kind() {
		case godTree.ExprAssignment:
			if err := p.evaluateAssignment(expr.(*godTree.AssignmentExpression)); err != nil {
				return ggErrs.Runtime("expr %d: %v", i, err)
			}

		case godTree.ExprFuncDecl:
			decl := expr.(*godTree.FunctionDeclExpression)
			newVar := variables.Variable{
				Name: decl.Target.Raw,
				Value: &variables.RuntimeValue{
					Val: decl,
					Typ: variables.Function,
				},
			}
			p.top.variables[newVar.Name] = &newVar

		case godTree.ExprFunctionCall:
			fcall := expr.(*godTree.FunctionCallExpression)
			err := p.funcCall(fcall)
			if err != nil {
				return ggErrs.Runtime("expr %d: %v", i, err)
			}
		default:
			return ggErrs.Crit("Invalid top-level expression: %s", expr.Kind().String())
		}
	}

	return nil
}

func (p *Program) enterNewScope() {
	ns := &Scope{
		Parent:    p.current,
		variables: make(map[string]*variables.Variable),
	}

	p.current = ns
}

func (p *Program) exitScope() {
	if p.current == p.top {
		fmt.Println("exitScope called on top scope. Goodbye!")
		os.Exit(1)
	}
	p.current = p.current.Parent
}

func (p *Program) findVariable(name string) *variables.Variable {
	s := p.current
	for {
		res, ok := s.variables[name]
		if ok {
			return res
		}

		if s.Parent == nil {
			return nil
		}

		s = s.Parent
	}
}

func (p *Program) evaluateAssignment(expr *godTree.AssignmentExpression) error {
	switch expr.Target.Kind() {
	case godTree.ExprVariable:
	default:
		return ggErrs.Runtime("invalid assignment target: %s", expr.Target.Raw)
	}

	if expr.Value.Kind() > godTree.SentinelValueExpression {
		return ggErrs.Runtime("cannot make value for %v", expr)
	}

	val, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}

	existing := p.findVariable(expr.Target.Name())
	if existing != nil {
		existing.Value = val
		//existing.Typ = val.Typ
		return nil
	}

	newVar := variables.Variable{Name: expr.Target.Raw}
	newVar.Value = val
	//newVar.Typ = val.Typ

	p.current.variables[newVar.Name] = &newVar
	return nil
}
