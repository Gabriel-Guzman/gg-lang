//go:generate stringer -type=IdExprKind
package gg_ast

import (
	"fmt"
	"gg-lang/src/token"
	"strings"
)

// utility to prevent wrong assignments
type IdExprKind int

const (
	IdExprNumber    = IdExprKind(ExprIntLiteral)
	IdExprString    = IdExprKind(ExprStringLiteral)
	IdExprBool      = IdExprKind(ExprBoolLiteral)
	IdExprVariable  = IdExprKind(ExprVariable)
	IdExprDotAccess = IdExprKind(ExprDotAccess)
)

type Literal struct {
	Tok token.Token
}

func (l *Literal) Name() string {
	return l.Tok.Symbol
}

func getLiteralKind(tok token.Token) ExpressionKind {
	switch tok.TokenType {
	case token.FalseLiteral:
		fallthrough
	case token.TrueLiteral:
		return ExprBoolLiteral
	case token.IntLiteral:
		return ExprIntLiteral
	case token.StringLiteral:
		return ExprStringLiteral
	default:
		panic(fmt.Sprintf("unknown token type: %v", tok.TokenType))
	}
}
func (l *Literal) Kind() ExpressionKind {
	return getLiteralKind(l.Tok)
}

// a
type Identifier struct {
	Tok    token.Token
	idKind IdExprKind
}

func (id *Identifier) Name() string {
	return id.Tok.Symbol
}
func (id *Identifier) Kind() ExpressionKind {
	switch id.idKind {
	case IdExprNumber:
		return ExprIntLiteral
	case IdExprString:
		return ExprStringLiteral
	case IdExprBool:
		return ExprBoolLiteral
	case IdExprVariable:
		return ExprVariable
	case IdExprDotAccess:
		return ExprDotAccess
	}

	panic(fmt.Sprintf("unknown identifier kind: %v", id.idKind))
}

// { print(a) }
type BlockStatement []Expression

func (bs BlockStatement) SetStatements(s []Expression) { copy(bs, s) }
func (bs BlockStatement) Kind() ExpressionKind         { return ExprBlock }

// -b
type UnaryExpression struct {
	Op  token.Token
	Rhs ValueExpression
}

func (be *UnaryExpression) Name() string         { return be.Op.Symbol }
func (be *UnaryExpression) Kind() ExpressionKind { return ExprUnary }

// (a)
type ParenthesizedExpression struct {
	Expr ValueExpression
}

func (pe *ParenthesizedExpression) Name() string         { return fmt.Sprintf("(%s)", pe.Expr.Name()) }
func (pe *ParenthesizedExpression) Kind() ExpressionKind { return ExprParenthesized }

// a + b
type BinaryExpression struct {
	Lhs ValueExpression
	Op  token.Token
	Rhs ValueExpression
}

func (be *BinaryExpression) Name() string         { return be.Op.Symbol }
func (be *BinaryExpression) Kind() ExpressionKind { return ExprBinary }

// a(b, c)
type FunctionCallExpression struct {
	Id   *Identifier
	Args []ValueExpression
}

func (fce *FunctionCallExpression) Name() string         { return fce.Id.Name() }
func (fce *FunctionCallExpression) Kind() ExpressionKind { return ExprFunctionCall }

// try { a = 32 } catch (e) { print(e) }
type TryCatchExpression struct {
	Try     *BlockStatement
	Catch   *CatchExpression
	Finally *BlockStatement
}

func (t TryCatchExpression) Kind() ExpressionKind {
	return ExprTryCatch
}

// catch (e) { print(e) }
type CatchExpression struct {
	ErrorParam string
	Body       *BlockStatement
}

// a = 32
type AssignmentExpression struct {
	Target *Identifier
	Value  ValueExpression
}

func (ae *AssignmentExpression) Kind() ExpressionKind { return ExprAssignment }

// a.b = 5
type DotAccessAssignmentExpression struct {
	Target *DotAccessExpression
	Value  ValueExpression
}

func (d *DotAccessAssignmentExpression) Kind() ExpressionKind { return ExprDotAccessAssignment }

// routine a(b, c) {
type FunctionDeclExpression struct {
	Target *Identifier
	Params []token.Token
	Body   BlockStatement
}

