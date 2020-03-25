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
	// TODO: Add rewriters (macros)
	// TODO: Feed macros from referenced libs to macro rewriter.

	rdr := compiler.NewReader()
	prs := compiler.NewParser()
	cmp := compiler.NewCompiler()
	qrw := compiler.NewQuoteRewriter()

	for _, inFileName := range flag.Args() {
		inFileBytes, err := ioutil.ReadFile(inFileName)
		if err != nil {
			panic(err)
		}

		rdr.Load(string(inFileBytes))
		n := prs.Parse(rdr)
		n = qrw.Rewrite(n)
		lib := cmp.CompileLib(n)

		outFile, err := os.Create(inFileName + ".lib")
		if err != nil {
			panic(err)
		}

		bin.WriteLib(lib, outFile)
	}
}
