package asm

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mhoertnagl/noodles/internal/vm"
)

type Disassembler struct {
	ip    int64
	code  []byte
	lines []string
}

func NewDisassembler() *Disassembler {
	return &Disassembler{
		ip: 0,
	}
}

func (m *Disassembler) Disassemble(code []byte) []string {
	m.lines = make([]string, 0)
	m.code = code
	ln := int64(len(code))
	for m.ip = 0; m.ip < ln; {
		m.writeInstr(m.readOp())
	}
	return m.lines
}

func (m *Disassembler) DisassembleToStr(code []byte) string {
	var buf bytes.Buffer
	m.Disassemble(code)
	for _, line := range m.lines {
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m *Disassembler) readOp() vm.Op {
	op := m.code[m.ip]
	m.ip++
	return op
}

func (m *Disassembler) readArg(sz int) uint64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+int64(sz)])
	m.ip += int64(sz)
	return v
}

func (m *Disassembler) readString(l int64) string {
	s := string(m.code[m.ip : m.ip+l])
	m.ip += l
	return s
}

func (m *Disassembler) write(format string, a ...interface{}) {
	m.lines = append(m.lines, fmt.Sprintf(format, a...))
}

func (m *Disassembler) writeInstr(op vm.Op) {
	if meta, ok := vm.LookupMeta(op); ok == nil {
		switch op {
		case vm.OpStr:
			slen := int64(m.readArg(meta.Args[0]))
			str := m.readString(slen)
			m.write("%s '%s'", meta.Name, str)
		default:
			var buf bytes.Buffer
			buf.WriteString(meta.Name)
			for _, sz := range meta.Args {
				buf.WriteString(" ")
				buf.WriteString(fmt.Sprintf("%d", m.readArg(sz)))
			}
			m.write("%s", buf.String())
		}
	} else {
		m.write("Invalid [%d]", op)
	}
}
