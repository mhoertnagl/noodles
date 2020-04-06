package cmp

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type fnDef struct {
	addr uint64
	code vm.Ins
}

type primDef struct {
	op    vm.Op
	nargs int
	rev   bool
}

type primDefs map[string]primDef

func (d primDefs) add(name string, op vm.Op, nargs int, rev bool) {
	d[name] = primDef{op: op, nargs: nargs, rev: rev}
}

type Compiler struct {
	hg    hash.Hash64
	prims primDefs
	fns   []*fnDef
}

func NewCompiler() *Compiler {
	prims := primDefs{}
	prims.add("fst", vm.OpFst, 1, false)
	prims.add("rest", vm.OpRest, 1, false)
	prims.add("len", vm.OpLength, 1, false)
	prims.add("<", vm.OpLT, 2, false)
	prims.add("<=", vm.OpLE, 2, false)
	prims.add(">", vm.OpLT, 2, true)
	prims.add(">=", vm.OpLE, 2, true)
	prims.add("=", vm.OpEQ, 2, false)
	prims.add("!=", vm.OpNE, 2, false)
	prims.add("not", vm.OpNot, 1, false)
	prims.add("::", vm.OpCons, 2, true)
	prims.add("dissolve", vm.OpDissolve, 1, false)

	return &Compiler{
		hg:    fnv.New64(),
		prims: prims,
		fns:   make([]*fnDef, 0),
	}
}

// TODO: Alternative scheme
//       OpCallNew - generates a new environment.
//       OpCall    - does not generate an environment.
//       With this we can implement leaf functions and tail calls.
// TODO: Variadic +, *, do, and, or, vector, list, ...
// TODO: Closure
// TODO: static scoping?

func (c *Compiler) Compile(node Node) vm.Ins {
	code := NewCodeGen()
	code.Append(c.compile(node))
	// Marks the end of non-function code. This will halt the CPU. Code beyond
	// that point contains function definitions only.
	code.Instr(vm.OpHalt)
	code.AppendFunctions(c.fns)
	code.CorrectFunctionCalls(c.fns)
	return code.Emit()
}

func (c *Compiler) compile(node Node) vm.Ins {
	switch n := node.(type) {
	case bool:
		return c.compileBooleanLiteral(n)
	case int64:
		return vm.Instr(vm.OpConst, uint64(n))
	case string:
		return vm.Str(n)
	case *SymbolNode:
		return c.compileSymbol(n)
	case []Node:
		return c.compileVector(n)
	case *ListNode:
		return c.compileList(n)
	}
	panic(fmt.Sprintf("Compiler: Unsupported node [%v:%T]", node, node))
}

func (c *Compiler) compileBooleanLiteral(n bool) vm.Ins {
	if n {
		return vm.Instr(vm.OpTrue)
	}
	return vm.Instr(vm.OpFalse)
}

func (c *Compiler) compileSymbol(n *SymbolNode) vm.Ins {
	return vm.Instr(vm.OpGetLocal, c.hashSymbol(n))
}

func (c *Compiler) compileVector(n []Node) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEmptyVector)
	for i := len(n) - 1; i >= 0; i-- {
		code.Append(c.compile(n[i]))
		code.Instr(vm.OpCons)
	}
	return code.Emit()
}

func (c *Compiler) compileList(n *ListNode) vm.Ins {
	if n.Empty() {
		panic("Empty list")
	}
	switch x := n.First().(type) {
	case *SymbolNode:
		// Special forms.
		switch x.Name {
		case "+":
			return c.compileAdd(n.Rest())
		case "-":
			return c.compileSub(n.Rest())
		case "*":
			return c.compileMul(n.Rest())
		case "/":
			return c.compileDiv(n.Rest())
		case "let":
			return c.compileLet(n.Rest())
		case "def":
			return c.compileDef(n.Rest())
		case "if":
			return c.compileIf(n.Rest())
		case "do":
			return c.compileDo(n.Rest())
		case "fn":
			return c.compileFn(n.Rest())
		case "and":
			return c.compileAnd(n.Rest())
		case "or":
			return c.compileOr(n.Rest())
		}
		// Primitive functions follow the same pattern: first compile all the
		// arguments, then append a single VM instruction.
		if prim, ok := c.prims[x.Name]; ok {
			return c.compilePrim(x.Name, prim, n.Rest())
		}
		return c.compileCall(x, n.Rest())
	case *ListNode:
		return c.compileListCall(x, n.Rest())
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", x))
	}
}

func (c *Compiler) compilePrim(name string, prim primDef, args []Node) vm.Ins {
	if len(args) != prim.nargs {
		panic(fmt.Sprintf("[%s] requires exactly [%d] arguments", name, prim.nargs))
	}
	code := NewCodeGen()
	if prim.rev {
		c.compileNodesReverse(code, args)
	} else {
		c.compileNodes(code, args)
	}
	code.Instr(prim.op)
	return code.Emit()
}

