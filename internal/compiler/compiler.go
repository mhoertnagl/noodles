package compiler

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type Compiler interface {
	Compile(node Node) vm.Ins
}

type fnDef struct {
	addr uint64
	code vm.Ins
}

type compiler struct {
	hg  hash.Hash64
	fns []*fnDef
}

func NewCompiler() Compiler {
	return &compiler{
		hg:  fnv.New64(),
		fns: make([]*fnDef, 0),
	}
}

// TODO: Alternative scheme
//       OpCallNew - generates a new environment.
//       OpCall    - does not generate an environment.
//       With this we can implement leaf functions and tail calls.

func (c *compiler) Compile(node Node) vm.Ins {
	code := NewCodeGen()
	code.Append(c.compile(node))
	// Marks the end of non-function code. This will halt the CPU. Code beyond
	// that point contains function definitions only.
	code.Instr(vm.OpHalt)
	code.AppendFunctions(c.fns)
	code.CorrectFunctionCalls(c.fns)
	return code.Emit()
}

func (c *compiler) compile(node Node) vm.Ins {
	switch n := node.(type) {
	case bool:
		return c.compileBooleanLiteral(n)
	case int64:
		return c.compileIntegerLiteral(n)
	case *SymbolNode:
		return c.compileSymbol(n)
	case *VectorNode:
		return c.compileVector(n)
	case *ListNode:
		return c.compileList(n)
	}
	panic(fmt.Sprintf("Unsupported node [%v]", node))
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

func (c *compiler) compileSymbol(n *SymbolNode) vm.Ins {
	return vm.Instr(vm.OpGetLocal, c.hashSymbol(n))
}

func (c *compiler) compileVector(n *VectorNode) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEmptyVector)
	for i := len(n.Items) - 1; i >= 0; i-- {
		code.Append(c.compile(n.Items[i]))
		code.Instr(vm.OpCons)
	}
	return code.Emit()
}

func (c *compiler) compileList(n *ListNode) vm.Ins {
	items := n.Items
	if len(items) == 0 {
		panic("Empty list")
	}
	args := items[1:]
	switch x := items[0].(type) {
	case *SymbolNode:
		switch x.Name {
		case "+":
			return c.compileAdd(args)
		case "-":
			return c.compileSub(args)
		case "*":
			return c.compileMul(args)
		case "/":
			return c.compileDiv(args)
		case "let":
			return c.compileLet(args)
		case "def":
			return c.compileDef(args)
		case "if":
			return c.compileIf(args)
		case "do":
			return c.compileDo(args)
		case "fn":
			return c.compileFn(args)
		case "::":
			return c.compileCons(args)
		default:
			return c.compileCall(x, args)
		}
	case *ListNode:
		return c.compileListCall(x, args)
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", x))
	}
}

func (c *compiler) compileAdd(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty sum (+) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton sum (+ x) yields x.
		return c.compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// addition operations except for the first pair. The resulting sequence of
		// instructions is then:
		//
		//   <(+ x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpAdd, <x3>, OpAdd, <x4>, OpAdd, ...
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0]))
		for _, arg := range args[1:] {
			code.Append(c.compile(arg))
			code.Instr(vm.OpAdd)
		}
		return code.Emit()
	}
}

func (c *compiler) compileSub(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty difference (-) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton difference (- x) yields (- 0 x) which if effectively -x.
		code := NewCodeGen()
		code.Instr(vm.OpConst, 0)
		code.Append(c.compile(args[0]))
		code.Instr(vm.OpSub)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their difference.
		//
		//   <(- x1 x2)> := <x1>, <x2>, OpSub
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0]))
		code.Append(c.compile(args[1]))
		code.Instr(vm.OpSub)
		return code.Emit()
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
		return c.compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// multiplication operations except for the first pair. The resulting
		// sequence of instructions is then:
		//
		//   <(* x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpMul, <x3>, OpMul, <x4>, OpMul, ...
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0]))
		for _, arg := range args[1:] {
			code.Append(c.compile(arg))
			code.Instr(vm.OpMul)
		}
		return code.Emit()
	}
}

func (c *compiler) compileDiv(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty division (/) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton difference (/ x) yields (/ 1 x) which if effectively 1/x.
		code := NewCodeGen()
		code.Instr(vm.OpConst, 1)
		code.Append(c.compile(args[0]))
		code.Instr(vm.OpDiv)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their quotient.
		//
		//   <(/ x1 x2)> := <x1>, <x2>, OpDiv
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0]))
		code.Append(c.compile(args[1]))
		code.Instr(vm.OpDiv)
		return code.Emit()
	default:
		panic("Too many arguments")
	}
}

