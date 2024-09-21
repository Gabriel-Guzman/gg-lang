package gg_ast

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/operators"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

/*
This is the top-level expression parser.
Its job is to
 1. pick the appropriate expression parser (all defined below) based on

the first two tokens in the parser argument.
 2. parse the expression
 3. handle block statements (if any)

The parser should be at its initial state, i.e. with the current index set to 0 and
pointing to the first token in the expression.
After a successful parse, the parser should be pointing to the token after the expression
*/
func parseTopLevelExpr(p *parser.Parser[token.Token]) (Expression, error) {
	// first expression should be an identifier, a reserved keyword, or a value expression
	// check reserved keywords first
	if p.Curr.TokenType == token.Function {
		funcDecl, err := _parseFuncDecl(p)
		if err != nil {
			return nil, err
		}
		return funcDecl, nil
	}

	// now it could be a function call or an assignment expression, both of which
	// have to start with an identifier. this means no unassigned value expressions
	// other than function calls are allowed at the top level.
	id, err := _parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == token.RAssign {
		p.Advance()
		return _parseAssignmentExpr(id, p)
	}

	// if no operator, it's a function call
	return _parseFuncCallExpr(id, p)
}

/*
The parser argument to an expression parser should
be pointing to the first token in the expression to be parsed,
meaning parser.Curr == <first relevant token>
This will almost always be the token right after the first identifier in the expression
i.e. the "(" in "funCall();" or the "=" in "x = 1 + 2;"
After a successful parse, the parser should be pointing to the token after the expression
*/
func _parseAssignmentExpr(target *Identifier, p *parser.Parser[token.Token]) (*AssignmentExpression, error) {
	expr, err := _parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	return &AssignmentExpression{Target: target, Value: expr}, nil
}

func _parseFuncDecl(p *parser.Parser[token.Token]) (*FunctionDeclExpression, error) {
	p.Advance() // eat the function keyword

	id, err := _parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.ROpenParen) {
		return nil, ggErrs.Runtime("expected '(' after function name\n%s", p.String())
	}

	var params []string
	for {
		if !p.HasCurr {
			return nil, ggErrs.Runtime("Unexpected end of param list\n%s", p.String())
		}
		param := p.Curr
		if param.TokenType == token.RCloseParen {
			break
		}
		if param.TokenType == token.RComma {
			continue
		}
		if param.TokenType != token.Ident {
			return nil, ggErrs.Runtime("Unexpected token\n%s", p.String())
		}

		params = append(params, param.Str)
		p.Advance()
	}

	if !advanceIfCurrIs(p, token.ROpenBrace) {
		return nil, ggErrs.Runtime("expected '{' after function declaration\n%s", p.String())
	}

	return &FunctionDeclExpression{
		Target: id,
		Params: params,
	}, nil
}

func _parseIdentifier(p *parser.Parser[token.Token]) (*Identifier, error) {
	t := p.Curr
	var ik IdExprKind
	switch t.TokenType {
	case token.IntLiteral:
		ik = IdExprNumber
	case token.Ident:
		ik = IdExprVariable
	case token.StringLiteral:
		ik = IdExprString
	case token.TrueLiteral:
		ik = IdExprBool
	case token.FalseLiteral:
		ik = IdExprBool
	default:
		return nil, ggErrs.Runtime("invalid identifier %s", t.Str)
	}
	p.Advance()
	return &Identifier{Raw: t.Str, idKind: ik}, nil
}

// returns a simple value expression or a binary expression
func _parseValueExpr(p *parser.Parser[token.Token]) (IValExpr, error) {
	if !p.HasCurr {
		return nil, ggErrs.Runtime("unexpected end of expression\n%s", p.String())
	}

	// build initial binary tree
	lhsNonBinary, err := _parseSingleValueExpr(p)
	if err != nil {
		return nil, err
	}

	if !p.Curr.TokenType.IsOperator() {
		return lhsNonBinary, nil
	}

	op := p.Curr

	rhs, err := _parseSingleValueExpr(p)
	if err != nil {
		return nil, err
	}

	lhs := &BinaryExpression{
		Lhs: lhsNonBinary,
		Op:  op.Str,
		Rhs: rhs,
	}

	// add on to the initial tree
	for {
		op = p.Curr
		if !op.TokenType.IsMathOperator() {
			break
		}

		rhs, err := _parseSingleValueExpr(p)
		if err != nil {
			return nil, err
		}

		if operators.LeftFirst(lhs.Op, op.Str) {
			// left needs to be evaluated first and therefor deeper into the tree
			lhs = &BinaryExpression{
				Lhs: lhs,
				Op:  op.Str,
				Rhs: rhs,
			}
		} else {
			lhs.Rhs = &BinaryExpression{
				Lhs: lhs.Rhs,
				Op:  op.Str,
				Rhs: rhs,
			}
		}
	}
	return lhs, nil
}

// A single value expression is either an identifier, a literal, or a function call
func _parseSingleValueExpr(p *parser.Parser[token.Token]) (IValExpr, error) {
	if !p.HasCurr {
		return nil, ggErrs.Runtime("unexpected end of expression\n%s", p.String())
	}

	// unary operators
	if p.Curr.TokenType == token.RMinus {
		p.Advance()
		id, err := _parseIdentifier(p)
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Lhs: &Identifier{
				Raw:    "-1",
				idKind: IdExprNumber,
			},
			Op:  "*",
			Rhs: id,
		}, nil
	} else {
		id, err := _parseIdentifier(p)
		if err != nil {
			return nil, err
		}

		if p.Curr.TokenType == token.ROpenParen {
			return _parseFuncCallExpr(id, p)
		} else {
			return id, nil
		}
	}
}

func _parseFuncCallExpr(id *Identifier, p *parser.Parser[token.Token]) (IValExpr, error) {
	if !advanceIfCurrIs(p, token.ROpenParen) {
		return nil, ggErrs.Runtime("expected '(' after function name\n%s", p.String())
	}

	var args []IValExpr
	for {
		if !p.HasCurr {
			return nil, ggErrs.Runtime("Unexpected end of arg list\n%s", p.String())
		}
		arg := p.Curr
		if arg.TokenType == token.RCloseParen {
			break
		}
		if arg.TokenType == token.RComma {
			continue
		}
		expr, err := _parseValueExpr(p)
		if err != nil {
			return nil, err
		}

		args = append(args, expr)
		p.Advance()
	}

	return &FunctionCallExpression{
		Id:   id,
		Args: args,
	}, nil
}

// Advances the parser if the current token matches the given token type
func advanceIfCurrIs(p *parser.Parser[token.Token], tt token.TokenType) bool {
	if p.HasCurr && p.Curr.TokenType == tt {

		p.Advance()
		return true
	}

	return false
}
