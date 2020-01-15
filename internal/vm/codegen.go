package vm

type CodeGen interface {
	Append(i Ins)
	Instr(op Op, args ...uint64)
	OpAt(offset int) Op
	Emit() Ins
	Len() uint64
}

type codeGen struct {
	code Ins
}

func NewCodeGen() CodeGen {
	return &codeGen{
		make(Ins, 0),
	}
}

func (c *codeGen) Append(i Ins) {
	c.code = ConcatVar(c.code, i)
}

func (c *codeGen) Instr(op Op, args ...uint64) {
	c.code = ConcatVar(c.code, Instr(op, args...))
}

func (c *codeGen) Emit() Ins {
	return c.code
}

func (c *codeGen) OpAt(offset int) Op {
	if offset >= 0 {
		return c.code[offset]
	}
	return c.code[len(c.code)+offset]
}

func (c *codeGen) Len() uint64 {
	return uint64(len(c.code))
}
