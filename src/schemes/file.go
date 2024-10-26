package schemes

import (
	"encoding/json"
	"fmt"
	"gg-lang/src/gg"
	"gg-lang/src/gg_ast"
	"gg-lang/src/program"
	"gg-lang/src/token"
	"os"
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// execute a GG program from a file and output the AST to a file
func Exec(filename string) {
	t := makeTimestamp()
	fmt.Println("Reading file:", filename)
	out, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	// tokenize the input manually so we can save the tokens to a file for debugging
	stmts, err := token.TokenizeRunes([]rune(string(out)))
	gg.Handle(err)

	stmtsJson, err := json.MarshalIndent(stmts, "", "    ")
	gg.Handle(err)

	err = os.WriteFile("out/stmts.json", stmtsJson, 0644)
	gg.Handle(err)

	ast, err := gg_ast.BuildFromTokens(stmts)
	gg.Handle(err)

	tree, err := json.MarshalIndent(ast, "", "    ")
	gg.Handle(err)

	err = os.WriteFile("out/ast.json", tree, 0644)
	gg.Handle(err)

	fmt.Println("AST saved to out/ast.json")

	fmt.Println("Running program...")
	sess := program.New()
	err = sess.Run(ast)
	gg.Handle(err)

	fmt.Println(makeTimestamp()-t, "ms")
}
