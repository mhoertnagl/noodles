package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
)

// TODO: Repl struct with parameters for the header text, input prefix, ...
// TODO: Unit tests.

// Start initiates a new REPL session taking input form in and outputting
// it to out. The parameter args modifies the behavior of the REPL.
func Start(in io.Reader, out io.Writer, args Args) {
	scanner := bufio.NewScanner(in)
	reader := read.NewReader()
	parser := read.NewParser()
	env := eval.NewEnv(nil)
	eval := eval.NewEvaluator(env)
	printer := print.NewPrinter()

	for {
		fmt.Fprintf(out, ">> ")
		if ok := scanner.Scan(); !ok {
			return
		}
		// TODO: Print environment.
		input := scanner.Text()
		reader.Load(input)
		src := parser.Parse(reader)
		errors := printer.PrintErrors(parser.Errors())
		fmt.Fprintf(out, "\n%s", errors)
		res := eval.Eval(src)
		output := printer.Print(res)
		fmt.Fprintf(out, "%s\n", output)
	}
}
