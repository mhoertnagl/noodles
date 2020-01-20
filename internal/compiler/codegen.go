package compiler

import (
	"encoding/binary"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type CodeGen interface {
	Append(i vm.Ins)
	Instr(op vm.Op, args ...uint64)
	OpAt(offset int) vm.Op
	AppendFunctions(fns []*fnDef)
	CorrectFunctionCalls(fns []*fnDef)
	Emit() vm.Ins
	Len() uint64
}

type codeGen struct {
	code vm.Ins
}

func NewCodeGen() CodeGen {
	return &codeGen{
		make(vm.Ins, 0),
	}
}

func (c *codeGen) Append(i vm.Ins) {
	c.code = vm.ConcatVar(c.code, i)
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
		if op == vm.OpRef {
			id := binary.BigEndian.Uint64(c.code[i+1 : i+9])
			fn := fns[id]
			vm.Correct(c.code, i+1, fn.addr)
		}
		i += mt.Size() + 1
	}
}

func (c *codeGen) Emit() vm.Ins {
	return c.code
}

func (c *codeGen) OpAt(offset int) vm.Op {
	if offset >= 0 {
		return c.code[offset]
	}
	return c.code[len(c.code)+offset]
}

func (c *codeGen) Len() uint64 {
	return uint64(len(c.code))
}
