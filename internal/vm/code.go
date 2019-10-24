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
	OpPop

	OpAdd
	OpSub
	OpMul
	OpDiv

	OpEq
	OpNe
	OpLT
	OpLE

	OpTrue
	OpFalse
	OpAnd
	OpOr

	OpJump        // <uint64>
	OpJumpIfTrue  // <bool> <uint64>
	OpJumpIfFalse // <bool> <uint64>

	OpNewEnv
	OpPopEnv
	OpSet
	OpGet
)

// OpMeta contains the human-readable name of the operation and the length in
// bytes of each of its arguments.
type OpMeta struct {
	Name string
	Args []int
}

var meta = map[Op]*OpMeta{
	OpConst:       {"const", []int{8}},
	OpPop:         {"pop", []int{}},
	OpAdd:         {"add", []int{}},
	OpSub:         {"sub", []int{}},
	OpMul:         {"mul", []int{}},
	OpDiv:         {"div", []int{}},
	OpTrue:        {"true", []int{}},
	OpFalse:       {"false", []int{}},
	OpJumpIfFalse: {"jumpf", []int{8}},
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
func LookupMeta(op byte) (*OpMeta, error) {
	if m, ok := meta[op]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("opcode [%d] undefined", op)
}

// Instr creates a new instruction from an opcode and a variable number of
// arguments.
func Instr(op Op, args ...uint64) Ins {
	m := meta[op]
	sz := m.Size() + 1
	ins := make(Ins, sz)
	pos := 1

	ins[0] = op
	for i, as := range m.Args {
		switch as {
		case 1:
			ins[pos] = uint8(args[i])
		case 2:
			binary.BigEndian.PutUint16(ins[pos:pos+4], uint16(args[i]))
		case 4:
			binary.BigEndian.PutUint32(ins[pos:pos+4], uint32(args[i]))
		case 8:
			binary.BigEndian.PutUint64(ins[pos:pos+8], args[i])
		}
		pos += as
	}
	return ins
}

// Concat joins an array of instructions.
func Concat(is []Ins) Ins {
	return bytes.Join(is, []byte{})
}

// ConcatVar joins a variable number of instructions.
func ConcatVar(is ...Ins) Ins {
	return Concat(is)
}
