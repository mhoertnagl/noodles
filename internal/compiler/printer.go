package compiler

import (
	"bytes"
	"fmt"
	"strconv"
)

type Printer interface {
	Print(node Node) string
	PrintErrors(errs ...*ErrorNode) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(node Node) string {
	p.buf.Reset()
	p.print(node)
	return p.buf.String()
}

func (p *printer) print(n Node) {
	switch {
	case IsError(n):
		p.buf.WriteString("  [ERROR]  ")
	case IsNil(n):
		p.buf.WriteString("nil")
	case IsBool(n):
		p.buf.WriteString(strconv.FormatBool(n.(bool)))
	case IsInteger(n):
		p.buf.WriteString(strconv.FormatInt(n.(int64), 10))
	case IsNumber(n):
		p.buf.WriteString(strconv.FormatFloat(n.(float64), 'f', -1, 64))
	case IsString(n):
		p.printString(n.(string))
	case IsSymbol(n):
		p.buf.WriteString(n.(*SymbolNode).Name)
	case IsList(n):
		p.printSeq(n.(*ListNode).Items, "(", ")")
	case IsVector(n):
		p.printSeq(n.(*VectorNode).Items, "[", "]")
	case IsHashMap(n):
		p.printHashMap(n.(*HashMapNode).Items)
		// case IsFuncNode(n):
		// 	p.buf.WriteString(n.(*FuncNode).Name)
	}
}

func (p *printer) printString(s string) {
	p.buf.WriteString(`"`)
	p.buf.WriteString(s)
	p.buf.WriteString(`"`)
}

func (p *printer) printSeq(items []Node, start string, end string) {
	p.buf.WriteString(start)
	for i, item := range items {
		if i > 0 {
			p.buf.WriteString(" ")
		}
		p.print(item)
	}
	p.buf.WriteString(end)
}

func (p *printer) printHashMap(items Map) {
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

func (p *printer) PrintErrors(errs ...*ErrorNode) string {
	var errBuf bytes.Buffer
	for _, e := range errs {
		s := fmt.Sprintf("ERROR: %s\n", e.Msg)
		errBuf.WriteString(s)
	}
	return errBuf.String()
}

// func (p *printer) PrintEnv(e Env) string {
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
