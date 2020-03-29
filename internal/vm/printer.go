package vm

import (
	"bytes"
	// "fmt"
	"strconv"
)

type Printer interface {
	Print(val Val) string
	// PrintErrors(errs ...*ErrorNode) string
}

type printer struct {
	buf bytes.Buffer
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(val Val) string {
	p.buf.Reset()
	p.print(val)
	return p.buf.String()
}

func (p *printer) print(val Val) {
	switch n := val.(type) {
	case bool:
		p.buf.WriteString(strconv.FormatBool(n))
	case int64:
		p.buf.WriteString(strconv.FormatInt(n, 10))
	// case float64:
	// 	p.buf.WriteString(strconv.FormatFloat(n, 'f', -1, 64))
	case string:
		p.buf.WriteString(`"`)
		p.buf.WriteString(n)
		p.buf.WriteString(`"`)
	case []Val:
		p.buf.WriteString("[")
		for i, v := range n {
			if i > 0 {
				p.buf.WriteString(" ")
			}
			p.print(v)
		}
		p.buf.WriteString("]")
	}
}

// func (p *printer) printSeq(vals []Val, start string, end string) {
// 	p.buf.WriteString(start)
// 	for i, val := range vals {
// 		if i > 0 {
// 			p.buf.WriteString(" ")
// 		}
// 		p.print(val)
// 	}
// 	p.buf.WriteString(end)
// }

// func (p *printer) printHashMap(items Map) {
// 	p.buf.WriteString("{")
// 	// TODO: Unfortunate.
// 	init := false
// 	for key, val := range items {
// 		if init {
// 			p.buf.WriteString(" ")
// 		}
// 		init = true
// 		p.print(key)
// 		p.buf.WriteString(" ")
// 		p.print(val)
// 	}
// 	p.buf.WriteString("}")
// }

// func (p *printer) PrintErrors(errs ...*ErrorNode) string {
// 	var errBuf bytes.Buffer
// 	for _, e := range errs {
// 		s := fmt.Sprintf("ERROR: %s\n", e.Msg)
// 		errBuf.WriteString(s)
// 	}
// 	return errBuf.String()
// }
