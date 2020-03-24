package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mhoertnagl/splis2/internal/bin"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func main() {

	// interactive := flag.Bool("i", false, "interactive mode")

	// flag.Parse()

	// args := repl.Args{
	// 	Interactive: *interactive,
	// 	Files:       flag.Args(),
	// }
	//
	// repl.Start(os.Stdin, os.Stdout, os.Stderr, args)

	flag.Parse()

	vm := vm.New(1024, 512, 256, 128)

	for _, inFileName := range flag.Args() {
		inFile, err := os.Open(inFileName)
		if err != nil {
			panic(err)
		}
		vm.Run(bin.ReadStatic(inFile))
	}

	fmt.Println("-----")
	for i := int64(0); i < vm.StackSize(); i++ {
		fmt.Printf("  %v\n", vm.InspectStack(i))
	}
	fmt.Println("-----")
}
