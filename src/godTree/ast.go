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

func tokStringer(t tokenizer.Token) string {
	return t.Str
}

func (a *Ast) ParseStmts(tokens [][]tokenizer.Token) error {
	iter := iterator.New(tokens)
	var currStmt []tokenizer.Token
	tokIter := iterator.New(currStmt)

	nextStmt := func() bool {
		currStmt, ok := iter.Next()
		if !ok {
			return ok
		}

		tokIter = iterator.New(currStmt)
		tokIter.Stringer = tokStringer
		return true
	}

outer:
	for {
		ok := nextStmt()
		if !ok {
			break
		}
		expr, err := parseStmt(tokIter)
		if err != nil {
			return err
		}
		// trap for function declaration
		if casted, ok := expr.(*FunctionDeclExpression); ok {
			for {
				ok = nextStmt()
				if !ok {
					return ggErrs.Runtime("missing } in function decl\n%s", tokIter.String())
				}

				// handle close brace (end function decl)
				if next, ok := tokIter.Peek(); ok {
					if next.TokenType == tokenizer.RCloseBrace {
						if tokIter.Len() > 1 {
							return ggErrs.Runtime("unexpected expr after }\n%s", tokIter.String())
						}

						a.Body = append(a.Body, casted)
						tokIter.Next()
						continue outer
					}
				}

				funcBodyExpr, err := parseStmt(tokIter)
				if err != nil {
					return err
				}
				casted.Value = append(casted.Value, funcBodyExpr)
			}
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

	if curr.TokenType == tokenizer.Function {

		mbIdent, ok := tokIter.Next()
		if !ok {
			return nil, ggErrs.Runtime("Expected func name\n%s", tokIter.String())
		}

		id, err := newIdentifier(mbIdent)
		if err != nil {
			return nil, err
		}

		mbOpenParen, ok := tokIter.Next()
		if !ok || mbOpenParen.TokenType != tokenizer.ROpenParen {
			return nil, ggErrs.Runtime("Expected (\n%s", tokIter.String())
		}

		var parms []string
		for {
			parm, ok := tokIter.Next()
			if !ok {
				return nil, ggErrs.Runtime("Unexpected end of param list\n%s", tokIter.String())
			}
			if parm.TokenType == tokenizer.RCloseParen {
				break
			}
			if parm.TokenType == tokenizer.RComma {
				continue
			}
			if parm.TokenType != tokenizer.Var {
				return nil, ggErrs.Runtime("Unexpected token\n%s", tokIter.String())
			}

			parms = append(parms, parm.Str)
		}

		mbOpenBrack, ok := tokIter.Next()
		if !ok || mbOpenBrack.TokenType != tokenizer.ROpenBrace {
			return nil, ggErrs.Runtime("Expected {\n%s", tokIter.String())
		}

		return &FunctionDeclExpression{
			Target: *id,
			Parms:  parms,
			Value:  nil,
		}, nil
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

// tokIter should be pointing to the token right before the value expression
func parseAssignmentExpr(id *Identifier, tokIter *iterator.Iter[tokenizer.Token]) (Expression, error) {
	expr, err := parseValueExpr(tokIter)
	if err != nil {
		return nil, err
	}

	at := newAssignmentExpression(id, expr)
	return at, nil
}

// tokIter should be pointing to the token right before the value expression
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
	return firstExpr, nil
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
