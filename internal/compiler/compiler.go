package compiler

import "github.com/mhoertnagl/splis2/internal/vm"

type Compiler interface {
	Compile(node Node) vm.Ins
}

type compiler struct {
}

func NewCompiler() Compiler {
	return &compiler{}
}

func (c *compiler) Compile(node Node) vm.Ins {
	switch n := node.(type) {
	case bool:
		return c.compileBooleanLiteral(n)
	case int64:
		return c.compileIntegerLiteral(n)
	}
	return nil
}

func (c *compiler) compileBooleanLiteral(n bool) vm.Ins {
	if n {
		return vm.Instr(vm.OpTrue)
	}
	return vm.Instr(vm.OpFalse)
}

func (c *compiler) compileIntegerLiteral(n int64) vm.Ins {
	return vm.Instr(vm.OpConst, uint64(n))
}

// func (c *compiler) compileList(n ListNode) vm.Ins {
//
// }
//
// func (c *compiler) compileBasicOp(n int64) vm.Ins {
// 	return vm.Instr(vm.OpConst, uint64(n))
// }
