package asm

import (
	"bytes"
	"fmt"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type AsmPrinter struct {
	code  AsmCode
	lines []string
}

func NewAsmPrinter() *AsmPrinter {
	return &AsmPrinter{}
}

func (m *AsmPrinter) Print(code AsmCode) []string {
	m.lines = make([]string, 0)
	for _, line := range code {
		switch x := line.(type) {
		case *AsmLabel:
			m.write("%s:", x.Name)
		case *AsmLabeled:
			m.write("  %s %s", m.opName(x.Op), x.Name)
		case *AsmIns:
			m.writeInstr(x)
		case *AsmStr:
			m.write("  %s '%s'", m.opName(vm.OpStr), x.Str)
		}
	}
	return m.lines
}

func (m *AsmPrinter) PrintToStr(code AsmCode) string {
	var buf bytes.Buffer
	for _, line := range m.Print(code) {
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m *AsmPrinter) opMeta(op vm.Op) (*vm.OpMeta, bool) {
	meta, err := vm.LookupMeta(op)
	return meta, err == nil
}

func (m *AsmPrinter) opName(op vm.Op) string {
	if meta, ok := m.opMeta(op); ok {
		return meta.Name
	}
	return fmt.Sprintf("Invalid [%d]", op)
}

func (m *AsmPrinter) write(format string, a ...interface{}) {
	m.lines = append(m.lines, fmt.Sprintf(format, a...))
}

func (m *AsmPrinter) writeInstr(ins *AsmIns) {
	var buf bytes.Buffer
	buf.WriteString("  ")
	buf.WriteString(m.opName(ins.Op))
	for _, arg := range ins.Args {
		buf.WriteString(" ")
		buf.WriteString(fmt.Sprintf("%d", arg))
	}
	m.lines = append(m.lines, buf.String())
}
