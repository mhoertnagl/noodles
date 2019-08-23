package print

import (
	"bytes"
	"fmt"
	"github.com/mhoertnagl/splis2/internal/read"
	"strconv"
)

type Printer interface {
	Print(node read.Node) string
	PrintErrors(parser read.Parser) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(node read.Node) string {
	p.buf.Reset()
	p.print(node)
	return p.buf.String()
}

func (p *printer) print(n read.Node) {
	switch {
	case read.IsError(n):
		p.buf.WriteString("  [ERROR]  ")
	case read.IsNil(n):
		p.buf.WriteString("nil")
	case read.IsBool(n):
		p.buf.WriteString(strconv.FormatBool(n.(bool)))
	case read.IsNumber(n):
		p.buf.WriteString(strconv.FormatFloat(n.(float64), 'g', -1, 64))
	case read.IsString(n):
		p.printString(n.(string))
	case read.IsSymbol(n):
		p.buf.WriteString(n.(*read.SymbolNode).Name)
	case read.IsList(n):
		p.printSeq(n.(*read.ListNode).Items, "(", ")")
	case read.IsVector(n):
		p.printSeq(n.(*read.VectorNode).Items, "[", "]")
	case read.IsHashMap(n):
		p.printHashMap(n.(*read.HashMapNode).Items)
	}
}

func (p *printer) printString(s string) {
	p.buf.WriteString(`"`)
	p.buf.WriteString(s)
	p.buf.WriteString(`"`)
}

func (p *printer) printSeq(items []read.Node, start string, end string) {
	p.buf.WriteString(start)
	for i, item := range items {
		if i > 0 {
			p.buf.WriteString(" ")
		}
		p.print(item)
	}
	p.buf.WriteString(end)
}

func (p *printer) printHashMap(items map[read.Node]read.Node) {
	p.buf.WriteString("{")
	// TODO: Unfortunate.
	init := false
	for key, val := range items {
		if init {
			p.buf.WriteString(" ")
		}
		init = true
		p.print(key)
		p.buf.WriteString(" ")
		p.print(val)
	}
	p.buf.WriteString("}")
}

func (p *printer) PrintErrors(parser read.Parser) string {
	var errBuf bytes.Buffer
	for _, e := range parser.Errors() {
		errBuf.WriteString(p.printError(e))
	}
	return errBuf.String()
}

func (p *printer) printError(n *read.ErrorNode) string {
	return fmt.Sprintf("ERROR: %s\n", n.Msg)
}
