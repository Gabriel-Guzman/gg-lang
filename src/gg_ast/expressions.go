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

type Identifier struct {
	Raw    string
	idKind IdExprKind
}

func (id *Identifier) Name() string {
	return id.Raw
}
func (id *Identifier) Kind() ExpressionKind { return ExpressionKind(id.idKind) }

// a + b
type BinaryExpression struct {
	Lhs ValExpression
	Op  string
	Rhs ValExpression
}

func (be *BinaryExpression) Name() string         { return be.Op }
func (be *BinaryExpression) Kind() ExpressionKind { return ExprBinary }

// a(b, c)
type FunctionCallExpression struct {
	Id   *Identifier
	Args []ValExpression
}

func (fce *FunctionCallExpression) Name() string         { return fce.Id.Name() }
func (fce *FunctionCallExpression) Kind() ExpressionKind { return ExprFunctionCall }

// a = 32
type AssignmentExpression struct {
	Target *Identifier
	Value  ValExpression
}

func (ae *AssignmentExpression) Kind() ExpressionKind { return ExprAssignment }

// routine a(b, c) {
type FunctionDeclExpression struct {
	Target *Identifier
	Params []string
	Value  []Expression
}

func (fde *FunctionDeclExpression) Kind() ExpressionKind { return ExprFuncDecl }

func ind(count int) string {
	var spaces []rune
	for i := 0; i < count*4; i++ {
		spaces = append(spaces, ' ')
	}

	return string(spaces)
}

func ExprString(e Expression, d int, sb *strings.Builder) {
	w := func(s string) {
		sb.WriteString(ind(d) + s)
	}
	sb.WriteString("\n")
	switch e.Kind() {
	case ExprAssignment:
		val := e.(*AssignmentExpression)
		w("assign of ")
		ExprString(val.Value, d+1, sb)
		sb.WriteString("\n")
		w(" to")
		ExprString(val.Target, d+1, sb)
	case ExprBinary:
		val := e.(*BinaryExpression)
		w("operation of ")
		ExprString(val.Lhs, d+1, sb)
		sb.WriteString("\n")
		w(val.Op)
		ExprString(val.Rhs, d+1, sb)
	case ExprIntLiteral:
		fallthrough
	case ExprVariable:
		fallthrough
	case ExprStringLiteral:
		fallthrough
	case ExprBoolLiteral:
		id := e.(*Identifier)
		w("Ident " + id.idKind.String() + " " + id.Raw)
	case ExprFunctionCall:
		val := e.(*FunctionCallExpression)
		w("call to " + val.Id.Name())
		for _, param := range val.Args {
			ExprString(param, d+1, sb)
		}
	case ExprFuncDecl:
		val := e.(*FunctionDeclExpression)
		w("decl of " + val.Target.Name() + fmt.Sprintf("(%s)\n", strings.Join(val.Params, ", ")))
		w(" to do")
		for _, expr := range val.Value {
			ExprString(expr, d+1, sb)
		}
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}
	return
}
