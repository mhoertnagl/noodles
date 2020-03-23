package main

import (
	"flag"
	"os"

	"github.com/mhoertnagl/splis2/internal/repl"
)

func main() {

	interactive := flag.Bool("i", false, "interactive mode")

	flag.Parse()

	args := repl.Args{
		Interactive: *interactive,
		Files:       flag.Args(),
	}

	repl.Start(os.Stdin, os.Stdout, os.Stderr, args)
}
