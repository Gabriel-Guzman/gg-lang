package main

import (
	"encoding/json"
	"fmt"
	"github.com/gabriel-guzman/gg-lang/src/ggErrs"
	"github.com/gabriel-guzman/gg-lang/src/godTree"
	"github.com/gabriel-guzman/gg-lang/src/program"
	"github.com/gabriel-guzman/gg-lang/src/tokenizer"
	"os"
)

func handle(err error) {
	if err == nil {
		return
	}
	switch err.(type) {
	case *ggErrs.Runtime:
		fmt.Printf("Runtime error: %s\n", err.Error())
	case *ggErrs.Internal:
		panic(fmt.Sprintf("Internal error: %s\n", err.Error()))
	default:
		panic(fmt.Sprintf("Unknown error (please wrap in ggErrs): %v\n", err))
	}
}

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

	stmts, err := tokenizer.BuildStmts([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed %d statements\n", len(stmts))

	ast := godTree.New()
	err = ast.ParseStmts(stmts)
	handle(err)

	sess := program.New()
	err = sess.Run(ast)
	handle(err)

	fmt.Printf("Program: \n%v\n", sess.String())
	tree, err := json.MarshalIndent(ast, "", "    ")
	handle(err)

	err = os.WriteFile("out/ast.json", tree, 0644)
	handle(err)
	fmt.Println("AST saved to out/ast.json")
}
