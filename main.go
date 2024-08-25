package main

import (
	"fmt"
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

	tokens, err := BuildTokens([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BuildTokens output: \n%v\n", tokens)

	ast, err := newAST(tokens)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AST: \n%v\n", ast.String())

	sess := &session{variables: make(map[string]variable), omap: defaultOpMap()}
	err = sess.run(ast)
	if err != nil {
		fmt.Println("Error running session:", err)
		return
	}

	fmt.Println("Final variables:")
	for k, v := range sess.variables {
		fmt.Printf("%s = %v\n", k, v.value)
	}
}
