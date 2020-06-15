package main

import (
	"flag"
	"os"

	"github.com/mhoertnagl/noodles/internal/util"
	"github.com/mhoertnagl/noodles/internal/vm"
)

func main() {
	flag.Parse()

	vm := vm.NewVM(1024, 512, 512)
	vm.AddDefaultGlobals()

	for _, inFileName := range flag.Args() {
		inFile, err := os.Open(inFileName)
		if err != nil {
			panic(err)
		}
		vm.Run(util.ReadStatic(inFile))
	}
}
