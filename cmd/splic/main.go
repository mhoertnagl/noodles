package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mhoertnagl/splis2/internal/asm"
	"github.com/mhoertnagl/splis2/internal/cmp"
	"github.com/mhoertnagl/splis2/internal/rwr"
	"github.com/mhoertnagl/splis2/internal/util"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		panic("provide exactly one input file")
	}

	srcPath := args[0]

	// The compiler will search for used modules in '$(SPLIS_HOME)/lib' and the
	// directory that contians the input source file.
	dirs := []string{
		util.SplisLibPath(),
		filepath.Dir(srcPath),
	}

	rdr := cmp.NewReader()
	prs := cmp.NewParser()
	urw := rwr.NewUseRewriter(dirs)
	qrw := rwr.NewQuoteRewriter()
	mrw := rwr.NewMacroRewriter()
	cmp := cmp.NewCompiler()
	asm := asm.NewAssembler()

	cmp.AddDefaultGlobals()

	srcBytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if len(srcBytes) == 0 {
		fmt.Println("source file is empty")
		os.Exit(-1)
	}

	rdr.Load(string(srcBytes))
	n := prs.Parse(rdr)

	if len(prs.Errors()) > 0 {
		for _, err := range prs.Errors() {
			fmt.Println(err.Msg)
		}
		os.Exit(-1)
	}

	n = urw.Rewrite(n)
	n = qrw.Rewrite(n)
	n = mrw.Rewrite(n)
	a := cmp.Compile(n)
	c := asm.Assemble(a)

	outPath := util.FilePathWithoutExt(srcPath)
	outFile, err := os.Create(outPath + ".nob")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	util.WriteStatic(c, outFile)
}
