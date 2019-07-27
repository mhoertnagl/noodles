package repl

import (
	"bufio"
	"fmt"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, ">> ")
		if ok := s.Scan(); !ok {
			return
		}
		input := s.Text()
		if input == ":exit" {
			fmt.Fprintf(out, "Bye.\n")
			return
		}
		fmt.Fprintf(out, input)
	}
}
