package godTree

import (
	"errors"
	"fmt"
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
	switch curr.TokenType {
	// var decl
	case tokenizer.VAR:
		_, exists := tokIter.Next()
		if !exists {
			return nil, fmt.Errorf("solo expressions are not allowed: %s", curr.Str)
		}

		expr, err := parseValueExpr(tokIter)
		if err != nil {
			return nil, err
		}

		id, err := newIdentifier(curr)
		if err != nil {
			return nil, err
		}

		at := newAssignmentExpression(*id, expr)
		return at, nil
	}

	return nil, ggErrs.NewRuntime("invalid statement")
}

func parseValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	curr, exists := tokIter.Current()
	if !exists {
		return nil, ggErrs.NewRuntime("expected a value expression")
	}

	switch curr.TokenType {
	case tokenizer.Operator:
		return nil, ggErrs.NewRuntime("unexpected op %s at %d", curr.Str, curr.Start)
	}

	// if this value expression is only one token
	_, exists = tokIter.Peek()
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
		return nil, ggErrs.NewRuntime("expected an op after name expression %s, got %s at %d", curr.Str, expectOp.Str, expectOp.Start)
	}

	_, exists = tokIter.Next()
	if !exists {
		return nil, ggErrs.NewRuntime("expected a value expression after op")
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
	case tokenizer.VAR:
		ik = IdExprVariable
	case tokenizer.StringLiteral:
		ik = IdExprString
	default:
		return nil, ggErrs.NewRuntime("invalid identifier %s at %d: ", t.Str, t.Start)
	}
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}
