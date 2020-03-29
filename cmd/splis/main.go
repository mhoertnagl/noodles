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
	vm := vm.New(1024, 512, 256, 128)

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
