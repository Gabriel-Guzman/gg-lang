//go:generate stringer -type=IdExprKind
package gg_ast

import (
	"fmt"
	"strings"
)

// utility to prevent wrong assignments
type IdExprKind int

const (
	IdExprNumber   = IdExprKind(ExprIntLiteral)
	IdExprString   = IdExprKind(ExprStringLiteral)
	IdExprBool     = IdExprKind(ExprBoolLiteral)
	IdExprVariable = IdExprKind(ExprVariable)
)

// a
type Identifier struct {
	Raw    string
	idKind IdExprKind
}

func (id *Identifier) Name() string {
	return id.Raw
}
func (id *Identifier) Kind() ExpressionKind { return ExpressionKind(id.idKind) }

// { print(a) }
type BlockStatement []Expression

func (bs BlockStatement) SetStatements(s []Expression) { copy(bs, s) }
func (bs BlockStatement) Kind() ExpressionKind         { return ExprBlock }

// a + b
type BinaryExpression struct {
	Lhs ValueExpression
	Op  string
	Rhs ValueExpression
}

func (be *BinaryExpression) Name() string         { return be.Op }
func (be *BinaryExpression) Kind() ExpressionKind { return ExprBinary }

// a(b, c)
type FunctionCallExpression struct {
	Id   *Identifier
	Args []ValueExpression
}

func (fce *FunctionCallExpression) Name() string         { return fce.Id.Name() }
func (fce *FunctionCallExpression) Kind() ExpressionKind { return ExprFunctionCall }

// a = 32
type AssignmentExpression struct {
	Target *Identifier
	Value  ValueExpression
}

func (ae *AssignmentExpression) Kind() ExpressionKind { return ExprAssignment }

// routine a(b, c) {
type FunctionDeclExpression struct {
	Target *Identifier
	Params []string
	Body   BlockStatement
}

func (fde *FunctionDeclExpression) Kind() ExpressionKind { return ExprFuncDecl }
func (fde *FunctionDeclExpression) SetStatements(s []Expression) {
	fde.Body = s
}
func (fde *FunctionDeclExpression) Name() string {
	return fde.Target.Name()
}

// if a == b { } else if a == c { } else { }
type IfElseStatement struct {
	Condition      ValueExpression
	Body           BlockStatement
	ElseExpression Expression // optional
}

func (ife *IfElseStatement) Kind() ExpressionKind { return ExprIfElse }
func (ife *IfElseStatement) SetStatements(s []Expression) {
	ife.Body = s
}

// for i != 10 {
type ForLoopExpression struct {
	Condition ValueExpression
	Body      BlockStatement
}

func (fle *ForLoopExpression) Kind() ExpressionKind { return ExprForLoop }
func (fle *ForLoopExpression) SetStatements(s []Expression) {
	fle.Body = s
}

type ReturnStatement struct {
	Value ValueExpression
}

func (rs *ReturnStatement) Kind() ExpressionKind { return ExprReturn }

func ind(count int) string {
	var spaces []rune
	for i := 0; i < count*4; i++ {
		spaces = append(spaces, ' ')
	}

	return string(spaces)
}

func NoBuilderExprString(e Expression) string {
	sb := &strings.Builder{}
	ExprString(e, 0, sb)
	return sb.String()
}

func ExprString(e Expression, d int, sb *strings.Builder) {
	w := func(s string) {
		sb.WriteString(ind(d) + s)
	}
	sb.WriteString("\n")

	switch e.(type) {
	case *AssignmentExpression:
		val := e.(*AssignmentExpression)
		w("assign of ")
		ExprString(val.Value, d+1, sb)
		sb.WriteString("\n")
		w(" to")
		ExprString(val.Target, d+1, sb)
	case *BinaryExpression:
		val := e.(*BinaryExpression)
		w("operation of ")
		ExprString(val.Lhs, d+1, sb)
		sb.WriteString("\n")
		w(val.Op)
		ExprString(val.Rhs, d+1, sb)
	case *Identifier:
		id := e.(*Identifier)
		w("Ident " + id.idKind.String() + " " + id.Raw)
	case *FunctionCallExpression:
		val := e.(*FunctionCallExpression)
		w("call to " + val.Id.Name())
		for _, param := range val.Args {
			ExprString(param, d+1, sb)
		}
	case *FunctionDeclExpression:
		val := e.(*FunctionDeclExpression)
		w("decl of " + val.Target.Name() + fmt.Sprintf("(%s)\n", strings.Join(val.Params, ", ")))
		w(" to do")
		for _, expr := range val.Body {
			ExprString(expr, d+1, sb)
		}
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}
	return
}
