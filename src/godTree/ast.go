package godTree

import (
	"errors"
	"fmt"
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
		sb.WriteString(ExprString(expr, 0))
		sb.WriteString("\n")
	}
	return sb.String()
}

func NewAST() (*Ast, error) {
	a := &Ast{}
	return a, nil
}

func (a *Ast) ExecStmts(tokens [][]tokenizer.Token) error {
	for _, stmt := range tokens {
		tokIter := iterator.NewIter(stmt)
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
	case tokenizer.VAR:
		// var decl
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

	return nil, errors.New("invalid statement")
}

func parseValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	curr, exists := tokIter.Current()
	if !exists {
		return nil, fmt.Errorf("expected a value expression")
	}

	switch curr.TokenType {
	case tokenizer.OPERATOR:
		return nil, fmt.Errorf("unexpected op %s at %d", curr.Str, curr.Start)
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
	if expectOp.TokenType != tokenizer.OPERATOR {
		return nil, fmt.Errorf("expected an op after name expression %s, got %s at %d", curr.Str, expectOp.Str, expectOp.Start)
	}

	tokIter.Next()
	rhs, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	return newBinaryExpression(lhs, expectOp.Str, rhs), nil
}

func newIdentifier(t tokenizer.Token) (*Identifier, error) {
	var ik idKind
	switch t.TokenType {
	case tokenizer.NUMBER_LITERAL:
		ik = IdExprNumber
	case tokenizer.VAR:
		ik = IdVariable
	case tokenizer.STRING_LITERAL:
		ik = IdExprString
	default:
		return nil, fmt.Errorf("invalid value expression %s at %d: ", t.Str, t.Start)
	}
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}
