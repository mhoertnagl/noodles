package vm

type CodeGen interface {
	Append(i Ins)
	Instr(op Op, args ...uint64)
	Emit() Ins
}

type codeGen struct {
	code []Ins
}

func NewCodeGen() CodeGen {
	return &codeGen{
		make([]Ins, 0),
	}
}

func (c *codeGen) Append(i Ins) {
	c.code = append(c.code, i)
}

func (c *codeGen) Instr(op Op, args ...uint64) {
	c.code = append(c.code, Instr(op, args...))
}

func (c *codeGen) Emit() Ins {
	return Concat(c.code)
}
