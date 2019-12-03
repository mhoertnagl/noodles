package print

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/mhoertnagl/splis2/internal/data"
)

type Printer interface {
	Fprint(out io.Writer, node data.Node)
	Fprintln(out io.Writer, node data.Node)
	Print(node data.Node) string
	FprintErrors(err io.Writer, errs []*data.ErrorNode)
	PrintErrors(errs ...*data.ErrorNode) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Fprint(out io.Writer, node data.Node) {
	res := p.Print(node)
	if len(res) > 0 {
		fmt.Fprint(out, res)
	}
}

func (p *printer) Fprintln(out io.Writer, node data.Node) {
	res := p.Print(node)
	if len(res) > 0 {
		fmt.Fprintln(out, res)
	}
}

func (p *printer) Print(node data.Node) string {
	p.buf.Reset()
	p.print(node)
	return p.buf.String()
}

func (p *printer) FprintErrors(err io.Writer, errs []*data.ErrorNode) {
	if len(errs) > 0 {
		fmt.Fprintln(err, p.PrintErrors(errs...))
	}
}

func (p *printer) PrintErrors(errs ...*data.ErrorNode) string {
	var errBuf bytes.Buffer
	for _, e := range errs {
		s := fmt.Sprintf("ERROR: %s\n", e.Msg)
		errBuf.WriteString(s)
	}
	return errBuf.String()
}

func (p *printer) print(n data.Node) {
	switch x := n.(type) {
	case *data.ErrorNode:
		p.buf.WriteString("  [ERROR]  ")
	case nil:
		p.buf.WriteString("nil")
	case bool:
		p.buf.WriteString(strconv.FormatBool(x))
	case float64:
		p.buf.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	case string:
		p.printString(x)
		//p.buf.WriteString(x)
	case *data.SymbolNode:
		p.buf.WriteString(x.Name)
	case *data.ListNode:
		p.printSeq(x.Items, "(", ")")
	case *data.VectorNode:
		p.printSeq(x.Items, "[", "]")
	case *data.HashMapNode:
		p.printHashMap(x.Items)
		// case data.IsFuncNode(n):
		// 	p.buf.WriteString(n.(*data.FuncNode).Name)
	}
}

func (p *printer) printString(s string) {
	p.buf.WriteString(`"`)
	p.buf.WriteString(s)
	p.buf.WriteString(`"`)
}

func (p *printer) printSeq(items []data.Node, start string, end string) {
	p.buf.WriteString(start)
	for i, item := range items {
		if i > 0 {
			p.buf.WriteString(" ")
		}
		p.print(item)
	}
	p.buf.WriteString(end)
}

func (p *printer) printHashMap(m data.Map) {
	p.buf.WriteString("{")
	// TODO: Unfortunate.
	init := false
	for key, val := range m {
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

// func (p *printer) PrintEnv(e data.Env) string {
// 	var buf bytes.Buffer
// 	w := NewPrinter()
// 	buf.WriteString("- DEFS ---------------------------------\n")
// 	for k, v := range e.defs {
// 		buf.WriteString("  ")
// 		buf.WriteString(k)
// 		buf.WriteString(" = ")
// 		buf.WriteString(w.Print(v))
// 		buf.WriteString("\n")
// 	}
// 	buf.WriteString("- SPECIALS -----------------------------\n")
// 	for k := range e.specials {
// 		buf.WriteString("  ")
// 		buf.WriteString(k)
// 		buf.WriteString("\n")
// 	}
// 	return buf.String()
// }
