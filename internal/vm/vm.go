package vm

import (
	"encoding/binary"
	"fmt"
	//"fmt"
)

// TODO: Turn an identifier into a 64bit hash (https://golang.org/pkg/hash/fnv/)
//       Maintain a stack of environments (type Env map[int64]Val)
//       Lookup can then be implemented the usual way.
// TODO: https://yourbasic.org/golang/bitwise-operator-cheat-sheet/

type VM interface {
	Run(code Ins)
	InspectStack(offset int64) Val
	// InspectLocals(offset int64) Val
	InspectEnvs(offset int64) Env
	StackSize() int64
}

type vm struct {
	ip     int64
	sp     int64
	ep     int64
	fp     int64
	stack  []Val
	envs   []Env
	frames []int64
	code   Ins
}

func New(
	stackSize int64,
	localsSize int64,
	envStackSize int64,
	frameStackSize int64) VM {
	m := &vm{
		ip:     0,
		sp:     0,
		ep:     0,
		fp:     0,
		stack:  make([]Val, stackSize),
		envs:   make([]Env, envStackSize),
		frames: make([]int64, frameStackSize),
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

func (m *vm) StackSize() int64 {
	return m.sp
}

func (m *vm) InspectEnvs(offset int64) Env {
	return m.envs[offset]
}

func (m *vm) InspectFrames(offset int64) int64 {
	return m.frames[offset]
}

func (m *vm) Run(code Ins) {
	m.ip = 0
	m.code = code
	len := int64(len(code))
	for m.ip < len {
		switch m.readOp() {
		case OpConst:
			m.push(m.readInt64())
		case OpFalse:
			m.push(false)
		case OpTrue:
			m.push(true)
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
		case OpSll:
			r := m.popUInt64()
			l := m.popUInt64()
			m.push(l << r)
		case OpSrl:
			r := m.popUInt64()
			l := m.popUInt64()
			m.push(l >> r)
		case OpSra:
			r := m.popUInt64()
			l := m.popInt64()
			m.push(l >> r)
		case OpJump:
			m.ip += m.readInt64()
		case OpJumpIf:
			d := m.readInt64()
			if m.popBool() {
				m.ip += d
			}
		case OpJumpIfNot:
			d := m.readInt64()
			if m.popBool() == false {
				m.ip += d
			}
		case OpNewEnv:
			m.newEnv()
		case OpPopEnv:
			m.ep--
		case OpSetLocal:
			m.bindLocal(m.readInt64(), m.pop())
		case OpGetLocal:
			m.push(m.lookupLocal(m.readInt64()))
		case OpSetGlobal:
			m.bindGlobal(m.readInt64(), m.pop())
		case OpGetGlobal:
			m.push(m.lookupGlobal(m.readInt64()))
		case OpCall:
			m.frames[m.fp] = m.ip + 1
			m.fp++
			m.ip += m.readInt64()
		case OpReturn:
			m.fp--
			m.ip = m.frames[m.fp]
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

func (m *vm) popUInt64() uint64 {
	return m.pop().(uint64)
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

func (m *vm) bindLocal(a int64, v Val) {
	env := m.envs[m.ep-1]
	if x, ok := env[a]; ok {
		panic(fmt.Sprintf("Local symbol [%d] already bound to [%v]", a, x))
	}
	env[a] = v
}

func (m *vm) lookupLocal(a int64) Val {
	for i := m.ep - 1; i >= 0; i-- {
		if v, ok := m.envs[i][a]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("Unbound local symbol [%d]", a))
}

func (m *vm) bindGlobal(a int64, v Val) {
	env := m.envs[0]
	if x, ok := env[a]; ok {
		panic(fmt.Sprintf("Global symbol [%d] already bound to [%v]", a, x))
	}
	env[a] = v
}

func (m *vm) lookupGlobal(a int64) Val {
	if v, ok := m.envs[0][a]; ok {
		return v
	}
	panic(fmt.Sprintf("Unbound global symbol [%d]", a))
}
