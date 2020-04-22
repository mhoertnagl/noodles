package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Op defines the type of the opcode as a byte. Hence there are 256 different
// operations at most.
type Op = byte

// Ins defines the type of a sequence of instructions as an array of bytes.
type Ins = []byte

const (
	OpConst Op = iota
	// OpNil
	OpTrue
	OpFalse
	OpEmptyList
	OpEmptyVector
	OpStr

	OpAdd
	OpSub
	OpMul
	OpDiv

	OpList
	OpCons
	OpNth
	OpDrop
	OpLength
	OpDissolve

	// OpAnd
	// OpOr
	// OpInv
	// OpSll
	// OpSrl
	// OpSra

	OpNot
	OpEQ
	OpNE
	OpLT
	OpLE

	OpJump
	OpJumpIf
	OpJumpIfNot

	OpPop
	// OpDup

	OpSetGlobal
	OpGetGlobal

	OpPushArgs
	OpDropArgs
	OpGetArg

	OpRef
	OpCall
	OpRecCall
	// OpTailCall
	OpReturn
	OpEnd

	OpRead
	OpWrite
	// TODO: Perhaps obsolete if we switch script and function portion and start
	//       vm at script start.
	OpHalt
	OpDebug
)

// Arguments to OpDebug.
const (
	// DbgStack prints the stack.
	DbgStack = uint64(1 << 0)
	// DbgFrames prints the frames stack.
	DbgFrames = uint64(1 << 1)
)

// OpMeta contains the human-readable name of the operation and the length in
// bytes of each of its arguments.
type OpMeta struct {
	Name string
	Args []int
}

var meta = map[Op]*OpMeta{
	OpConst: {"Const", []int{8}},
	// OpNil:         {"Nil", []int{}},
	OpTrue:        {"True", []int{}},
	OpFalse:       {"False", []int{}},
	OpEmptyList:   {"EmptyList", []int{}},
	OpEmptyVector: {"EmptyVector", []int{}},
	OpStr:         {"String", []int{8}},

	OpAdd: {"Add", []int{}},
	OpSub: {"Sub", []int{}},
	OpMul: {"Mul", []int{}},
	OpDiv: {"Div", []int{}},

	OpList:     {"List", []int{}},
	OpCons:     {"Cons", []int{}},
	OpNth:      {"Nth", []int{}},
	OpDrop:     {"Drop", []int{}},
	OpLength:   {"Tail", []int{}},
	OpDissolve: {"Dissolve", []int{}},
	// OpAnd:         {"And", []int{}},
	// OpOr:          {"Or", []int{}},
	// OpInv:         {"Inv", []int{}},
	// OpSll:         {"Sll", []int{}},
	// OpSrl:         {"Srl", []int{}},
	// OpSra:         {"Sra", []int{}},
	OpNot:       {"Not", []int{}},
	OpEQ:        {"EQ", []int{}},
	OpNE:        {"NE", []int{}},
	OpLT:        {"LT", []int{}},
	OpLE:        {"LE", []int{}},
	OpJump:      {"Jump", []int{8}},
	OpJumpIf:    {"JumpIf", []int{8}},
	OpJumpIfNot: {"JumpIfNot", []int{8}},

	OpPop: {"Pop", []int{}},

	OpSetGlobal: {"SetGlobal", []int{8}},
	OpGetGlobal: {"GetGlobal", []int{8}},

	OpPushArgs: {"PushArgs", []int{8}},
	OpDropArgs: {"DropArgs", []int{8}},
	OpGetArg:   {"GetArg", []int{8}},

	OpRef:     {"Ref", []int{8}},
	OpCall:    {"Call", []int{}},
	OpRecCall: {"RecCall", []int{}},
	OpReturn:  {"Return", []int{}},

	OpRead:  {"Read", []int{}},
	OpWrite: {"Write", []int{}},

	OpEnd:   {"End", []int{}},
	OpHalt:  {"Halt", []int{}},
	OpDebug: {"Debug", []int{8}},
}

// Size returns the number of bytes for all arguments of an instruction.
func (m *OpMeta) Size() int {
	sz := 0
	for _, as := range m.Args {
		sz += as
	}
	return sz
}

// LookupMeta returns meta data for an opcode or an error if the code is
// undefined. The meta data contains the human-readable name of the operation
// and the length in bytes of each of its arguments.
func LookupMeta(op Op) (*OpMeta, error) {
	if m, ok := meta[op]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("opcode [%d] undefined", op)
}

// // TODO: ~> bin?
// func Correct(code []byte, pos int, new uint64) {
// 	for i := 0; i < 8; i++ {
// 		code[pos+7-i] = byte(new >> uint64(8*i))
// 	}
// }

// Instr creates a new instruction from an opcode and a variable number of
// arguments.
func Instr(op Op, args ...uint64) []byte {
	m := meta[op]
	sz := 1 + m.Size()
	ins := make([]byte, sz)
	pos := 1

	ins[0] = op
	for i, as := range m.Args {
		switch as {
		case 1:
			ins[pos] = uint8(args[i])
		case 2:
			binary.BigEndian.PutUint16(ins[pos:pos+2], uint16(args[i]))
		case 4:
			binary.BigEndian.PutUint32(ins[pos:pos+4], uint32(args[i]))
		case 8:
			binary.BigEndian.PutUint64(ins[pos:pos+8], args[i])
		}
		pos += as
	}
	return ins
}

func Bool(n bool) []byte {
	if n {
		return Instr(OpTrue)
	}
	return Instr(OpFalse)
}

func Str(s string) []byte {
	b := []byte(s)
	ln := len(b)
	sz := 9 + ln
	ins := make([]byte, sz)
	ins[0] = OpStr
	binary.BigEndian.PutUint64(ins[1:9], uint64(ln))
	copy(ins[9:sz], b)
	return ins
}

// Concat joins an array of instructions.
func Concat(is [][]byte) []byte {
	return bytes.Join(is, []byte{})
}

// ConcatVar joins a variable number of instructions.
func ConcatVar(is ...[]byte) []byte {
	return Concat(is)
}
