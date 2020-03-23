package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mhoertnagl/splis2/internal/bin"
)

func main() {

	outFileName := flag.String("o", "out.splis.lib", "output file name")

	flag.Parse()

	lnk := bin.NewLinker()

	for _, inFileName := range flag.Args() {
		inFile, err := os.Open(inFileName)
		if err != nil {
			panic(err)
		}
		lnk.Add(bin.ReadLib(inFile))
	}

	outFile, err := os.Create(*outFileName)
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasSuffix(*outFileName, "splis.lib") {
		bin.WriteLib(lnk.Lib(), outFile)
	}
	if strings.HasSuffix(*outFileName, "splis.bin") {
		bin.WriteStatic(lnk.Code(), outFile)
	}
}
