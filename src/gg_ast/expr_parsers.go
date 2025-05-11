package gg_ast

import (
	"gg-lang/src/gg"
	"gg-lang/src/operators"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

/*
This is the top-level expression parser.
Its job is to
 1. pick the appropriate expression parser (all defined below) based on the
		first two tokens in the parser argument.
 2. parse the expression

The parser should be at its initial state, i.e. with the current index set to 0 and
pointing to the first token in the expression.
After a successful parse, the parser should be pointing to the token after the expression
*/

type tokenParser = *parser.Parser[token.Token]

// a statement is a function call, declaration, assignment expression, or a for loop expression
// it is up to the builder to disallow these expressions if it's not parsing the top level
func parseExpression(p tokenParser) (Expression, error) {
	// first expression should be an identifier, a reserved keyword, or a value expression
	// check reserved keywords first
	if p.Curr.TokenType == token.Function {
		return parseFuncDecl(p)
	}
	if p.Curr.TokenType == token.For {
		return parseForLoopExpr(p)
	}
	if p.Curr.TokenType == token.If {
		return parseIfElseExpr(p)
	}
	if p.Curr.TokenType == token.Return {
		return parseReturnExpr(p)
	}

	// now it could be a function call or an assignment expression, both of which
	// have to start with an identifier. this means no unassigned value expressions
	// other than function calls are allowed at the top level.
	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == token.Assign {
		p.Advance()
		return parseAssignmentExpr(id, p)
	}

	// if no operator, it's a function call
	expr, err := parseFuncCallExpr(id, p)
	if err != nil {
		return nil, err
	}
	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after function call\n%s", p.String())
	}
	return expr, nil
}

/*
The parser argument to an expression parser should
be pointing to the first token in the expression to be parsed,
meaning parser.Curr == <first relevant token>
This will almost always be the token right after the first identifier in the expression
i.e. the "(" in "funCall();" or the "=" in "x = 1 + 2;"
After a successful parse, the parser should be pointing to the token after the expression
*/
func parseAssignmentExpr(target *Identifier, p tokenParser) (*AssignmentExpression, error) {
	expr, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after assignment expression\n%s", p.String())
	}

	return &AssignmentExpression{Target: target, Value: expr}, nil
}

func parseForLoopExpr(p tokenParser) (*ForLoopExpression, error) {
	if !advanceIfCurrIs(p, token.For) { // eat the for keyword
		return nil, gg.Crit("expected 'for' keyword in expression parser\n%s", p.String())
	}

	condition, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	body, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}
	return &ForLoopExpression{Condition: condition, Body: body}, nil
}

func parseIfElseExpr(p tokenParser) (*IfElseStatement, error) {
	if !advanceIfCurrIs(p, token.If) { // eat the if keyword
		return nil, gg.Crit("expected 'if' keyword in expression parser\n%s", p.String())
	}
	res := &IfElseStatement{}
	condition, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}
	res.Condition = condition

	body, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}
	res.Body = body

	if advanceIfCurrIs(p, token.Else) {
		if p.Curr.TokenType == token.If {
			alt, err := parseIfElseExpr(p)
			if err != nil {
				return nil, err
			}
			res.ElseExpression = alt
		} else {
			b, err := parseBlockStatement(p)
			if err != nil {
				return nil, err
			}
			res.ElseExpression = b
		}
	}

	return res, nil
}

func parseReturnExpr(p tokenParser) (*ReturnStatement, error) {
	if !advanceIfCurrIs(p, token.Return) { // eat the return keyword
		return nil, gg.Crit("expected 'return' keyword in expression parser\n%s", p.String())
	}

	expr, err := parseValueExpr(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.Term) {
		return nil, gg.Syntax("expected ; after return expression\n%s", p.String())
	}

	return &ReturnStatement{Value: expr}, nil
}

