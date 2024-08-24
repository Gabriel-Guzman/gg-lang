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

	inter, err := TokenizeRunes([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("TokenizeRunes output: \n%v\n", inter)

	ast, err := fromTokens(inter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AST: \n%v\n", ast.String())

	sess := &session{variables: make(map[string]variable)}
	err = sess.run(ast)
	if err != nil {
		fmt.Println("Error running session:", err)
		return
	}

	fmt.Println("Final variables:", sess.variables)
}
