package compiler

import (
	"fmt"

	"github.com/mhoertnagl/splis2/internal/vm"
)

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
	case *ListNode:
		return c.compileList(n)
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

func (c *compiler) compileList(n *ListNode) vm.Ins {
	args := n.Items[1:]

	switch len(args) {
	case 0:
		// Empty sum (+) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton sum (+ x) yields x.
		return c.Compile(args[0])
	default:
		code := make([]vm.Ins, 0)
		code = append(code, c.Compile(args[0]))
		// fmt.Printf("1: %v\n", code)
		for _, arg := range args[1:] {
			code = append(code, c.Compile(arg))
			// fmt.Printf("n: %v\n", code)
		}
		code = append(code, c.compileCoreOp(n.Items[0]))
		return vm.Concat(code)
	}
}

func (c *compiler) compileCoreOp(op Node) vm.Ins {
	switch sym := op.(type) {
	case *SymbolNode:
		switch sym.Name {
		case "+":
			return vm.Instr(vm.OpAdd)
		case "-":
			return vm.Instr(vm.OpSub)
		default:
			panic(fmt.Sprintf("Cannot compile core function [%v]", op))
		}
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", op))
	}
}
