package print

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/read"
)

type Printer interface {
	Print(node data.Node) string
	PrintEnv(node data.Node, env data.Env) string
	PrintErrors(parser read.Parser) string
	PrintError(n *data.ErrorNode) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(node data.Node) string {
	p.buf.Reset()
	// p.print(node)
	p.print(node, nil)
	return p.buf.String()
}

// func (p *printer) print(n data.Node) {
// 	switch {
// 	case data.IsError(n):
// 		p.buf.WriteString("  [ERROR]  ")
// 	case data.IsNil(n):
// 		p.buf.WriteString("nil")
// 	case data.IsBool(n):
// 		p.buf.WriteString(strconv.FormatBool(n.(bool)))
// 	case data.IsNumber(n):
// 		p.buf.WriteString(strconv.FormatFloat(n.(float64), 'g', -1, 64))
// 	case data.IsString(n):
// 		p.printString(n.(string))
// 	case data.IsSymbol(n):
// 		p.buf.WriteString(n.(*data.SymbolNode).Name)
// 	case data.IsList(n):
// 		p.printSeq(n.(*data.ListNode).Items, "(", ")")
// 	case data.IsVector(n):
// 		p.printSeq(n.(*data.VectorNode).Items, "[", "]")
// 	case data.IsHashMap(n):
// 		p.printHashMap(n.(*data.HashMapNode).Items)
// 		// case eval.IsFuncNode(n):
// 		// 	p.buf.WriteString("#<fn>")
// 	}
// }

func (p *printer) PrintEnv(node data.Node, env data.Env) string {
	p.buf.Reset()
	p.print(node, env)
	return p.buf.String()
}

func (p *printer) print(n data.Node, env data.Env) {
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
		// sym := n.(*data.SymbolNode).Name
		// if env == nil {
		// 	p.buf.WriteString(sym)
		// } else if val, ok := env.Lookup(sym); ok {
		// 	p.print(val, env)
		// } else {
		// 	p.buf.WriteString(sym)
		// }
	case data.IsList(n):
		p.printSeq(n.(*data.ListNode).Items, env, "(", ")")
	case data.IsVector(n):
		p.printSeq(n.(*data.VectorNode).Items, env, "[", "]")
	case data.IsHashMap(n):
		p.printHashMap(n.(*data.HashMapNode).Items, env)
		// case data.IsFuncNode(n):
		// 	p.buf.WriteString(n.(*data.FuncNode).Name)
	}
}

func (p *printer) printString(s string) {
	p.buf.WriteString(`"`)
	p.buf.WriteString(s)
	p.buf.WriteString(`"`)
}

func (p *printer) printSeq(items []data.Node, env data.Env, start string, end string) {
	p.buf.WriteString(start)
	for i, item := range items {
		if i > 0 {
			p.buf.WriteString(" ")
		}
		p.print(item, env)
	}
	p.buf.WriteString(end)
}

func (p *printer) printHashMap(items data.Map, env data.Env) {
	p.buf.WriteString("{")
	// TODO: Unfortunate.
	init := false
	for key, val := range items {
		if init {
			p.buf.WriteString(" ")
		}
		init = true
		p.print(key, env)
		p.buf.WriteString(" ")
		p.print(val, env)
	}
	p.buf.WriteString("}")
}

func (p *printer) PrintErrors(parser read.Parser) string {
	var errBuf bytes.Buffer
	for _, e := range parser.Errors() {
		errBuf.WriteString(p.PrintError(e))
	}
	return errBuf.String()
}

func (p *printer) PrintError(n *data.ErrorNode) string {
	return fmt.Sprintf("ERROR: %s\n", n.Msg)
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
