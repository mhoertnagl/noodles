package cmp

import (
	"encoding/binary"

	"github.com/mhoertnagl/splis2/internal/vm"
)

// CodeGen defines the interface of the code generator.
type CodeGen interface {

	// Prepend adds a code sequence to the beginning of the code sequence.
	Prepend(i vm.Ins)

	// Append adds a code sequence to the end of the code sequence.
	Append(i vm.Ins)

	// Instr adds a new instruction to the beginning of the code sequence.
	PrependInstr(op vm.Op, args ...uint64)

	// Instr adds a new instruction to the end of the code stream.
	Instr(op vm.Op, args ...uint64)

	// AppendFunctions appends a function definition the the list of definitions.
	AppendFunctions(fns []*fnDef)

	// CorrectFunctionCalls updates all function references in the code sequence
	// to point to the real function addresses.
	CorrectFunctionCalls(fns []*fnDef)

	// OpAt returns the op code at byte index offset.
	OpAt(offset int) vm.Op

	// Emit returns the actual code sequence.
	Emit() vm.Ins

	// Len returns the length of the code sequence.
	Len() uint64
}

type codeGen struct {
	code vm.Ins
}

// NewCodeGen creates a new code generator.
func NewCodeGen() CodeGen {
	return &codeGen{
		make(vm.Ins, 0),
	}
}

func (c *codeGen) Prepend(i vm.Ins) {
	c.code = vm.ConcatVar(i, c.code)
}

func (c *codeGen) Append(i vm.Ins) {
	c.code = vm.ConcatVar(c.code, i)
}

func (c *codeGen) PrependInstr(op vm.Op, args ...uint64) {
	c.code = vm.ConcatVar(vm.Instr(op, args...), c.code)
}

func (c *codeGen) Instr(op vm.Op, args ...uint64) {
	c.code = vm.ConcatVar(c.code, vm.Instr(op, args...))
}

func (c *codeGen) AppendFunctions(fns []*fnDef) {
	for _, fd := range fns {
		fd.addr = c.Len()
		c.Append(fd.code)
	}
}

func (c *codeGen) CorrectFunctionCalls(fns []*fnDef) {
	for i := 0; i < len(c.code); {
		op := c.OpAt(i)
		mt, err := vm.LookupMeta(op)
		if err != nil {
			panic(err)
		}
		switch op {
		case vm.OpRef:
			// The only argument of Ref is the index of the referenced function
			// definition entry in fns. The definitions entry contains the real
			// address of the function in the code stream. Replace the index in
			// the code stream with the real address.
			id := binary.BigEndian.Uint64(c.code[i+1 : i+9])
			fn := fns[id]
			vm.Correct(c.code, i+1, fn.addr)
		case vm.OpStr:
			// String commands are of variable length. The first argument is the
			// length of the string. We need to skip over the string as well and
			// thus add the string length to the position pointer i.
			l := binary.BigEndian.Uint64(c.code[i+1 : i+9])
			i += int(l)
		}
		i += mt.Size() + 1
	}
}

func (c *codeGen) OpAt(offset int) vm.Op {
	if offset >= 0 {
		return c.code[offset]
	}
	return c.code[len(c.code)+offset]
}

func (c *codeGen) Emit() vm.Ins {
	return c.code
}

func (c *codeGen) Len() uint64 {
	return uint64(len(c.code))
}
