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
		if tokParser.Curr.TokenType == token.RCloseBrace {
			if tokParser.HasNext {
				return nil, ggErrs.Crit("unexpected token after closing brace - check tokenizer\n%s", tokParser.String())
			}

			a.StmtPar.Advance() // consume closing brace, which should be its own statement in this version
			return expressions, nil
		}

		stmt, err := parseTopLevelExpr(tokParser)
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
	for _, tokens := range ins {
		newStmtParser := parser.New(tokens)
		newStmtParser.SetStringer(tokStringer)
		stmtParserMembers = append(stmtParserMembers, newStmtParser)
	}
	newParser := parser.New(stmtParserMembers)
	newParser.SetSeparator("")
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

func BuildFromStatements(ins [][]token.Token) (*Ast, error) {
	a := newAstBuilder(ins)

	var expressions []Expression
	for a.StmtPar.HasCurr {
		stmt, err := parseTopLevelExpr(a.StmtPar.Curr)
		if err != nil {
			return nil, err
		}
		if !a.StmtPar.Curr.IsDone() {
			return nil, ggErrs.Crit("could not finish parsing statement\n%s", a.StmtPar.String())
		}
		expressions = append(expressions, stmt)
		a.StmtPar.Advance() // consume the declaration

		// function trap, note continue at end of this block
		if stmt.Kind() == ExprFuncDecl {
			// function trap
			decl := stmt.(*FunctionDeclExpression)
			exprs, err := a.parseBlockStatement()
			if err != nil {
				return nil, err
			}
			decl.Value = exprs

			expressions = append(expressions, decl)
			continue
		}

	}

	return &Ast{Body: expressions}, nil
}
