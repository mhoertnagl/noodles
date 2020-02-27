package vm

import (
	"encoding/binary"
	"fmt"
)

// TODO: https://yourbasic.org/golang/bitwise-operator-cheat-sheet/

// This is a special marker that marks the end of a sequence on the stack.
// For instance an end is pushed onto the stack before the arguments to a
// function invocation are pushed. The virtual machine can leverage this marker
// to provide varargs support.
var end Val = nil

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
	ln := int64(len(code))
	for m.ip < ln {
		switch m.readOp() {
		case OpConst:
			m.push(m.readInt64())
			// case OpNil:
			//   m.push(nil)
		case OpRef:
			m.push(m.readInt64())
		case OpFalse:
			m.push(false)
		case OpTrue:
			m.push(true)
		case OpEmptyVector:
			m.push(make([]Val, 0))
		case OpStr:
			l := m.readUint64()
			m.push(m.readString(int64(l)))
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
		case OpNot:
			v := m.popBool()
			m.push(!v)
		case OpList:
			l := make([]Val, 0)
			for v := m.pop(); v != end; v = m.pop() {
				l = append(l, v)
			}
			m.push(end)
			m.push(l)
		case OpCons:
			v := m.pop()
			l := m.popVector()
			// TODO: This will not create a copy of the vector.
			m.push(prepend(v, l))
		case OpDissolve:
			l := m.popVector()
			for i := len(l) - 1; i >= 0; i-- {
				m.push(l[i])
			}
			// case OpHead:
			// 	l := m.popVector()
			// 	m.push(l[0])
			// case OpTail:
			// 	l := m.popVector()
			// 	m.push(l[1:])
			// case OpSll:
			// 	r := m.popUInt64()
			// 	l := m.popUInt64()
			// 	m.push(l << r)
			// case OpSrl:
			// 	r := m.popUInt64()
			// 	l := m.popUInt64()
			// 	m.push(l >> r)
			// case OpSra:
			// 	r := m.popUInt64()
			// 	l := m.popInt64()
			// 	m.push(l >> r)
		case OpEQ:
			r := m.pop()
			l := m.pop()
			m.push(m.eq(l, r))
		case OpNE:
			r := m.pop()
			l := m.pop()
			m.push(m.eq(l, r) == false)
		case OpLT:
			r := m.pop()
			l := m.pop()
			m.push(m.lt(l, r))
		case OpLE:
			r := m.pop()
			l := m.pop()
			m.push(m.le(l, r))
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
			m.bindEnv(m.ep-1, m.readInt64(), m.pop())
		case OpGetLocal:
			m.push(m.lookupEnv(m.ep-1, m.readInt64()))
		case OpSetGlobal:
			m.bindEnv(0, m.readInt64(), m.pop())
		case OpGetGlobal:
			m.push(m.lookupEnv(0, m.readInt64()))
		case OpCall:
			m.frames[m.fp] = m.ip
			m.fp++
			m.ip = m.popInt64()
		case OpReturn:
			m.fp--
			m.ip = m.frames[m.fp]
		case OpEnd:
			m.push(end)
		case OpHalt:
			return
		case OpDebug:
			mode := m.readUint64()
			// Bit 0 show stack.
			if mode&DbgStack == 1 {
				m.printStack()
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

func (m *vm) popVector() []Val {
	return m.pop().([]Val)
}

func (m *vm) readOp() Op {
	op := m.code[m.ip]
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

func (m *vm) readString(l int64) string {
	s := string(m.code[m.ip : m.ip+l])
	m.ip += l
	return s
}

func (m *vm) newEnv() {
	m.envs[m.ep] = make(Env)
	m.ep++
}

func (m *vm) bindEnv(ep int64, a int64, v Val) {
	env := m.envs[ep]
	if x, ok := env[a]; ok {
		panic(fmt.Sprintf("Symbol [%d] already bound to [%v]", a, x))
	}
	env[a] = v
}

func (m *vm) lookupEnv(ep int64, a int64) Val {
	for i := ep; i >= 0; i-- {
		if v, ok := m.envs[i][a]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("Unbound symbol [%d]", a))
}

func (m *vm) eq(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll == rr
		}
	}
	return false
}

func (m *vm) lt(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll < rr
		}
	}
	return false
}

func (m *vm) le(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll <= rr
		}
	}
	return false
}

func (m *vm) printStack() {
	fmt.Print("-|")
	for i := int64(0); i < m.sp; i++ {
		fmt.Printf(" %v", m.stack[i])
	}
	fmt.Print("\n")
}

func prepend(v Val, l []Val) []Val {
	return append([]Val{v}, l...)
}
