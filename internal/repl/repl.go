package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
)

// Start preprocesses all input files and may initiate a new REPL session taking
// input form in and outputting it to out if.
func Start(in io.Reader, out io.Writer, err io.Writer, args Args) {
	scanner := bufio.NewScanner(in)
	reader := read.NewReader()
	parser := read.NewParser()
	env := data.NewEnv(nil)
	eval := eval.NewEvaluator(env)
	printer := print.NewPrinter()

	for _, file := range args.Files {
		res := eval.EvalFile(file)
		printer.FprintErrors(err, eval.Errors())
		printer.Fprint(out, res)
	}

	if args.Interactive {
		for {
			fmt.Fprintf(out, ">> ")
			if ok := scanner.Scan(); !ok {
				return
			}
			input := scanner.Text()
			reader.Load(input)
			src := parser.Parse(reader)
			printer.FprintErrors(err, parser.Errors())
			res := eval.Eval(src)
			printer.FprintErrors(err, eval.Errors())
			printer.Fprint(out, res)
		}
	}
}
