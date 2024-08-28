package godTree

import (
	"errors"
	"github.com/gabriel-guzman/gg-lang/src/ggErrs"
	"github.com/gabriel-guzman/gg-lang/src/iterator"
	"github.com/gabriel-guzman/gg-lang/src/tokenizer"
	"strings"
)

type Ast struct {
	Body []Expression
}

func (a *Ast) String() string {
	var sb strings.Builder
	for _, expr := range a.Body {
		ExprString(expr, 0, &sb)
	}
	return sb.String()
}

func New() *Ast {
	a := &Ast{}
	return a
}

func (a *Ast) ParseStmts(tokens [][]tokenizer.Token) error {
	for _, stmt := range tokens {
		tokIter := iterator.New(stmt)
		expr, err := parseStmt(tokIter)
		if err != nil {
			return err
		}

		a.Body = append(a.Body, expr)
	}
	return nil
}

func parseStmt(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	curr, exists := tokIter.Next()
	if !exists {
		return nil, errors.New("expected a statement")
	}

	next, ok := tokIter.Peek()
	nextTokType := next.TokenType

	exprs := []Expression{
		&Identifier{},
		&BinaryExpression{},
		&FunctionCallExpression{},
		&AssignmentExpression{},
	}
	switch {
	case matchExpr(tokIter, &Identifier{}):
	case matchExpr(tokIter, &BinaryExpression{}):
	case matchExpr(tokIter, &FunctionCallExpression{}):
	case matchExpr(tokIter, &AssignmentExpression{}):
	}
	for _, expr := range exprs {
		if matchExpr(tokIter, expr) {
			switch expr.(type) {
			case *Identifier:
				return nil, ggErrs.Runtime("didnt expect an identifier: %s", tokIter.String())
			}
		}

	}

	switch curr.TokenType {
	case tokenizer.Var:
		if ok && nextTokType == tokenizer.ROpenParen { // checking nextTokType isnt necessary due to lookup, but its faster
			// function call

		}

		if ok && nextTokType.IsOperator() {

		}

	}

	return nil, ggErrs.Runtime("invalid statement")
}

func matchExpr(tokIter *iterator.Iter[tokenizer.Token], expression Expression) bool {
	iter := tokIter.Copy()
	ok := true
	for _, exprMask := range expression.MinShape() {
		if !ok {
			return false
		}
		tok := iter.Current()
		if exprMask&tok.TokenType == 0 {
			return false
		}

		_, ok = iter.Next()
	}
	return true
}

//func parseFuncCallExpr(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
//
//}

func parseAssignmentExpr(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	curr := tokIter.Current()
	id, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	expr, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	at := newAssignmentExpression(*id, expr)
	return at, nil
}

func parseValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	curr := tokIter.Current()

	switch curr.TokenType {
	case tokenizer.Operator:
		return nil, ggErrs.Runtime("unexpected op %s at %d", curr.Str, curr.Start)
	}

	// if this value expression is only one token
	_, exists := tokIter.Peek()
	if !exists {
		return newIdentifier(curr)
	}

	// if its more than one token, it should be a binary expression
	lhs, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	expectOp, _ := tokIter.Next()
	if expectOp.TokenType != tokenizer.Operator {
		return nil, ggErrs.Runtime("expected an op after name expression %s, got %s at %d", curr.Str, expectOp.Str, expectOp.Start)
	}

	_, exists = tokIter.Next()
	if !exists {
		return nil, ggErrs.Runtime("expected a value expression after op")
	}
	rhs, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	return newBinaryExpression(lhs, expectOp.Str, rhs), nil
}

func newIdentifier(t tokenizer.Token) (*Identifier, error) {
	var ik IdExprKind
	switch t.TokenType {
	case tokenizer.NumberLiteral:
		ik = IdExprNumber
	case tokenizer.Var:
		ik = IdExprVariable
	case tokenizer.StringLiteral:
		ik = IdExprString
	default:
		return nil, ggErrs.Runtime("invalid identifier %s at %d: ", t.Str, t.Start)
	}
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}
