// Credit to Thorsten Ball for:
//  - Idea of interfaces for AST and objects.
//      - Actual AST structures/object enumeration by me.
//  - Idea of object system to pass evaluated AST results.
//  - Variable argument error messages.
//      - Handling of all syntax/runtime errors decided by me.
//  - Applying expressions to function calls.
//      - Function scoping/reuse of identifier errors done by me.
//  - Built-in functions interface.
//  - Using a map to hold the contents of a state.

package repl

import (
	"ariel/color"
	"ariel/eval"
	"ariel/misc"
	"ariel/object"
	"ariel/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func REPL(debug bool) {
	prompt := color.Cyan + "ariel>> " + color.Reset

	scanner := bufio.NewScanner(os.Stdin)
	state := object.NewState()

	welcome := "Welcome to the Ariel programming language.\"\n"
	welcome += "\"Type \"ariel\" or \"help\" for more information.\"\n"
	welcome += "\"If you want to exit or quit, just say so."
	fmt.Println(misc.Flounder(welcome))

	easteregg := false
	bye := "Out of the C, wish I could be, programming in Ariel."
	sad := "I guess you're leaving me like Ariel did."

	for {
		fmt.Print(prompt)

		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		switch line {
		case "ariel":
			easteregg = true
			fmt.Println(misc.Flounder("She left me for a dude named Eric."))
		case "help":
			fmt.Println(misc.Flounder("Refer to Cahn's Axiom."))
		case "exit":
			if easteregg {
				fmt.Println(misc.Flounder(sad))
			} else {
				fmt.Println(misc.Flounder(bye))
			}
			os.Exit(0)
		case "quit":
			if easteregg {
				fmt.Println(misc.Flounder(sad))
			} else {
				fmt.Println(misc.Flounder(bye))
			}
			os.Exit(0)
		default:
			program := parser.ParseProgramString(line, debug)
			result := eval.Eval(program, state)
			if result != nil {
				fmt.Println(result.Eval())
			}
		}
	}
}
