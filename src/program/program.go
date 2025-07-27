package program

import (
	"fmt"
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/operators"
	"gg-lang/src/stack"
	"gg-lang/src/variable"
	"strings"
)

type Program struct {
	scopes *stack.Stack[*Scope]
	OpMap  *operators.OpMap

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
	sb.WriteString(p.OpMap.String())
	return sb.String()
}

// New initializes the top Scope, declares every default builtin.Func,
// and registers every default operators.Operator
func New() *Program {
	scopes := stack.New[*Scope]()
	prog := &Program{
		scopes: scopes,
		OpMap:  operators.Default(),
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
	gg.Handle(err)
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

// essentially the same as RunExpression.
func (p *Program) runBlockStmt(block gg_ast.BlockStatement) error {
	for _, stmt := range block {
		err := p.RunExpression(stmt)
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
