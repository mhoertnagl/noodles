package repl

import (
	"bufio"
	"fmt"
	//"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"io"
)

// TODO: Repl struct with parameters for the header text, input prefix, ...

// Start initiates a new REPL session taking input form in and outputting
// it to out. The parameter args modifies the behavior of the REPL.
func Start(in io.Reader, out io.Writer, args Args) {
	scanner := bufio.NewScanner(in)
	reader := read.NewReader()
	parser := read.NewParser()
	//env := eval.NewEnv()
	//eval := eval.NewEvaluator()
	printer := print.NewPrinter()

	for {
		fmt.Fprintf(out, ">> ")
		if ok := scanner.Scan(); !ok {
			return
		}
		input := scanner.Text()
		// if input == ":x" {
		// 	fmt.Fprintf(out, "Bye.\n")
		// 	return
		// }
		reader.Load(input)
		ast := parser.Parse(reader)
		output := printer.Print(ast)
		fmt.Fprintf(out, "%s\n", output)
		errors := printer.PrintErrors(parser)
		fmt.Fprintf(out, "\n%s", errors)
	}
}
