package main

import (
	"flag"
	"os"
	"repl"
)

func main() {

	lexOnly := flag.Bool("l", false, "a bool")
	parseOnly := flag.Bool("p", false, "a bool")

	flag.Parse()

	if len(flag.Args()) == 0 {

	} else {

	}
}
