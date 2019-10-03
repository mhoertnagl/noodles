package print

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mhoertnagl/splis2/internal/data"
)

type Printer interface {
	Print(node data.Node) string
	PrintErrors(errs ...*data.ErrorNode) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(node data.Node) string {
	p.buf.Reset()
	p.print(node)
	return p.buf.String()
}

func (p *printer) print(n data.Node) {
	switch {
	case data.IsError(n):
		p.buf.WriteString("  [ERROR]  ")
	case data.IsNil(n):
		p.buf.WriteString("nil")
	case data.IsBool(n):
		p.buf.WriteString(strconv.FormatBool(n.(bool)))
	case data.IsNumber(n):
		p.buf.WriteString(strconv.FormatFloat(n.(float64), 'f', -1, 64))
	case data.IsString(n):
		p.printString(n.(string))
	case data.IsSymbol(n):
		p.buf.WriteString(n.(*data.SymbolNode).Name)
	case data.IsList(n):
		p.printSeq(n.(*data.ListNode).Items, "(", ")")
	case data.IsVector(n):
		p.printSeq(n.(*data.VectorNode).Items, "[", "]")
	case data.IsHashMap(n):
		p.printHashMap(n.(*data.HashMapNode).Items)
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

func (p *printer) printHashMap(items data.Map) {
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

func (p *printer) PrintErrors(errs ...*data.ErrorNode) string {
	var errBuf bytes.Buffer
	for _, e := range errs {
		s := fmt.Sprintf("ERROR: %s\n", e.Msg)
		errBuf.WriteString(s)
	}
	return errBuf.String()
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
