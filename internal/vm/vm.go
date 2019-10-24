package vm

import (
	"encoding/binary"
)

type VM interface {
	Run(code Ins)
	Inspect(offset int64) Val
	StackSize() int64
}

type vm struct {
	ip    int64
	sp    int64
	ep    int64
	stack []Val
	envs  []Env
	code  Ins
}

func New(size uint32) VM {
	return &vm{
		ip:    0,
		sp:    0,
		ep:    0,
		stack: make([]Val, size),
		envs:  make([]Env, 1),
	}
}

func (m *vm) Inspect(offset int64) Val {
	return m.stack[m.sp-offset-1]
}

func (m *vm) StackSize() int64 {
	return m.sp
}

func (m *vm) Run(code Ins) {
	m.ip = 0
	m.code = code
	len := int64(len(code))
	for m.ip < len {
		switch m.readOp() {
		case OpConst:
			m.push(m.readUint64())
		case OpPop:
			m.pop()
		case OpAdd:
			r := m.pop().(uint64)
			l := m.pop().(uint64)
			m.push(l + r)
		case OpSub:
			r := m.pop().(uint64)
			l := m.pop().(uint64)
			m.push(l - r)
		case OpMul:
			r := m.pop().(uint64)
			l := m.pop().(uint64)
			m.push(l * r)
		case OpDiv:
			r := m.pop().(uint64)
			l := m.pop().(uint64)
			m.push(l / r)
		case OpFalse:
			m.push(false)
		case OpTrue:
			m.push(true)
		case OpJumpIfFalse:
			c := m.pop().(bool)
			t := m.readInt64()
			if !c {
				m.ip += t
			}
		default:
			panic("Unsupported operation.")
		}
	}
}

func (m *vm) push(v Val) {
	m.stack[m.sp] = v
	m.sp++
}

// func (m *vm) peek() Val {
// 	if m.sp <= 0 {
// 		return nil
// 	}
// 	return m.stack[m.sp-1]
// }

func (m *vm) pop() Val {
	v := m.stack[m.sp-1]
	m.sp--
	return v
}

// func (m *vm) peekOp() Op {
// 	return Op(m.code[m.ip])
// }

func (m *vm) readOp() Op {
	op := Op(m.code[m.ip])
	m.ip++
	return op
}

func (m *vm) readUint64() uint64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return v
}

func (m *vm) readInt64() int64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return int64(v)
}
