package vm

import (
	"encoding/binary"
	"fmt"
)

// TODO: *STDOUT*
// TODO: write
// TODO: str -> use printer to turn value into a string.
// TODO: *STDIN*
// TODO: read
// TODO: :::
// TODO: quot
// TODO: mod
// TODO: join (strings)

// TODO: https://yourbasic.org/golang/bitwise-operator-cheat-sheet/

// This is a special marker that marks the end of a sequence on the stack.
// For instance an end is pushed onto the stack before the arguments to a
// function invocation are pushed. The virtual machine can leverage this marker
// to provide varargs support.
var end Val = nil

type VM struct {
	ip     int64
	sp     int64
	ep     int64
	fp     int64
	fsp    int64
	stack  []Val
	envs   []Env
	frames []Val
	code   Ins
}

func New(
	stackSize int64,
	envStackSize int64,
	frameStackSize int64) *VM {
	m := &VM{
		ip:    0,
		sp:    0,
		ep:    0,
		fp:    0,
		fsp:   0,
		stack: make([]Val, stackSize),
		// TODO: We only need a single global environment in the future.
		envs:   make([]Env, envStackSize),
		frames: make([]Val, frameStackSize),
	}
	// Create the outermost environment.
	m.newEnv()
	return m
}

func (m *VM) InspectStack(offset int64) Val {
	a := m.sp - offset - 1
	if a >= 0 {
		return m.stack[a]
	}
	return nil
}

func (m *VM) StackSize() int64 {
	return m.sp
}

func (m *VM) InspectEnvs(offset int64) Env {
	return m.envs[offset]
}

func (m *VM) InspectFrames(offset int64) Val {
	a := m.fsp - offset - 1
	if a >= 0 {
		return m.frames[a]
	}
	return nil
}

func (m *VM) FramesSize() int64 {
	return m.fsp
}

func (m *VM) Run(code Ins) {
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
		case OpFst:
			l := m.popVector()
			// TODO: panic if null?
			m.push(l[0])
		case OpRest:
			l := m.popVector()
			if len(l) == 0 {
				// TODO: push fresh empty vector?
				m.push(l)
				// m.push(make([]Val, 0))
			} else {
				m.push(l[1:])
			}
		case OpLength:
			l := m.popVector()
			m.push(int64(len(l)))
		case OpDissolve:
			l := m.popVector()
			for i := len(l) - 1; i >= 0; i-- {
				m.push(l[i])
			}
		// case OpAnd:
		// 	a := ^int64(0)
		// 	for v := m.pop(); v != end; v = m.pop() {
		// 		a = a & v.(int64)
		// 	}
		// 	m.push(a)
		// case OpOr:
		// 	a := int64(0)
		// 	for v := m.pop(); v != end; v = m.pop() {
		// 		a = a | v.(int64)
		// 	}
		// 	m.push(a)
		// case OpInv:
		// 	m.push(^m.popInt64())
		// case OpSll:
		// 	r := m.popInt64()
		// 	l := m.popInt64()
		// 	m.push(l << uint64(r))
		// case OpSrl:
		// 	r := m.popInt64()
		// 	l := m.popInt64()
		// 	m.push(int64(uint64(l) >> uint64(r)))
		// case OpSra:
		// 	r := m.popInt64()
		// 	l := m.popInt64()
		// 	m.push(l >> uint64(r))
		case OpEQ:
			r := m.pop()
			l := m.pop()
			m.push(m.eq(l, r))
		case OpNE:
			r := m.pop()
			l := m.pop()
			m.push(!m.eq(l, r))
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

			// TODO: Deprecated.
		case OpNewEnv:
			m.newEnv()
			// TODO: Deprecated.
		case OpPopEnv:
			m.ep--
			// TODO: Deprecated.
		case OpSetLocal:
			m.bindEnv(m.ep-1, m.readInt64(), m.pop())
			// TODO: Deprecated.
		case OpGetLocal:
			m.push(m.lookupEnv(m.ep-1, m.readInt64()))
			// ----------------

		case OpPushArgs:
			// Pops n arguments from the stack and pushes them onto the frames stack.
			// This will reverse the order of the arguments.
			n := m.readInt64()
			for ; n > 0; n-- {
				m.pushFrame(m.pop())
			}
		case OpGetArg:
			// The first argument is at FRAMES[FP], the nth at FRAMES[FP + n].
			a := m.fp + m.readInt64()
			m.push(m.frames[a])
		case OpSetGlobal:
			m.bindEnv(0, m.readInt64(), m.pop())
		case OpGetGlobal:
			m.push(m.lookupEnv(0, m.readInt64()))
		case OpCall:
			m.pushFrame(m.ip)   // Push IP.
			m.pushFrame(m.fp)   // Push pointer to previous frame.
			m.fp = m.fsp        // Set pointer to new frame.
			m.ip = m.popInt64() // Call function.
		case OpReturn:
			m.fsp = m.fp             // Drop arguments.
			m.fp = m.popFrameInt64() // Restore pointer to previous frame.
			m.ip = m.popFrameInt64() // Restore IP.
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

func (m *VM) push(v Val) {
	m.stack[m.sp] = v
	m.sp++
}

// func (m *vm) peek() Val {
// 	return m.stack[m.sp-1]
// }

func (m *VM) pop() Val {
	m.sp--
	return m.stack[m.sp]
}

func (m *VM) popBool() bool {
	return m.pop().(bool)
}

func (m *VM) popInt64() int64 {
	return m.pop().(int64)
}

func (m *VM) popUInt64() uint64 {
	return m.pop().(uint64)
}

func (m *VM) popVector() []Val {
	return m.pop().([]Val)
}

func (m *VM) pushFrame(v Val) {
	m.frames[m.fsp] = v
	m.fsp++
}

func (m *VM) popFrame() Val {
	m.fsp--
	return m.frames[m.fsp]
}

func (m *VM) popFrameInt64() int64 {
	return m.popFrame().(int64)
}

func (m *VM) readOp() Op {
	op := m.code[m.ip]
	m.ip++
	return op
}

func (m *VM) readUint64() uint64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return v
}

