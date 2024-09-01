package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gg-lang/src/ggErrs"
	"gg-lang/src/godTree"
	"gg-lang/src/tokenizer"
	"os"
)

func handle(err error) {
	if err == nil {
		return
	}
	var chillErr *ggErrs.ChillErr
	var critErr *ggErrs.CritErr
	switch {
	case errors.As(err, &chillErr):
		fmt.Printf("Runtime error: %s\n", err.Error())
	case errors.As(err, &critErr):
		panic(fmt.Sprintf("Crit error: %s\n", err.Error()))
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

	stmts, err := tokenizer.TokenizeRunes([]rune(string(out)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed %d statements\n", len(stmts))
	fmt.Printf("stmts: %+v\n", stmts)

	ast := godTree.New()
	err = ast.ParseStmts(stmts)
	handle(err)
	//
	//sess := program.New()
	//err = sess.Run(ast)
	//handle(err)
	//
	//fmt.Printf("Program: \n%v\n", sess.String())
	tree, err := json.MarshalIndent(ast, "", "    ")
	handle(err)
	//
	err = os.WriteFile("out/ast.json", tree, 0644)
	handle(err)
	fmt.Println("AST saved to out/ast.json")
}
