package print

import (
	"bytes"
	"github.com/mhoertnagl/splis2/internal/read"
	"strconv"
)

type Printer interface {
	Print(n read.Node) string
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
		p.buf.WriteString("ERROR\n")
	case *read.ListNode:
		p.printSeq(n.Items, "(", ")")
	case *read.VectorNode:
		p.printSeq(n.Items, "[", "]")
	case *read.HashMapNode:
		p.printSeq(n.Items, "{", "}")
	case *read.StringNode:
		p.buf.WriteString(n.Val)
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
