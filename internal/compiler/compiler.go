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
	items := n.Items
	if len(items) == 0 {
		panic("Empty list")
	}
	args := n.Items[1:]
	switch sym := items[0].(type) {
	case *SymbolNode:
		switch sym.Name {
		case "+":
			return c.compileAdd(args)
		case "-":
			return c.compileSub(args)
		case "*":
			return c.compileMul(args)
		case "/":
			return c.compileDiv(args)
		default:
			panic(fmt.Sprintf("Cannot compile core function [%v]", sym))
		}
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", sym))
	}

}

func (c *compiler) compileAdd(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty sum (+) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton sum (+ x) yields x.
		return c.Compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// addition operations except for the first pair. The resulting sequence of
		// instructions is then:
		//
		//   <(+ x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpAdd, <x3>, OpAdd, <x4>, OpAdd, ...
		//
		code := make([]vm.Ins, 0)
		code = append(code, c.Compile(args[0]))
		for _, arg := range args[1:] {
			code = append(code, c.Compile(arg))
			code = append(code, vm.Instr(vm.OpAdd))
		}
		return vm.Concat(code)
	}
}

func (c *compiler) compileSub(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty difference (-) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton difference (- x) yields (- 0 x) which if effectively -x.
		code := make([]vm.Ins, 0)
		code = append(code, vm.Instr(vm.OpConst, 0))
		code = append(code, c.Compile(args[0]))
		code = append(code, vm.Instr(vm.OpSub))
		return vm.Concat(code)
	case 2:
		// Only supports at most two operands and computes their difference.
		//
		//   <(- x1 x2)> := <x1>, <x2>, OpSub
		//
		code := make([]vm.Ins, 0)
		code = append(code, c.Compile(args[0]))
		code = append(code, c.Compile(args[1]))
		code = append(code, vm.Instr(vm.OpSub))
		return vm.Concat(code)
	default:
		panic("Too many arguments")
	}
}

func (c *compiler) compileMul(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty product (*) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton product (* x) yields x.
		return c.Compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// multiplication operations except for the first pair. The resulting
		// sequence of instructions is then:
		//
		//   <(* x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpMul, <x3>, OpMul, <x4>, OpMul, ...
		//
		code := make([]vm.Ins, 0)
		code = append(code, c.Compile(args[0]))
		for _, arg := range args[1:] {
			code = append(code, c.Compile(arg))
			code = append(code, vm.Instr(vm.OpMul))
		}
		return vm.Concat(code)
	}
}

func (c *compiler) compileDiv(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty division (/) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton difference (/ x) yields (/ 1 x) which if effectively 1/x.
		code := make([]vm.Ins, 0)
		code = append(code, vm.Instr(vm.OpConst, 1))
		code = append(code, c.Compile(args[0]))
		code = append(code, vm.Instr(vm.OpDiv))
		return vm.Concat(code)
	case 2:
		// Only supports at most two operands and computes their quotient.
		//
		//   <(/ x1 x2)> := <x1>, <x2>, OpDiv
		//
		code := make([]vm.Ins, 0)
		code = append(code, c.Compile(args[0]))
		code = append(code, c.Compile(args[1]))
		code = append(code, vm.Instr(vm.OpDiv))
		return vm.Concat(code)
	default:
		panic("Too many arguments")
	}
}
