//go:generate stringer -type=IdExprKind
package godTree

import (
	"fmt"
	"gg-lang/src/tokenizer"
	"strings"
)

// utility to prevent wrong assignments
type IdExprKind int

const (
	IdExprNumber   = IdExprKind(ExprNumberLiteral)
	IdExprString   = IdExprKind(ExprStringLiteral)
	IdExprVariable = IdExprKind(ExprVariable)
)

type Identifier struct {
	Raw    string
	idKind IdExprKind
}

func (id *Identifier) MinShape() []tokenizer.TokenType {
	return []tokenizer.TokenType{tokenizer.IdentifierMask}
}

func (id *Identifier) Name() string {
	return id.Raw
}

func (id *Identifier) Kind() ExpressionKind { return ExpressionKind(id.idKind) }

// a + b
type BinaryExpression struct {
	Lhs  ValueExpression
	Op   string
	Rhs  ValueExpression
	Type ExpressionKind
}

func (be *BinaryExpression) MinShape() []tokenizer.TokenType {
	return []tokenizer.TokenType{
		tokenizer.IdentifierMask,
		tokenizer.OperatorMask,
		tokenizer.IdentifierMask,
	}
}

func (be *BinaryExpression) Name() string         { return be.Op }
func (be *BinaryExpression) Kind() ExpressionKind { return ExprBinary }

func newBinaryExpression(lhs ValueExpression, operator string, rhs ValueExpression) *BinaryExpression {
	return &BinaryExpression{
		Lhs: lhs,
		Op:  operator,
		Rhs: rhs,
	}
}

// a(b, c)
type FunctionCallExpression struct {
	Id     *Identifier
	Params []ValueExpression
}

func (fce *FunctionCallExpression) MinShape() []tokenizer.TokenType {
	return []tokenizer.TokenType{
		tokenizer.Var,
		tokenizer.ROpenParen,
		tokenizer.IdentifierMask,
	}
}

func (fce *FunctionCallExpression) Name() string         { return fce.Id.Name() }
func (fce *FunctionCallExpression) Kind() ExpressionKind { return ExprFunctionCall }

// a 32
type AssignmentExpression struct {
	Target Identifier
	Value  ValueExpression
}

func (ae *AssignmentExpression) MinShape() []tokenizer.TokenType {
	return []tokenizer.TokenType{
		tokenizer.Var,
		tokenizer.RAssign,
		tokenizer.IdentifierMask,
	}
}

func (ae *AssignmentExpression) Kind() ExpressionKind { return ExprAssignment }
func newAssignmentExpression(target *Identifier, value ValueExpression) *AssignmentExpression {
	return &AssignmentExpression{
		Target: *target,
		Value:  value,
	}
}

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
		ExprString(&val.Target, d+1, sb)
	case ExprBinary:
		val := e.(*BinaryExpression)
		w("operation of ")
		ExprString(val.Lhs, d+1, sb)
		sb.WriteString("\n")
		w(val.Op)
		ExprString(val.Rhs, d+1, sb)
	case ExprNumberLiteral:
		goto IdentifierStr
	case ExprVariable:
		goto IdentifierStr
	case ExprStringLiteral:
		goto IdentifierStr
	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}
	return

IdentifierStr:
	id := e.(*Identifier)
	w("Ident " + id.idKind.String() + " " + id.Raw)
}
