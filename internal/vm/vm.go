package vm

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
)

// This is a special marker that marks the end of a sequence on the stack.
// For instance an end is pushed onto the stack before the arguments to a
// function invocation are pushed. The virtual machine can leverage this marker
// to provide varargs support.
var end Val = nil

type VM struct {
	ip     int64
	sp     int64
	fp     int64
	fsp    int64
	defs   []Val
	stack  []Val
	frames []Val
	code   Ins
}

func NewVM(stackSize int64, envStackSize int64, frameStackSize int64) *VM {
	return &VM{
		ip:     0,
		sp:     0,
		fp:     0,
		fsp:    0,
		defs:   make([]Val, envStackSize),
		stack:  make([]Val, stackSize),
		frames: make([]Val, frameStackSize),
	}
}

func (m *VM) Run(code Ins) {
	m.code = code
	ln := int64(len(code))
	for m.ip = 0; m.ip < ln; {
		switch m.readOp() {
		// case OpNil:
		//   m.push(nil)
		case OpConst:
			c := m.readInt64()
			// fmt.Printf("Const %d\n", c)
			m.push(c)
		case OpConstF:
			c := m.readFloat64()
			// fmt.Printf("Const %d\n", c)
			m.push(c)
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
			// fmt.Printf("Pop\n")
		case OpAdd:
			var s Val = int64(0)
			for v := m.pop(); v != end; v = m.pop() {
				s = m.add(s, v)
			}
			m.push(s)
		case OpSub:
			r := m.pop()
			l := m.pop()
			m.push(m.sub(l, r))
		case OpMul:
			var s Val = int64(1)
			for v := m.pop(); v != end; v = m.pop() {
				s = m.mul(s, v)
			}
			m.push(s)
		case OpDiv:
			r := m.pop()
			l := m.pop()
			m.push(m.div(l, r))
			// fmt.Printf("Div\n")
		case OpMod:
			r := m.popInt64()
			l := m.popInt64()
			m.push(l % r)
		case OpNot:
			v := m.popBool()
			m.push(!v)
		case OpList:
			l := make([]Val, 0)
			for v := m.pop(); v != end; v = m.pop() {
				l = append(l, v)
			}
			m.push(l)
		case OpCons:
			v := m.pop()
			l := m.popVector()
			// TODO: This will not create a copy of the vector.
			m.push(prepend(v, l))
		case OpAppend:
			v := m.pop()
			l := m.popVector()
			// TODO: This will not create a copy of the vector.
			m.push(append(l, v))
		case OpConcat:
			l := make([]Val, 0)
			for v := m.pop(); v != end; v = m.pop() {
				if vl, ok := v.([]Val); ok {
					l = append(l, vl...)
				}
			}
			m.push(l)
		case OpNth:
			l := m.popVector()
			n := m.popInt64()
			// TODO: panic if null?
			m.push(l[n])
		case OpDrop:
			l := m.popVector()
			n := m.popInt64()
			if len(l) == 0 {
				// TODO: push fresh empty vector?
				m.push(l)
			} else {
				m.push(l[n:])
			}
		case OpLength:
			l := m.popVector()
			m.push(int64(len(l)))
		case OpDissolve:
			l := m.popVector()
			for i := len(l) - 1; i >= 0; i-- {
				m.push(l[i])
			}
		case OpJoin:
			// str := ""
			var sb strings.Builder
			for v := m.pop(); v != end; v = m.pop() {
				if s, ok := v.(string); ok {
					sb.WriteString(s)
					// str = s + str
				}
			}
			// m.push(str)
			m.push(sb.String())
		case OpExplode:
			s := m.popStr()
			l := make([]Val, 0)
			for _, c := range s {
				l = append(l, string(c))
			}
			m.push(l)
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
			m.ip = m.readInt64()
			// fmt.Printf("Jump\n")
		case OpJumpIf:
			ip := m.readInt64()
			if m.popBool() {
				m.ip = ip
			}
		case OpJumpIfNot:
			ip := m.readInt64()
			if m.popBool() == false {
				m.ip = ip
			}

		case OpPushArgs:
			// Pops n arguments from the stack and pushes them onto the frames stack.
			// This will reverse the order of the arguments.
			n := m.readInt64()
			// m.printFrames()
			// fmt.Printf("PushArgs %d\n", n)
			for ; n > 0; n-- {
				m.pushFrame(m.pop())
			}
			// m.printFrames()
		case OpDropArgs:
			m.fsp -= m.readInt64()
		case OpGetArg:
			// The first argument is at FRAMES[FP], the nth at FRAMES[FP + n].
			d := m.readInt64()
			a := m.fp + d
			// fmt.Printf("GetArg @[%d + %d] = %v\n", m.fp, d, m.frames[a])
			// m.printFrames()
			m.push(m.frames[a])
		case OpSetGlobal:
			m.defs[m.readInt64()] = m.pop()
			// fmt.Printf("SetGlobal\n")
		case OpGetGlobal:
			m.push(m.defs[m.readInt64()])
			// fmt.Printf("GetGlobal\n")
		case OpRef:
			n := m.readInt64()
			a := m.readInt64()
			r := NewRef(a)
			// Pop n closure arguments from the stack. We will save them in the Ref
			// struct and push them on the stack whenever the closure gets called.
			for ; n > 0; n-- {
				r.Add(m.pop())
			}
			// fmt.Printf("Ref %v @%v\n", r.cargs, r.addr)
			m.push(r)
		case OpCall:
			m.pushFrame(m.ip) // Push IP.
			m.pushFrame(m.fp) // Push pointer to previous frame.
			r := m.popRef()
			// Push closue arguments.
			for _, carg := range r.cargs {
				m.push(carg)
			}
			// fmt.Printf("Call @%d\n", r.addr)
			// m.printStack()
			// m.printFrames()
			m.fp = m.fsp  // Set pointer to new frame.
			m.ip = r.addr // Call function.
		case OpRecCall:
			r := m.popRef()
			// Push closue arguments.
			for _, carg := range r.cargs {
				m.push(carg)
			}
			m.fsp = m.fp  // Drop arguments.
			m.ip = r.addr // Call function.
		case OpReturn:
			m.fsp = m.fp             // Drop arguments.
			m.fp = m.popFrameInt64() // Restore pointer to previous frame.
			m.ip = m.popFrameInt64() // Restore IP.
			// fmt.Printf("Return\n")
		case OpEnd:
			m.push(end)
			// fmt.Printf("End\n")
		case OpHalt:
			return
		case OpWrite:
			f := m.popFileDesc()
			for v := m.pop(); v != end; v = m.pop() {
				fmt.Fprint(f, v)
			}
		case OpDebug:
			mode := m.readUint64()
			// Bit 0 show stack.
			if mode&DbgStack > 0 {
				m.printStack()
			}
			// Bit 1 show frames stack.
			if mode&DbgFrames > 0 {
				m.printFrames()
			}
			fmt.Print("\n")
		default:
			panic("Unsupported operation.")
		}
		// m.printStack()
		// m.printFrames()
		// fmt.Printf("---\n")
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

func (m *VM) popStr() string {
	return m.pop().(string)
}

func (m *VM) popVector() []Val {
	return m.pop().([]Val)
}

func (m *VM) popFileDesc() *os.File {
	return m.pop().(*os.File)
}

func (m *VM) popRef() *Ref {
	return m.pop().(*Ref)
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

func (m *VM) readFloat64() float64 {
	v := binary.BigEndian.Uint64(m.code[m.ip : m.ip+8])
	m.ip += 8
	return math.Float64frombits(v)
}

func (m *VM) readString(l int64) string {
	s := string(m.code[m.ip : m.ip+l])
	m.ip += l
	return s
}

func (m *VM) add(l Val, r Val) Val {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll + rr
		case float64:
			return float64(ll) + rr
		default:
			panic(fmt.Sprintf("Cannot add %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll + float64(rr)
		case float64:
			return ll + rr
		default:
			panic(fmt.Sprintf("Cannot add %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot add %v", ll))
	}
}

func (m *VM) sub(l Val, r Val) Val {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll - rr
		case float64:
			return float64(ll) - rr
		default:
			panic(fmt.Sprintf("Cannot subtract %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll - float64(rr)
		case float64:
			return ll - rr
		default:
			panic(fmt.Sprintf("Cannot subtract %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot subtract %v", ll))
	}
}

func (m *VM) mul(l Val, r Val) Val {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll * rr
		case float64:
			return float64(ll) * rr
		default:
			panic(fmt.Sprintf("Cannot multiply %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll * float64(rr)
		case float64:
			return ll * rr
		default:
			panic(fmt.Sprintf("Cannot multiply %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot multiply %v", ll))
	}
}

func (m *VM) div(l Val, r Val) Val {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return float64(ll) / float64(rr)
		case float64:
			return float64(ll) / rr
		default:
			panic(fmt.Sprintf("Cannot divide %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll / float64(rr)
		case float64:
			return ll / rr
		default:
			panic(fmt.Sprintf("Cannot divide %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot divide %v", ll))
	}
}

func (m *VM) eq(l Val, r Val) bool {
	switch ll := l.(type) {
	case bool:
		switch rr := r.(type) {
		case bool:
			return ll == rr
		}
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll == rr
		case float64:
			return float64(ll) == rr
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll == float64(rr)
		case float64:
			return ll == rr
		}
	case string:
		switch rr := r.(type) {
		case string:
			return ll == rr
		}
	case []Val:
		switch rr := r.(type) {
		case []Val:
			return m.eqSeq(ll, rr)
		}
	}
	return false
}

func (m *VM) eqSeq(l []Val, r []Val) bool {
	if len(l) != len(r) {
		return false
	}
	for i := 0; i < len(l); i++ {
		if m.eq(l[i], r[i]) == false {
			return false
		}
	}
	return true
}

func (m *VM) lt(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll < rr
		case float64:
			return float64(ll) < rr
		default:
			panic(fmt.Sprintf("Cannot < %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll < float64(rr)
		case float64:
			return ll < rr
		default:
			panic(fmt.Sprintf("Cannot < %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot < %v", ll))
	}
}

func (m *VM) le(l Val, r Val) bool {
	switch ll := l.(type) {
	case int64:
		switch rr := r.(type) {
		case int64:
			return ll <= rr
		case float64:
			return float64(ll) <= rr
		default:
			panic(fmt.Sprintf("Cannot <= %v", rr))
		}
	case float64:
		switch rr := r.(type) {
		case int64:
			return ll <= float64(rr)
		case float64:
			return ll <= rr
		default:
			panic(fmt.Sprintf("Cannot <= %v", rr))
		}
	default:
		panic(fmt.Sprintf("Cannot <= %v", ll))
	}
}

func prepend(v Val, l []Val) []Val {
	return append([]Val{v}, l...)
}
