package gg_ast

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/parser"
	"gg-lang/src/token"
)

type astBuilder struct {
	StmtPar *parser.Parser[*parser.Parser[token.Token]]
}

func (a *astBuilder) parseBlockStatement() ([]Expression, error) {
	var expressions []Expression
	for a.StmtPar.HasCurr {
		tokParser := a.StmtPar.Curr
		if tokParser.Curr.TokenType == token.CloseBrace {
			if tokParser.HasNext {
				return nil, ggErrs.Crit("unexpected token after closing brace - check tokenizer\n%s", tokParser.String())
			}

			a.StmtPar.Advance() // consume closing brace, which should be its own statement in this version
			return expressions, nil
		}

		stmt, err := parseStatement(tokParser)
		if _, ok := stmt.(*FunctionDeclExpression); ok {
			return nil, ggErrs.Runtime("function declaration inside block statement is not allowed\n%s", tokParser.String())
		}
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, stmt)

		a.StmtPar.Advance()
	}

	return nil, ggErrs.Syntax("no closing brace for block statement\n%s", a.StmtPar.String())
}

func newAstBuilder(ins [][]token.Token) *astBuilder {
	var stmtParserMembers []*parser.Parser[token.Token]
	// astBuilder has a list of parsers which each have a parser, set those up
	// TODO load parsers lazily
	for _, tokens := range ins {
		newStmtParser := parser.New(tokens)
		newStmtParser.SetStringer(tokStringer)
		stmtParserMembers = append(stmtParserMembers, newStmtParser)
	}
	newParser := parser.New(stmtParserMembers)
	// due to the stringer below, only one member will print and doesn't need a separator
	newParser.SetSeparator("")
	// set up the stringer to only show the current statement being parsed, not the entire input
	newParser.SetStringer(func(in *parser.Parser[token.Token]) string {
		if newParser.Curr == in {
			return in.String()
		} else {
			return ""
		}
	})

	return &astBuilder{
		StmtPar: newParser,
	}
}

func BuildFromString(ins string) (*Ast, error) {
	stmts, err := token.TokenizeRunes([]rune(ins))
	if err != nil {
		return nil, err
	}

	return BuildFromStatements(stmts)
}

func BuildFromStatements(ins [][]token.Token) (*Ast, error) {
	a := newAstBuilder(ins)

	var expressions []Expression
	for a.StmtPar.HasCurr {
		stmt, err := parseStatement(a.StmtPar.Curr)
		if err != nil {
			return nil, err
		}
		if !a.StmtPar.Curr.IsDone() {
			return nil, ggErrs.Crit("could not finish parsing statement\n%s", a.StmtPar.String())
		}
		a.StmtPar.Advance() // consume the statement

		// if the last statement was a function declaration, parse its block statement
		if stmt.Kind() == ExprFuncDecl {
			exprs, err := a.parseBlockStatement()
			if err != nil {
				return nil, err
			}

			decl := stmt.(*FunctionDeclExpression)
			decl.Value = exprs
		}

		if stmt.Kind() == ExprForLoop {
			exprs, err := a.parseBlockStatement()
			if err != nil {
				return nil, err
			}

			loop := stmt.(*ForLoopExpression)
			loop.Body = exprs
		}

		expressions = append(expressions, stmt)
	}

	return &Ast{Body: expressions}, nil
}
