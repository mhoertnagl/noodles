package main

import (
	"flag"
	"github.com/mhoertnagl/splis2/internal/repl"
	"os"
)

func main() {

	lexOnly := flag.Bool("l", false, "a bool")
	parseOnly := flag.Bool("p", false, "a bool")

	flag.Parse()

	args := repl.Args{
		LexOnly:   *lexOnly,
		ParseOnly: *parseOnly,
	}

	if len(flag.Args()) == 0 {
		repl.Start(os.Stdin, os.Stdout, args)
	} else {

	}
}
