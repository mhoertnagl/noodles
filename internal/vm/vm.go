package vm

import "encoding/binary"

type VM interface {
}

type vm struct {
	ip    uint64
	sp    uint64
	stack []Val
	code  Ins
}

func NewVM(size uint32) VM {
	return &vm{
		ip:    0,
		sp:    0,
		stack: make([]Val, size),
	}
}

func (m *vm) Run(code Ins) {
	m.ip = 0
	m.code = code
	for m.ip < uint64(len(code)) {
		switch m.readOp() {
		case OpConst:
			m.push(m.readUint64())
		}
	}
}

func (m *vm) push(v Val) {
	m.stack[m.sp] = v
	m.sp++
}

func (m *vm) peek() Val {
	if m.sp <= 0 {
		return nil
	}
	return m.stack[m.sp-1]
}

func (m *vm) pop() Val {
	v := m.stack[m.sp-1]
	m.sp--
	return v
}

func (m *vm) readOp() Op {
	return Op(m.code[m.ip])
}

func (m *vm) readUint64() uint64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return v
}
