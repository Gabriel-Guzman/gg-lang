package main

import (
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src"
	"github.com/gabriel-guzman/gg-lang/src/godTree"
	"github.com/gabriel-guzman/gg-lang/src/operators"
	"github.com/gabriel-guzman/gg-lang/src/tokenizer"
	"github.com/gabriel-guzman/gg-lang/src/variables"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: go run main.go <filename>")
	}

	// get arguments
	filename := os.Args[1]
	fmt.Println("Reading file:", filename)
	out, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	tokens, err := tokenizer.BuildTokens([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BuildTokens output: \n%v\n", tokens)

	ast, err := godTree.NewAST()
	err = ast.ExecStmts(tokens)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AST: %v\n", ast.String())

	sess := &src.Session{Variables: make(map[string]variables.Variable), Omap: operators.DefaultOpMap()}
	err = sess.Run(ast)
	if err != nil {
		fmt.Println("Error running session:", err)
		return
	}

	fmt.Println("Final variables:")
	for k, v := range sess.Variables {
		fmt.Printf("%s = %v\n", k, v.Value)
	}
}
