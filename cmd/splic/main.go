package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/mhoertnagl/splis2/internal/bin"
	"github.com/mhoertnagl/splis2/internal/compiler"
)

func main() {
	flag.Parse()

	// TODO: Load Usings.
	// TODO: Add rewriters (quotes, macros)
	// TODO: Feed macros from referenced libs to macro rewriter.

	rdr := compiler.NewReader()
	prs := compiler.NewParser()
	cmp := compiler.NewCompiler()

	for _, inFileName := range flag.Args() {
		inFileBytes, err := ioutil.ReadFile(inFileName)
		if err != nil {
			panic(err)
		}

		rdr.Load(string(inFileBytes))
		n := prs.Parse(rdr)
		l := cmp.Compile2(n)

		outFile, err := os.Create(inFileName + ".lib")
		if err != nil {
			panic(err)
		}

		bin.WriteLib(l, outFile)
	}
}
