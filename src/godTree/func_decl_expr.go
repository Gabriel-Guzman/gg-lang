package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/tokenizer"
)

func parseFuncDecl(tokIter *iterator.Iter[tokenizer.Token]) (*FunctionDeclExpression, error) {
	mbIdent, ok := tokIter.Next()
	if !ok {
		return nil, ggErrs.Runtime("Expected func name\n%s", tokIter.String())
	}

	id, err := newIdentifier(mbIdent)
	if err != nil {
		return nil, ggErrs.Runtime("Invalid identifier\n%s", tokIter.String())
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
