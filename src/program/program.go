package program

import (
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/operators"
	"gg-lang/src/stack"
	"gg-lang/src/variable"
	"os"
	"strings"
)

type Scope struct {
	Parent    *Scope
	variables map[string]*variable.Variable

	caller *RuntimeFunc
}

// declares a variables.Variable in the Program.current Scope and checks if it's already declared
func (s *Scope) declareVar(name string, value *variable.RuntimeValue) (*variable.Variable, error) {
	_, ok := s.variables[name]
	if ok {
		return nil, ggErrs.Runtime("variable '%s' already declared in this scope\n%v", name, s)
	}
	v := &variable.Variable{
		Name:         name,
		RuntimeValue: value,
	}
	s.variables[name] = v
	return v, nil
}

// declares a variables.Variable in the Program.current Scope without checking if it's already declared
func (s *Scope) softDeclareVar(name string, value *variable.RuntimeValue) (*variable.Variable, error) {
	v := &variable.Variable{
		Name:         name,
		RuntimeValue: value,
	}
	s.variables[name] = v
	return v, nil
}

type Program struct {
	scopes *stack.Stack[*Scope]
	opMap  *operators.OpMap

	returnValue *variable.RuntimeValue
}

func (p *Program) currentScope() *Scope {
	curr, _ := p.scopes.Peek()
	return curr
}

func (p *Program) String() string {
	var sb strings.Builder
	sb.WriteString("Variables:\n")
	for k, v := range p.currentScope().variables {
		sb.WriteString(fmt.Sprintf("\t%s: %+v\n", k, v))
	}
	sb.WriteString("\nOperators:\n")
	sb.WriteString(p.opMap.String())
	return sb.String()
}

// New initializes the top Scope, declares every default builtin.Func,
// and registers every default operators.Operator
func New() *Program {
	scopes := stack.New[*Scope]()
	prog := &Program{
		scopes: scopes,
		opMap:  operators.Default(),
	}
	prog.enterNewScope()

	for _, fn := range Defaults() {
		_, err := prog.currentScope().declareVar(fn.Name(), &variable.RuntimeValue{
			Val: fn,
			Typ: variable.BuiltinFunction,
		})
		if err != nil {
			return nil
		}
	}

	return prog
}

// a shortcut for executing a string of code
func (p *Program) RunString(code string) error {
	ast, err := gg_ast.BuildFromString(code)
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
		err := p.RunStmt(expr)
		if err != nil {
			return err
		}
	}

	return nil
}

// essentially the same as RunStmt.
func (p *Program) runBlockStmt(block gg_ast.BlockStatement) error {
	for _, stmt := range block {
		err := p.RunStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// a shortcut runBlockStmt in a new scope.
func (p *Program) runBlockStmtNewScope(block gg_ast.BlockStatement) error {
	p.enterNewScope()
	defer p.exitScope()
	return p.runBlockStmt(block)
}

func (p *Program) RunStmt(expr gg_ast.Expression) error {
	// dont execute anything if there's a return value right now
	if p.returnValue != nil {
		return nil
	}
	switch expr.(type) {
	case *gg_ast.ReturnStatement:
		val, err := p.evaluateValueExpr(expr.(*gg_ast.ReturnStatement).Value)
		if err != nil {
			return err
		}
		p.returnValue = val
	case gg_ast.BlockStatement:
		block := expr.(gg_ast.BlockStatement)
		err := p.runBlockStmtNewScope(block)
		if err != nil {
			return err
		}
	case *gg_ast.AssignmentExpression:
		if err := p.evaluateAssignment(expr.(*gg_ast.AssignmentExpression)); err != nil {
			return err
		}
	case *gg_ast.FunctionDeclExpression:
		decl := expr.(*gg_ast.FunctionDeclExpression)
		_, err := p.currentScope().declareVar(decl.Target.Raw, &variable.RuntimeValue{
			Val: RuntimeFuncFromDecl(decl, p.currentScope()),
			Typ: variable.Function,
		})
		if err != nil {
			return err
		}
	case *gg_ast.ForLoopExpression:
		loop := expr.(*gg_ast.ForLoopExpression)
		for {
			val, err := p.evaluateValueExpr(loop.Condition)
			if err != nil {
				return err
			}
			if _, ok := val.Val.(bool); !ok {
				return ggErrs.Runtime("loop condition must evaluate to bool\n%+v", expr)
			}
			if !val.Val.(bool) {
				break
			}
			err = p.runBlockStmtNewScope(loop.Body)
			if err != nil {
				return err
			}
		}
	case *gg_ast.IfElseStatement:
		ifElse := expr.(*gg_ast.IfElseStatement)
		err := p.execIfElse(ifElse)
		if err != nil {
			return err
		}
	case *gg_ast.FunctionCallExpression:
		call := expr.(*gg_ast.FunctionCallExpression)
		_, err := p.call(call)
		if err != nil {
			return err
		}
	default:
		return ggErrs.Crit("Invalid top-level expression: %s", expr.Kind().String())
	}
	return nil
}

func (p *Program) execIfElse(expr gg_ast.Expression) error {
	ifElse := expr.(*gg_ast.IfElseStatement)
	cond, err := p.evaluateValueExpr(ifElse.Condition)
	if err != nil {
		return err
	}
	if _, ok := cond.Val.(bool); !ok {
		return ggErrs.Runtime("if condition must evaluate to bool\n%+v", expr)
	}
	if cond.Val.(bool) {
		err = p.runBlockStmtNewScope(ifElse.Body)
		if err != nil {
			return err
		}
	} else {
		if ifElse.ElseExpression != nil {
			err = p.execIfElse(ifElse.ElseExpression)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Program) enterNewScope() {
	ns := &Scope{
		// NOTE: the current scope is the new scopes parent here.
		// this will not always be the case. when a function decl captures a scope,
		// it may have pushed a scope with a parent other than the current scope.
		Parent:    p.currentScope(),
		variables: make(map[string]*variable.Variable),
	}
	p.scopes.Push(ns)
}

func (p *Program) enterCapturedScope(scope *Scope) {
	p.scopes.Push(scope)
}

func (p *Program) exitScope() {
	p.scopes.Pop()
	if p.currentScope() == nil {
		fmt.Println("exitScope called on top scope. Goodbye!")
		os.Exit(1)
	}
}

func (s *Scope) findVariable(name string) *variable.Variable {
	c := s
	for {
		res, ok := c.variables[name]
		if ok {
			return res
		}

		if c.Parent == nil {
			return nil
		}

		c = c.Parent
	}
}

// a shortcut for querying a variable in the 'current' scope
func (p *Program) findVariable(name string) *variable.Variable {
	return p.currentScope().findVariable(name)
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

	_, err = p.currentScope().softDeclareVar(expr.Target.Raw, val)
	return err
}
