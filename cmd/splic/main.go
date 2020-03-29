package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mhoertnagl/splis2/internal/compiler"
	"github.com/mhoertnagl/splis2/internal/util"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		panic("Provide exactly one input file.")
	}

	srcPath := args[0]

	// The compiler will search for used modules in '$(SPLIS_HOME)/lib' and the
	// directory that contians the input source file.
	dirs := []string{
		util.SplisLibPath(),
		filepath.Dir(srcPath),
	}

	rdr := compiler.NewReader()
	prs := compiler.NewParser()
	urw := compiler.NewUseRewriter(dirs)
	qrw := compiler.NewQuoteRewriter()
	mrw := compiler.NewMacroRewriter()
	cmp := compiler.NewCompiler()

	srcBytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		panic(err)
	}

	rdr.Load(string(srcBytes))
	n := prs.Parse(rdr)
	n = urw.Rewrite(n)
	n = qrw.Rewrite(n)
	n = mrw.Rewrite(n)
	code := cmp.Compile(n)

	outPath := util.FilePathWithoutExt(srcPath)
	outFile, err := os.Create(outPath + ".splin")
	if err != nil {
		panic(err)
	}

	util.WriteStatic(code, outFile)
}
