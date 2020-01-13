package vm

type CodeGen interface {
	Append(i Ins)
	Instr(op Op, args ...uint64)
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

// type codeGen struct {
// 	code []Ins
// }
//
// func NewCodeGen() CodeGen {
// 	return &codeGen{
// 		make([]Ins, 0),
// 	}
// }
//
// func (c *codeGen) Append(i Ins) {
// 	c.code = append(c.code, i)
// }
//
// func (c *codeGen) Instr(op Op, args ...uint64) {
// 	c.code = append(c.code, Instr(op, args...))
// }
//
// func (c *codeGen) Emit() Ins {
// 	return Concat(c.code)
// }

func (c *codeGen) Len() uint64 {
	return uint64(len(c.code))
}
