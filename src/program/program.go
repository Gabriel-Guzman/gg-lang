package program

import (
	"fmt"
	"gg-lang/src/builtin"
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/operators"
	"gg-lang/src/token"
	"gg-lang/src/variables"
	"os"
	"strings"
)

type Scope struct {
	Parent    *Scope
	variables map[string]*variables.Variable
}

// declares a variable in the current Scope and checks if it's already declared
func (s *Scope) declareVar(name string, value *variables.RuntimeValue) (*variables.Variable, error) {
	_, ok := s.variables[name]
	if ok {
		return nil, ggErrs.Runtime("variable '%s' already declared in this scope\n", name)
	}
	v := &variables.Variable{
		Name:         name,
		RuntimeValue: value,
	}
	s.variables[name] = v
	return v, nil
}

// declares a variable in the current Scope without checking if it's already declared
func (s *Scope) softDeclareVar(name string, value *variables.RuntimeValue) (*variables.Variable, error) {
	v := &variables.Variable{
		Name:         name,
		RuntimeValue: value,
	}
	s.variables[name] = v
	return v, nil
}

// declares a variable in the top Scope and checks if it's already declared
func (p *Program) declareVarTop(name string, value *variables.RuntimeValue) (*variables.Variable, error) {
	return p.top.declareVar(name, value)
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
	prog := &Program{
		top:     top,
		current: top,
		opMap:   operators.Default(),
	}

	for _, fn := range builtin.Defaults() {
		_, err := prog.declareVarTop(fn.Name(), &variables.RuntimeValue{
			Val: fn,
			Typ: variables.BuiltinFunction,
		})
		if err != nil {
			return nil
		}
	}

	return &Program{
		top:     top,
		current: top,
		opMap:   operators.Default(),
	}
}

func (p *Program) RunString(code string) error {
	stmts, err := token.TokenizeRunes([]rune(code))
	if err != nil {
		return err
	}

	ast, err := gg_ast.BuildFromStatements(stmts)
	ggErrs.Handle(err)
	if err != nil {
		return err
	}

	err = p.Run(ast)
	if err != nil {
		return err
	}
	return nil
}

func (p *Program) Run(ast *gg_ast.Ast) error {
	for _, expr := range ast.Body {
		err := p.RunExpression(expr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) RunExpression(expr gg_ast.Expression) error {
	switch expr.Kind() {
	case gg_ast.ExprAssignment:
		if err := p.evaluateAssignment(expr.(*gg_ast.AssignmentExpression)); err != nil {
			return err
		}

	case gg_ast.ExprFuncDecl:
		decl := expr.(*gg_ast.FunctionDeclExpression)
		_, err := p.declareVarTop(decl.Target.Raw, &variables.RuntimeValue{
			Val: decl,
			Typ: variables.Function,
		})
		if err != nil {
			return err
		}

	case gg_ast.ExprFunctionCall:
		call := expr.(*gg_ast.FunctionCallExpression)
		err := p.call(call)
		if err != nil {
			return err
		}
	default:
		return ggErrs.Crit("Invalid top-level expression: %s", expr.Kind().String())
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

func (p *Program) evaluateAssignment(expr *gg_ast.AssignmentExpression) error {
	if expr.Target.Kind() != gg_ast.ExprVariable {
		return ggErrs.Runtime("invalid assignment target: %s", expr.Target.Raw)
	}

	if expr.Value.Kind() > gg_ast.SentinelValueExpression {
		return ggErrs.Runtime("cannot make value for %v", expr)
	}

	val, err := p.evaluateValueExpr(expr.Value)
	if err != nil {
		return err
	}

	existing := p.findVariable(expr.Target.Name())
	if existing != nil {
		existing.RuntimeValue = val // garbage collect old value
		return nil
	}

	_, err = p.current.softDeclareVar(expr.Target.Raw, val)
	return err
}
