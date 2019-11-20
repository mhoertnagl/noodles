package vm

import (
	"encoding/binary"
	"fmt"
	//"fmt"
)

// TODO: Turn an identifier into a 64bit hash (https://golang.org/pkg/hash/fnv/)
//       Maintain a stack of environments (type Env map[int64]Val)
//       Lookup can then be implemented the usual way.

type VM interface {
	Run(code Ins)
	InspectStack(offset int64) Val
	// InspectLocals(offset int64) Val
	InspectEnvs(offset int64) Env
	StackSize() int64
}

type vm struct {
	ip int64
	sp int64
	// lp    int64
	ep    int64
	stack []Val
	// locals []Val
	envs []Env
	code Ins
}

func New(stackSize int64, localsSize int64, envStackSize int64) VM {
	m := &vm{
		ip: 0,
		sp: 0,
		// lp:    0,
		ep:    0,
		stack: make([]Val, stackSize),
		// locals: make([]Val, localsSize),
		envs: make([]Env, envStackSize),
	}
	// Create the outermost environment.
	m.newEnv()
	return m
}

func (m *vm) InspectStack(offset int64) Val {
	a := m.sp - offset - 1
	if a >= 0 {
		return m.stack[a]
	}
	return nil
}

// func (m *vm) InspectLocals(offset int64) Val {
// 	return m.locals[offset]
// }

func (m *vm) InspectEnvs(offset int64) Env {
	return m.envs[offset]
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
			m.push(m.readInt64())
		case OpPop:
			m.pop()
		case OpAdd:
			r := m.popInt64()
			l := m.popInt64()
			m.push(l + r)
		case OpSub:
			r := m.popInt64()
			l := m.popInt64()
			m.push(l - r)
		case OpMul:
			r := m.popInt64()
			l := m.popInt64()
			m.push(l * r)
		case OpDiv:
			r := m.popInt64()
			l := m.popInt64()
			m.push(l / r)
		case OpFalse:
			m.push(false)
		case OpTrue:
			m.push(true)
		case OpJump:
			m.ip += m.readInt64()
		case OpJumpIfFalse:
			d := m.readInt64()
			if m.popBool() == false {
				m.ip += d
			}
		case OpJumpIfTrue:
			d := m.readInt64()
			if m.popBool() {
				m.ip += d
			}
		case OpNewEnv:
			m.newEnv()
		case OpPopEnv:
			m.ep--
		case OpSetLocal:
			m.bind(m.readInt64(), m.pop())
		case OpGetLocal:
			v := m.lookup(m.readInt64())
			m.push(v)
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
	m.sp--
	return m.stack[m.sp]
}

func (m *vm) popBool() bool {
	return m.pop().(bool)
}

func (m *vm) popInt64() int64 {
	return m.pop().(int64)
}

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

func (m *vm) newEnv() {
	m.envs[m.ep] = make(Env)
	m.ep++
}

func (m *vm) bind(a int64, v Val) {
	m.envs[m.ep-1][a] = v
}

func (m *vm) lookup(a int64) Val {
	for i := m.ep - 1; i >= 0; i-- {
		if v, ok := m.envs[i][a]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("Unbound symbol [%d]", a))
}
