package program

import (
	"fmt"
	"gg-lang/src/gg"
	"gg-lang/src/variable"
	"os"
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
		return nil, gg.Runtime("variable '%s' already declared in this scope\n%v", name, s)
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
