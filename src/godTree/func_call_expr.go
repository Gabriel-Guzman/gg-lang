package godTree

import (
	"gg-lang/src/ggErrs"
	"gg-lang/src/iterator"
	"gg-lang/src/tokenizer"
)

// iter should be pointing to the opening parenthesis here
func newFuncCallExpression(funcName *Identifier, iter *iterator.Iter[tokenizer.Token]) (ValueExpression, error) {
	nextTok, ok := iter.Peek()
	if !ok {
		return nil, ggErrs.Runtime("expected closing parenthesis ')' or args after function name\n%s", iter.String())
	}
	if nextTok.TokenType == tokenizer.RCloseParen {
		iter.Next() // consume the closing parenthesis ')'
		return &FunctionCallExpression{
			Id:   *funcName,
			Args: nil,
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
		Id:   *funcName,
		Args: args,
	}, nil
}
