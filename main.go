//go:generate cp parser/parser.y parser.y
//go:generate goyacc -o parser.go parser.y
//go:generate mv parser.go parser/.
//go:generate rm -f parser.y y.output

package main

import (
	"ariel/eval"
	"ariel/object"
	"ariel/parser"
	"ariel/repl"
	"flag"
	"fmt"
	"os"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode.")
	replit := flag.Bool("repl", false, "Enable the REPL.")
	flag.Parse()

	if *replit || len(os.Args) == 1 {
		repl.REPL(*debug)
	} else {
		infile, err := os.Open(os.Args[len(os.Args)-1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to open input.\n")
			os.Exit(1)
		}

		program := parser.ParseProgram(infile, *debug)
		state := object.NewState()
		result := eval.Eval(program, state)
		if eval.IsError(result) {
			fmt.Println(result.Eval())
		}
	}
}
