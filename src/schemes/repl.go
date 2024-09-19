package schemes

import (
	"bufio"
	"fmt"
	"gg-lang/src/program"
	"os"
	"strings"
)

func Repl() {
	sess := program.New()
	fmt.Println("Welcome to the GG programming language!")
	for {
		fmt.Print("gg> ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		text := input.Text()
		text = strings.TrimSpace(text)

		err := sess.RunString(text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
