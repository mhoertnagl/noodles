package repl

import (
	"bufio"
	"fmt"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"io"
)

// TODO: Repl struct with parameters for the header text, input prefix, ...

// Start initiates a new REPL session taking input form in and outputting
// it to out. The parameter args modifies the behavior of the REPL.
func Start(in io.Reader, out io.Writer, args Args) {
	s := bufio.NewScanner(in)
	r := read.NewReader()
	p := read.NewParser()
	w := print.NewPrinter()

	for {
		fmt.Fprintf(out, ">> ")
		if ok := s.Scan(); !ok {
			return
		}
		input := s.Text()
		// if input == ":x" {
		// 	fmt.Fprintf(out, "Bye.\n")
		// 	return
		// }
		r.Load(input)
		// for t := r.Next(); t != ""; t = r.Next() {
		// 	fmt.Fprintf(out, "%s\n", t)
		// }
		ast := p.Parse(r)
		output := w.Print(ast)
		fmt.Fprintf(out, "%s\n", output)
	}
}