func (fde *FunctionDeclExpression) Kind() ExpressionKind { return ExprFuncDecl }
func (fde *FunctionDeclExpression) SetStatements(s []Expression) {
	fde.Body = s
}
func (fde *FunctionDeclExpression) Name() string {
	return fde.Target.Name()
}

// [1, 2, 3]
type ArrayDeclExpression struct {
	Elements []ValueExpression
}

func (a ArrayDeclExpression) Kind() ExpressionKind { return ExprArrayDecl }
func (a ArrayDeclExpression) Name() string {
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Name())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// a[1]
type ArrayIndexExpression struct {
	Array *Identifier
	Index ValueExpression
}

// a[1] = 1
type ArrayIndexAssignmentExpression struct {
	*ArrayIndexExpression
	Value ValueExpression
}

func (ai ArrayIndexAssignmentExpression) Kind() ExpressionKind { return ExprArrayIndexAssignment }
func (ai ArrayIndexAssignmentExpression) Name() string {
	return fmt.Sprintf("%s[%s] = %s", ai.Array.Name(), ai.Index.Name(), ai.Value.Name())
}

func (ai ArrayIndexExpression) Kind() ExpressionKind { return ExprArrayIndex }
func (ai ArrayIndexExpression) Name() string {
	return fmt.Sprintf("%s[%s]", ai.Array.Name(), ai.Index.Name())
}

// { x: 1, y: 2, z: 3 }
type ObjectExpression struct {
	Properties map[string]ValueExpression
}

func (o ObjectExpression) Kind() ExpressionKind {
	return ExprObject
}

func (o ObjectExpression) Name() string {
	return "[object]"
}

type DotAccessExpression struct {
	AccessChain []string
}

func (d DotAccessExpression) Kind() ExpressionKind {
	return ExprDotAccess
}

func (d DotAccessExpression) Name() string {
	return strings.Join(d.AccessChain, ".")
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

	switch val := e.(type) {
	case *AssignmentExpression:
		w("assign of ")
		ExprString(val.Value, d+1, sb)
		sb.WriteString("\n")
		w(" to")
		ExprString(val.Target, d+1, sb)
	case *BinaryExpression:
		w("operation of ")
		ExprString(val.Lhs, d+1, sb)
		sb.WriteString("\n")
		w(val.Op.Symbol)
		ExprString(val.Rhs, d+1, sb)
	case *Identifier:
		w("Ident " + val.idKind.String() + " " + val.Tok.Symbol + "\n")
	case *FunctionCallExpression:
		w("call to " + val.Id.Name())
		for _, param := range val.Args {
			ExprString(param, d+1, sb)
		}
	case *FunctionDeclExpression:
		w("decl of " + val.Target.Name())
		w(" to do")
		for _, expr := range val.Body {
			ExprString(expr, d+1, sb)
		}
	case *UnaryExpression:
		w("operation of " + val.Op.Symbol)
		ExprString(val.Rhs, d+1, sb)
	case BlockStatement:
		w("Block {")
		for _, expr := range val {
			ExprString(expr, d+1, sb)
		}

		w("} end block")
	case *DotAccessExpression:
		w("access to " + val.Name())
	case *DotAccessAssignmentExpression:
		w("assign of ")
		ExprString(val.Value, d+1, sb)
		sb.WriteString("\n")
		w("to")
		ExprString(val.Target, d+1, sb)
		w(" (dot access)")
		w("\n")
	case *ParenthesizedExpression:
		w("parenthesized expression of ")
		ExprString(val.Expr, d+1, sb)
	case *ArrayDeclExpression:
		w("array declaration of ")
		for i, expr := range val.Elements {
			ExprString(expr, d+1, sb)
			if i < len(val.Elements)-1 {
				w(", ")
			}
		}
	case *ArrayIndexExpression:
		w("access to array index of ")
		ExprString(val.Array, d+1, sb)
		w(" [")
		ExprString(val.Index, d+1, sb)
		w("]")
		w("\n")
	case *ArrayIndexAssignmentExpression:
		w("assign of ")
		ExprString(val.Value, d+1, sb)
		sb.WriteString("\n")
		w("to")
		ExprString(val.ArrayIndexExpression, d+1, sb)
		w(" (array index assignment)")
		w("\n")

	default:
		panic(fmt.Sprintf("unknown expression type: %T", e))
	}
	return
}
