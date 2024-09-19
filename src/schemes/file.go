package schemes

import (
	"encoding/json"
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/program"
	"gg-lang/src/tokenizer"
	"os"
)

// execute a GG program from a file and output the AST to a file
func Exec(filename string) {
	fmt.Println("Reading file:", filename)
	out, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	stmts, err := tokenizer.TokenizeRunes([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed %d statements\n\n", len(stmts))
	stmtsJson, err := json.MarshalIndent(stmts, "", "    ")
	ggErrs.Handle(err)
	err = os.WriteFile("out/stmts.json", stmtsJson, 0644)
	ggErrs.Handle(err)

	ast := godTree.New()
	err = ast.ParseStmts(stmts)
	ggErrs.Handle(err)
	if err != nil {
		return
	}

	fmt.Println("\nRunning program...")
	sess := program.New()
	err = sess.Run(ast)
	ggErrs.Handle(err)

	tree, err := json.MarshalIndent(ast, "", "    ")
	ggErrs.Handle(err)

	err = os.WriteFile("out/ast.json", tree, 0644)
	ggErrs.Handle(err)
	fmt.Println("AST saved to out/ast.json")
}