func parseFuncDecl(p tokenParser) (*FunctionDeclExpression, error) {
	p.Advance() // eat the function keyword

	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if !advanceIfCurrIs(p, token.OpenParen) {
		return nil, gg.Runtime("expected '(' after function name\n%s", p.String())
	}

	var params []token.Token
	for {
		if !p.HasCurr {
			return nil, gg.Runtime("Unexpected end of param list\n%s", p.String())
		}
		param := p.Curr
		if param.TokenType == token.CloseParen {
			p.Advance() // eat the closing parenthesis ')'
			break
		}
		if param.TokenType == token.Comma {
			p.Advance() // eat the comma
			continue
		}
		if param.TokenType != token.Ident {
			return nil, gg.Runtime("Unexpected token\n%s", p.String())
		}

		params = append(params, param)
		p.Advance()
	}

	block, err := parseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	return &FunctionDeclExpression{
		Target: id,
		Params: params,
		Body:   block,
	}, nil
}

func parseIdentifier(p tokenParser) (*Identifier, error) {
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
		return nil, gg.Runtime("invalid identifier %s", t.Symbol)
	}
	p.Advance()
	return &Identifier{Tok: t, idKind: ik}, nil
}

// returns a primary expression or a binary expression
func parseValueExpr(p tokenParser) (ValueExpression, error) {
	if !p.HasCurr {
		return nil, gg.Runtime("unexpected end of expression\n%s", p.String())
	}

	// build initial binary tree
	lhsNonBinary, err := parsePrimaryExpr(p)
	if err != nil {
		return nil, err
	}

	if !p.Curr.TokenType.IsOperator() {
		return lhsNonBinary, nil
	}

	op := p.Curr
	p.Advance() // eat the operator token

	rhs, err := parsePrimaryExpr(p)
	if err != nil {
		return nil, err
	}

	lhs := &BinaryExpression{
		Lhs: lhsNonBinary,
		Op:  op,
		Rhs: rhs,
	}

	// add on to the initial tree
	for {
		op = p.Curr
		if !op.TokenType.IsOperator() {
			break
		}
		p.Advance() // eat the operator token

		rhs, err := parsePrimaryExpr(p)
		if err != nil {
			return nil, err
		}

		if operators.LeftFirst(lhs.Op.Symbol, op.Symbol) {
			// left needs to be evaluated first and therefor deeper into the tree
			lhs = &BinaryExpression{
				Lhs: lhs,
				Op:  op,
				Rhs: rhs,
			}
		} else {
			lhs.Rhs = &BinaryExpression{
				Lhs: lhs.Rhs,
				Op:  op,
				Rhs: rhs,
			}
		}
	}
	return lhs, nil
}

// A primary expression is either an identifier, a literal, a function call, a unary binary expression, or a function declaration
func parsePrimaryExpr(p tokenParser) (ValueExpression, error) {
	if !p.HasCurr {
		return nil, gg.Runtime("unexpected end of expression\n%s", p.String())
	}

	// unary operators
	if p.Curr.TokenType.IsOperator() {
		op := p.Curr
		p.Advance()
		toNegate, err := parsePrimaryExpr(p)
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{
			Op:  op,
			Rhs: toNegate,
		}, nil
	}

	if p.Curr.TokenType == token.Function {
		return parseFuncDecl(p)
	}

	id, err := parseIdentifier(p)
	if err != nil {
		return nil, err
	}

	if p.Curr.TokenType == token.OpenParen {
		return parseFuncCallExpr(id, p)
	} else {
		return id, nil
	}

}

func parseFuncCallExpr(id *Identifier, p tokenParser) (ValueExpression, error) {
	if !advanceIfCurrIs(p, token.OpenParen) {
		return nil, gg.Runtime("expected '(' after function name\n%s", p.String())
	}
	if p.Curr.TokenType == token.CloseParen {
		p.Advance() // consume the ')'
		return &FunctionCallExpression{Id: id}, nil
	}

	var args []ValueExpression
	for {
		if !p.HasCurr {
			return nil, gg.Runtime("Unexpected end of arg list\n%s", p.String())
		}

		expr, err := parseValueExpr(p)
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if advanceIfCurrIs(p, token.CloseParen) { // consume the ')'
			break
		}
		if !advanceIfCurrIs(p, token.Comma) {
			return nil, gg.Runtime("expected ',' or ')' after argument\n%s", p.String())
		}
	}

	return &FunctionCallExpression{
		Id:   id,
		Args: args,
	}, nil
}

// Advances the parser if the current token matches the given token type
func advanceIfCurrIs(p tokenParser, tt token.Type) bool {
	return p.AdvanceIf(func(t token.Token) bool { return t.TokenType == tt })
}
