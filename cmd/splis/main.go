package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mhoertnagl/splis2/internal/util"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func main() {
	flag.Parse()

	pr := vm.NewPrinter()
	vm := vm.NewVM(1024, 512, 512)

	vm.AddGlobal(0, os.Stdin)  // *STD-IN*
	vm.AddGlobal(1, os.Stdout) // *STD-OUT*
	vm.AddGlobal(2, os.Stderr) // *STD-ERR*

	for _, inFileName := range flag.Args() {
		inFile, err := os.Open(inFileName)
		if err != nil {
			panic(err)
		}
		vm.Run(util.ReadStatic(inFile))
	}

	fmt.Println("-----")
	for i := int64(0); i < vm.StackSize(); i++ {
		v := vm.InspectStack(i)
		fmt.Println(pr.Print(v))
	}
	fmt.Println("-----")
}