func (c *compiler) compileLet(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[let] requires exactly two arguments")
	}
	if bs, ok := args[0].(*ListNode); ok {
		if len(bs.Items)%2 == 1 {
			panic("[let] reqires an even number of bindings")
		}
		code := NewCodeGen()
		code.Instr(vm.OpNewEnv)
		// TODO: Separate function.
		for i := 0; i < len(bs.Items); i += 2 {
			if sym, ok2 := bs.Items[i].(*SymbolNode); ok2 {
				code.Append(c.compile(bs.Items[i+1]))
				hsh := c.hashSymbol(sym)
				code.Instr(vm.OpSetLocal, hsh)
			} else {
				panic(fmt.Sprintf("[let] cannot bind to [%v]", sym))
			}
		}
		code.Append(c.compile(args[1]))
		code.Instr(vm.OpPopEnv)
		return code.Emit()
	}
	panic("[let] requires first argument to be a list of bindings")
}

func (c *compiler) compileDef(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[def] requires exactly two arguments")
	}
	if sym, ok := args[0].(*SymbolNode); ok {
		code := NewCodeGen()
		code.Append(c.compile(args[1]))
		hsh := c.hashSymbol(sym)
		code.Instr(vm.OpSetGlobal, hsh)
		return code.Emit()
	}
	panic("[def] requires first argument to be a symbol")
}

func (c *compiler) compileIf(args []Node) vm.Ins {
	if len(args) != 2 && len(args) != 3 {
		panic("[if] requires exactly two or three arguments")
	}
	code := NewCodeGen()
	switch len(args) {
	case 2:
		cnd := c.compile(args[0])
		cns := c.compile(args[1])
		cnsLen := uint64(len(cns))
		code.Append(cnd)
		code.Instr(vm.OpJumpIfNot, cnsLen)
		code.Append(cns)
	case 3:
		cnd := c.compile(args[0])
		cns := c.compile(args[1])
		alt := c.compile(args[2])
		cnsLen := uint64(len(cns)) + 9 // Add the length of the jmp instruction.
		altLen := uint64(len(alt))
		code.Append(cnd)
		code.Instr(vm.OpJumpIfNot, cnsLen)
		code.Append(cns)
		code.Instr(vm.OpJump, altLen)
		code.Append(alt)
	}
	return code.Emit()
}

func (c *compiler) compileDo(args []Node) vm.Ins {
	code := NewCodeGen()
	for _, arg := range args {
		code.Append(c.compile(arg))
	}
	return code.Emit()
}

// func (c *compiler) compileCond(args []Node) vm.Ins {
//
// }

func (c *compiler) compileFn(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[fn] expects exactly 2 arguments")
	}
	fd := &fnDef{}
	switch x := args[0].(type) {
	case *ListNode:
		fd.code = c.compileFn2(x.Items, args[1])
	case *VectorNode:
		fd.code = c.compileFn2(x.Items, args[1])
	default:
		panic("[fn] first argument must be a list or vector")
	}
	id := len(c.fns)
	c.fns = append(c.fns, fd)
	return vm.Instr(vm.OpRef, uint64(id))
}

func (c *compiler) compileFn2(params []Node, body Node) vm.Ins {
	code := NewCodeGen()
	switch len(params) {
	case 0:
		code.Append(c.compile(body))
	default:
		code.Instr(vm.OpNewEnv)
		c.compileFnParams(code, params)
		code.Append(c.compile(body))
		code.Instr(vm.OpPopEnv)
	}
	code.Instr(vm.OpReturn)
	return code.Emit()
}

func (c *compiler) compileFnParams(code CodeGen, params []Node) {
	for i := len(params) - 1; i >= 0; i-- {
		switch x := params[i].(type) {
		case *SymbolNode:
			code.Instr(vm.OpSetLocal, c.hashSymbol(x))
		default:
			panic(fmt.Sprintf("[fn] parameter [%d] is not a symbol", i))
		}
	}
}

func (c *compiler) compileCons(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[::] expects exactly 2 arguments")
	}
	code := NewCodeGen()
	code.Append(c.compile(args[1]))
	code.Append(c.compile(args[0]))
	code.Instr(vm.OpCons)
	return code.Emit()
}

func (c *compiler) compileCall(sym *SymbolNode, args []Node) vm.Ins {
	code := NewCodeGen()
	c.compileNodes(code, args)
	code.Instr(vm.OpGetGlobal, c.hashSymbol(sym))
	code.Instr(vm.OpCall)
	return code.Emit()
}

func (c *compiler) compileListCall(lst *ListNode, args []Node) vm.Ins {
	code := NewCodeGen()
	c.compileNodes(code, args)
	code.Append(c.compileList(lst))
	// TODO: Is this correct in any case?
	// What about: (((fn [] (fn [x] (+ x 1)))) 1) we need to emit OpCall twice.
	if code.OpAt(-1) != vm.OpCall {
		code.Instr(vm.OpCall)
	}
	return code.Emit()
}

func (c *compiler) compileNodes(code CodeGen, nodes []Node) {
	for _, node := range nodes {
		code.Append(c.compile(node))
	}
}

func (c *compiler) hashSymbol(sym *SymbolNode) uint64 {
	c.hg.Reset()
	c.hg.Write([]byte(sym.Name))
	return c.hg.Sum64()
}