func (c *Compiler) compileAdd(args []Node) vm.Ins {
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

func (c *Compiler) compileSub(args []Node) vm.Ins {
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

func (c *Compiler) compileMul(args []Node) vm.Ins {
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

func (c *Compiler) compileDiv(args []Node) vm.Ins {
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

func (c *Compiler) compileAnd(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty and (and) yields true.
		return vm.Instr(vm.OpTrue)
	case 1:
		// Singleton and (and x) yields x.
		return c.compile(args[0])
	default:
		// Compiles an and expressh at least two aguments. The expression is
		// compiled in reverse. This way, we can tell the distance for the jumps at
		// after each agrument evaluation if it yields False.
		//
		//   <(and x1 x2 ... xn)> :=
		//        <x1>
		//        OpJumpIfNot @A
		//        <x2>
		//        OpJumpIfNot @A
		//        ...
		//        <xn>
		//        OpJump 1
		//     A: OpFalse
		//
		code := NewCodeGen()
		// Compile the second to last instruction that jumps over the False
		// constant. If all arguments evaluated to True then there will be True on
		// the stack. If the last argument evaluated to False then the expression is
		// False and False is on the stack.
		code.Instr(vm.OpJump, 1)
		// Prepend all but the first argument in reverse. The length of the code
		// in each iteration equals the distance to jump if evaluation of a argument
		// yields False.
		for i := len(args) - 1; i > 0; i-- {
			code.Prepend(c.compile(args[i]))
			code.PrependInstr(vm.OpJumpIfNot, code.Len())
		}
		// Preprend the first argument.
		code.Prepend(c.compile(args[0]))
		// Each evaluation except for the last one that yielded False will jump to
		// this istruction that puts False on the stack.
		code.Instr(vm.OpFalse)
		return code.Emit()
	}
}

func (c *Compiler) compileOr(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty or (or) yields false.
		return vm.Instr(vm.OpFalse)
	case 1:
		// Singleton and (or x) yields x.
		return c.compile(args[0])
	default:
		// Analogous to and compilation
		code := NewCodeGen()
		code.Instr(vm.OpJump, 1)
		for i := len(args) - 1; i > 0; i-- {
			code.Prepend(c.compile(args[i]))
			code.PrependInstr(vm.OpJumpIf, code.Len())
		}
		code.Prepend(c.compile(args[0]))
		code.Instr(vm.OpTrue)
		return code.Emit()
	}
}

func (c *Compiler) compileLet(args []Node) vm.Ins {
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

// compileDef compiles a global definition. Global definitions will be bound in
// the root environment and are available in the entire codebase for the entire
// lifetime of the program.
//
//   <(def x y)> :=
//        <y>
//        OpSetGlobal #x
//
func (c *Compiler) compileDef(args []Node) vm.Ins {
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

func (c *Compiler) compileIf(args []Node) vm.Ins {
	if len(args) != 2 && len(args) != 3 {
		panic("[if] requires either two or three arguments")
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

// func (c *compiler) compileCond(args []Node) vm.Ins {
//
// }

func (c *Compiler) compileDo(args []Node) vm.Ins {
	code := NewCodeGen()
	c.compileNodes(code, args)
	return code.Emit()
}

func (c *Compiler) compileFn(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[fn] expects exactly 2 arguments")
	}
	fd := &fnDef{}
	switch x := args[0].(type) {
	case *ListNode:
		fd.code = c.compileFnBody(x.Items, args[1])
	case []Node:
		fd.code = c.compileFnBody(x, args[1])
	default:
		panic("[fn] first argument must be a list or vector")
	}
	id := len(c.fns)
	c.fns = append(c.fns, fd)
	return vm.Instr(vm.OpRef, uint64(id))
}

func (c *Compiler) compileFnBody(params []Node, body Node) vm.Ins {
	code := NewCodeGen()
	switch len(params) {
	case 0:
		// Removes the function argument's end marker from the stack.
		code.Instr(vm.OpPop)
		code.Append(c.compile(body))
	default:
		code.Instr(vm.OpNewEnv)
		c.compileFnParams(code, params)
		// Removes the function argument's end marker from the stack.
		code.Instr(vm.OpPop)
		code.Append(c.compile(body))
		code.Instr(vm.OpPopEnv)
	}
	code.Instr(vm.OpReturn)
	return code.Emit()
}

func (c *Compiler) compileFnParams(code CodeGen, params []Node) {
	for pos, param := range params {
		switch sym := param.(type) {
		case *SymbolNode:
			c.compileFnParam(code, sym)
		default:
			panic(fmt.Sprintf("[fn] parameter [%d] is not a symbol", pos))
		}
	}
}

func (c *Compiler) compileFnParam(code CodeGen, sym *SymbolNode) {
	switch sym.Name {
	case "&":
		code.Instr(vm.OpList)
	default:
		code.Instr(vm.OpSetLocal, c.hashSymbol(sym))
	}
}

func (c *Compiler) compileCall(sym *SymbolNode, args []Node) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEnd)
	c.compileNodesReverse(code, args)
	code.Instr(vm.OpGetGlobal, c.hashSymbol(sym))
	code.Instr(vm.OpCall)
	return code.Emit()
}

func (c *Compiler) compileListCall(lst *ListNode, args []Node) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEnd)
	c.compileNodesReverse(code, args)
	code.Append(c.compileList(lst))
	code.Instr(vm.OpCall)
	return code.Emit()
}

func (c *Compiler) compileNodes(code CodeGen, nodes []Node) {
	for _, node := range nodes {
		code.Append(c.compile(node))
	}
}

func (c *Compiler) compileNodesReverse(code CodeGen, nodes []Node) {
	for i := len(nodes) - 1; i >= 0; i-- {
		code.Append(c.compile(nodes[i]))
	}
}

func (c *Compiler) hashSymbol(sym *SymbolNode) uint64 {
	c.hg.Reset()
	c.hg.Write([]byte(sym.Name))
	return c.hg.Sum64()
}
