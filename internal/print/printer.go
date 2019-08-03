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

func (p *printer) print(node read.Node) {
	switch n := node.(type) {
	case *read.ErrorNode:
		p.buf.WriteString("  [ERROR]  ")
	case *read.ListNode:
		p.printSeq(n.Items, "(", ")")
	case *read.VectorNode:
		p.printSeq(n.Items, "[", "]")
	case *read.HashMapNode:
		p.printHashMap(n.Items)
	case *read.StringNode:
		p.printString(n)
		//p.buf.WriteString(n.Val)
	case *read.NumberNode:
		p.buf.WriteString(strconv.FormatFloat(n.Val, 'g', -1, 64))
	case *read.SymbolNode:
		p.buf.WriteString(n.Name)
	case *read.TrueNode:
		p.buf.WriteString("true")
	case *read.FalseNode:
		p.buf.WriteString("false")
	case *read.NilNode:
		p.buf.WriteString("nil")
	}
}

func (p *printer) printString(n *read.StringNode) {
	p.buf.WriteString(`"`)
	p.buf.WriteString(n.Val)
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
