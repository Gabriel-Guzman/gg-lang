package godTree

import (
	"errors"
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/tokenizer"
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
		if tokIter.HasNext() {
			return ggErrs.Runtime("couldnt finish parsing statement\n%s", tokIter.String())
		}

		a.Body = append(a.Body, expr)
	}
	return nil
}

func parseStmt(tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	// move to first token in stmt
	curr, exists := tokIter.Next()
	if !exists {
		return nil, errors.New("expected a statement")
	}
	if !curr.TokenType.IsIdentifier() {
		return nil, ggErrs.Runtime("expected an identifier\n%s", tokIter.String())
	}

	if !(curr.TokenType == tokenizer.Var) {
		return nil, ggErrs.Runtime("expected a variable identifier\n%s", tokIter.String())
	}

	id, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	next, ok := tokIter.Peek()
	nextTokType := next.TokenType

	if ok && nextTokType == tokenizer.RAssign {
		tokIter.Next() // consume the '=' token
		expr, err := parseAssignmentExpr(id, tokIter)
		return expr, err
	}

	// not assign, its value. this could be a lone expression i.e. "swag";
	tokIter.Reset()
	expr, err := parseValueExpr(tokIter)
	return expr, err
}

func parseAssignmentExpr(id *Identifier, tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	expr, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	at := newAssignmentExpression(id, expr)
	return at, nil
}

func parseValueExpr(tokIter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	// first word in value
	curr, exists := tokIter.Next()

	if curr.TokenType.IsOperator() {
		return nil, ggErrs.Runtime("unexpected op %s at %d", curr.Str, curr.Start)
	}

	firstId, err := newIdentifier(curr)
	if err != nil {
		return nil, err
	}

	// try for next word
	token2, exists := tokIter.Next()
	if !exists {
		// no next word, return identifier
		return firstId, nil
	}

	var firstExpr ValueExpression
	firstExpr = firstId
	// try for function call
	if token2.TokenType == tokenizer.ROpenParen {
		funcName := firstId
		firstExpr, err = newFuncCallExpression(funcName, tokIter)
		if err != nil {
			return nil, err
		}
		// tokIter is currently at the closing parenthesis ')'
		// if there's something next, it could be an operator
		_, ok := tokIter.Next()
		if !ok {
			return firstExpr, nil
		}
	}

	// operator, try for binary expression
	if token2.TokenType.IsOperator() {
		if token2.TokenType == tokenizer.RAssign {
			return nil, ggErrs.Runtime("invalid = in value expression: %s", tokIter.String())
		}

		lhs := firstExpr
		op := token2.Str
		rhs, err := parseValueExpr(tokIter)
		if err != nil {
			return nil, err
		}

		return newBinaryExpression(lhs, op, rhs), nil

	}
	//if token2.TokenType.IsSeparator() {
	return firstExpr, nil
	//}

	// this is now a value expression followed by a value expression
	//return nil, ggErrs.Runtime("invalid expression\n%s", tokIter.String())
}

func newFuncCallExpression(funcName *Identifier, iter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	//params, err := paramsList(iter)
	//if err != nil {
	//	return nil, err
	//}
	//var args []ValueExpression
	//for _, param := range params {
	//	ve, err := parseValueExpr(iterator.New(param))
	//	if err != nil {
	//		return nil, err
	//	}
	//	args = append(args, ve)
	//}
	// iter is current at the opening parenthesis '('.

	// if the next char is a closing paren, dont check for value expr

	nextTok, ok := iter.Peek()
	if !ok {
		return nil, ggErrs.Runtime("expected closing parenthesis ')' or args after function name\n%s", iter.String())
	}
	if nextTok.TokenType == tokenizer.RCloseParen {
		iter.Next() // consume the closing parenthesis ')'
		return &FunctionCallExpression{
			Id:     funcName,
			Params: nil,
		}, nil
	}

	var args []ValueExpression
	for {
		val, err := parseValueExpr(iter)
		if err != nil {
			return nil, err
		}
		if !iter.HasCurrent() {
			return nil, ggErrs.Runtime("unexpected end of arg list\n%s", iter.String())
		}
		args = append(args, val)
		mbComma := iter.Current()
		if mbComma.TokenType == tokenizer.RComma {
			continue
		}
		if mbComma.TokenType == tokenizer.RCloseParen {
			break
		}
	}
	return &FunctionCallExpression{
		Id:     funcName,
		Params: args,
	}, nil
}
func paramsList(tok *iterator.Iter[tokenizer.Token]) ([][]tokenizer.Token, error) {
	var list [][]tokenizer.Token

	var currList []tokenizer.Token
	for {
		curr, ok := tok.Next()
		if !ok {
			return nil, ggErrs.Runtime("unterminated arg list\n%s", tok.String())
		}
		switch curr.TokenType {
		case tokenizer.RComma:
			list = append(list, currList)
			currList = nil
		case tokenizer.RCloseParen:
			list = append(list, currList)
			return list, nil
		default:
			// this could be illegal tokens, but parseValueExpr() will handle them
			currList = append(currList, curr)
		}
	}
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
