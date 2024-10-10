package gg_ast

import (
	"fmt"
)

func LeftFirst(l, r string) bool {
	pl, ok := PrecedenceMap[l]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s", l))
	}
	pr, ok := PrecedenceMap[r]
	if !ok {
		panic(fmt.Sprintf("checked precedence on nonexistent op %s", r))
	}

	if pl == pr {
		return true
	}

	return pl > pr
}
