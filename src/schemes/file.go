package schemes

import (
	"encoding/json"
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/gg_ast"
	"gg-lang/src/program"
	"gg-lang/src/token"
	"os"
)

// execute a GG program from a file and output the AST to a file
func Exec(filename string) {
	fmt.Println("Reading file:", filename)
	out, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	// tokenize the input manually so we can save the tokens to a file for debugging
	stmts, err := token.TokenizeRunes([]rune(string(out)))
	ggErrs.Handle(err)

	stmtsJson, err := json.MarshalIndent(stmts, "", "    ")
	ggErrs.Handle(err)

	err = os.WriteFile("out/stmts.json", stmtsJson, 0644)
	ggErrs.Handle(err)

	ast, err := gg_ast.BuildFromStatements(stmts)
	ggErrs.Handle(err)

	tree, err := json.MarshalIndent(ast, "", "    ")
	ggErrs.Handle(err)

	err = os.WriteFile("out/ast.json", tree, 0644)
	ggErrs.Handle(err)

	fmt.Println("AST saved to out/ast.json")

	fmt.Println("Running program...")
	sess := program.New()
	err = sess.Run(ast)
	ggErrs.Handle(err)
}