func (m *VM) readInt64() int64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return int64(v)
}

func (m *VM) readString(l int64) string {
	s := string(m.code[m.ip : m.ip+l])
	m.ip += l
	return s
}

func (m *VM) newEnv() {
	m.envs[m.ep] = make(Env)
	m.ep++
}

func (m *VM) bindEnv(ep int64, a int64, v Val) {
	env := m.envs[ep]
	if x, ok := env[a]; ok {
		panic(fmt.Sprintf("Symbol [%d] already bound to [%v]", a, x))
	}
	env[a] = v
}

func (m *VM) lookupEnv(ep int64, a int64) Val {
	for i := ep; i >= 0; i-- {
		if v, ok := m.envs[i][a]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("Unbound symbol [%d]", a))
}

func (m *VM) eq(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll == rr
		}
	}
	return false
}

func (m *VM) lt(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll < rr
		}
	}
	return false
}

func (m *VM) le(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll <= rr
		}
	}
	return false
}

func (m *VM) printStack() {
	fmt.Print("-|")
	for i := int64(0); i < m.sp; i++ {
		fmt.Printf(" %v", m.stack[i])
	}
	fmt.Print("\n")
}

func prepend(v Val, l []Val) []Val {
	return append([]Val{v}, l...)
}

// func eq(e Evaluator, env data.Env, args []data.Node) data.Node {
// 	if len(args) != 2 {
// 		return e.Error("[=] expects 2 arguments.")
// 	}
// 	return eq2(e, env, args[0], args[1])
// }

// func eq2(e Evaluator, env data.Env, a, b data.Node) data.Node {
// 	if reflect.TypeOf(a) != reflect.TypeOf(b) {
// 		return false
// 	}
// 	switch x := a.(type) {
// 	case *data.SymbolNode:
// 		y := b.(*data.SymbolNode)
// 		return x.Name == y.Name
// 	case *data.ListNode:
// 		y := b.(*data.ListNode)
// 		return eqSeq(e, env, x.Items, y.Items)
// 	case *data.VectorNode:
// 		y := b.(*data.VectorNode)
// 		return eqSeq(e, env, x.Items, y.Items)
// 	case *data.HashMapNode:
// 		y := b.(*data.HashMapNode)
// 		return eqHashMap(e, env, x.Items, y.Items)
// 	default:
// 		return a == b
// 	}
// }

// func eqSeq(e Evaluator, env data.Env, as, bs []data.Node) data.Node {
// 	if len(as) != len(bs) {
// 		return false
// 	}
// 	for i := 0; i < len(as); i++ {
// 		if eq2(e, env, as[i], bs[i]) == false {
// 			return false
// 		}
// 	}
// 	return true
// }

// func eqHashMap(e Evaluator, env data.Env, as, bs data.Map) data.Node {
// 	if len(as) != len(bs) {
// 		return false
// 	}
// 	for k, va := range as {
// 		vb, ok := bs[k]
// 		if !ok {
// 			return false
// 		}
// 		if eq2(e, env, va, vb) == false {
// 			return false
// 		}
// 	}
// 	return true
// }

// func join(e Evaluator, env data.Env, args []data.Node) data.Node {
// 	var sb strings.Builder
// 	for _, arg := range args {
// 		if s, ok := arg.(string); ok {
// 			sb.WriteString(s)
// 		}
// 	}
// 	return sb.String()
// }
