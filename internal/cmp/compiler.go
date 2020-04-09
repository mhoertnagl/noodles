package cmp

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mhoertnagl/splis2/internal/util"
	"github.com/mhoertnagl/splis2/internal/vm"
)

// TODO: let bindings should have their own symbol table.
//       Fix Indexof to account for this fact.
// TODO: Define a function that return the special forms for (+, *) and (-, /)
// TODO: Variadic +, *, list, ...
// TODO: TR
// TODO: TCO
// TODO: Closure

// TODO: Funzt nicht für beliebige Funktionsaufrufe.

// TODO: *STDOUT*
// TODO: write
// TODO: str -> use printer to turn value into a string.
// TODO: *STDIN*
// TODO: read
// TODO: :::
// TODO: quot
// TODO: mod
// TODO: join (strings)

// TODO: https://yourbasic.org/golang/bitwise-operator-cheat-sheet/

type fnDef struct {
	addr uint64
	code vm.Ins
}

type specFun func([]Node, *SymTable) vm.Ins

type specDefs map[string]specFun

func (d specDefs) add(name string, fn specFun) {
	d[name] = fn
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
	specs specDefs
	prims primDefs
	fns   []*fnDef
	//sym   *SymTable
}

func NewCompiler() *Compiler {
	c := &Compiler{
		hg:  fnv.New64(),
		fns: make([]*fnDef, 0),
		//sym:   NewSymTable(),
	}

	c.specs = specDefs{}
	c.specs.add("+", c.compileAdd)
	c.specs.add("-", c.compileSub)
	c.specs.add("*", c.compileMul)
	c.specs.add("/", c.compileDiv)
	c.specs.add("let", c.compileLet)
	c.specs.add("def", c.compileDef)
	c.specs.add("if", c.compileIf)
	c.specs.add("do", c.compileDo)
	c.specs.add("fn", c.compileFn)
	c.specs.add("and", c.compileAnd)
	c.specs.add("or", c.compileOr)

	c.prims = primDefs{}
	c.prims.add("fst", vm.OpFst, 1, false)
	c.prims.add("rest", vm.OpRest, 1, false)
	c.prims.add("len", vm.OpLength, 1, false)
	c.prims.add("<", vm.OpLT, 2, false)
	c.prims.add("<=", vm.OpLE, 2, false)
	c.prims.add(">", vm.OpLT, 2, true)
	c.prims.add(">=", vm.OpLE, 2, true)
	c.prims.add("=", vm.OpEQ, 2, false)
	c.prims.add("!=", vm.OpNE, 2, false)
	c.prims.add("not", vm.OpNot, 1, false)
	c.prims.add("::", vm.OpCons, 2, true)
	c.prims.add("dissolve", vm.OpDissolve, 1, false)

	return c
}

func (c *Compiler) Compile(node Node) vm.Ins {
	sym := NewSymTable()
	code := NewCodeGen()
	code.Append(c.compile(node, sym))
	// Marks the end of non-function code. This will halt the CPU. Code beyond
	// that point contains function definitions only.
	code.Instr(vm.OpHalt)
	code.AppendFunctions(c.fns)
	code.CorrectFunctionCalls(c.fns)
	return code.Emit()
}

func (c *Compiler) compile(node Node, sym *SymTable) vm.Ins {
	switch n := node.(type) {
	case bool:
		return vm.Bool(n)
	case int64:
		return vm.Instr(vm.OpConst, uint64(n))
	case string:
		return vm.Str(n)
	case *SymbolNode:
		return c.compileSymbol(n, sym)
	case []Node:
		return c.compileVector(n, sym)
	case *ListNode:
		return c.compileList(n, sym)
	}
	panic(fmt.Sprintf("Compiler: Unsupported node [%v:%T]", node, node))
}

func (c *Compiler) compileSymbol(n *SymbolNode, sym *SymTable) vm.Ins {
	if idx, ok := sym.IndexOf(n.Name); ok {
		return vm.Instr(vm.OpGetArg, uint64(idx))
	}
	return vm.Instr(vm.OpGetGlobal, c.hashSymbol(n))
}

func (c *Compiler) compileVector(n []Node, sym *SymTable) vm.Ins {
	switch len(n) {
	case 0:
		return vm.Instr(vm.OpEmptyVector)
	default:
		code := NewCodeGen()
		code.Instr(vm.OpEnd)
		for i := len(n) - 1; i >= 0; i-- {
			code.Append(c.compile(n[i], sym))
		}
		code.Instr(vm.OpList)
		return code.Emit()
	}
}

func (c *Compiler) compileList(n *ListNode, sym *SymTable) vm.Ins {
	if n.Empty() {
		panic("Empty list")
	}
	switch x := n.First().(type) {
	case *SymbolNode:
		if x.Name == "debug" {
			mode := n.Rest()[0].(int64)
			return vm.Instr(vm.OpDebug, uint64(mode))
		}
		// Special forms handle their arguments in various ways. The arguments
		// may not get compiled in sequence.
		if spec, ok := c.specs[x.Name]; ok {
			return spec(n.Rest(), sym)
		}
		// Primitive functions follow the same pattern: first compile all the
		// arguments (either ascending or descending), then append a single VM
		// instruction.
		if prim, ok := c.prims[x.Name]; ok {
			return c.compilePrim(x.Name, prim, n.Rest(), sym)
		}
		// Compile a call to the global function.
		return c.compileCall(x, n.Rest(), sym)
	case *ListNode:
		// Compile the list. We expect the result of the computation to be a
		// REF value which we can then call.
		return c.compileListCall(x, n.Rest(), sym)
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", x))
	}
}

func (c *Compiler) compilePrim(name string, prim primDef, args []Node, sym *SymTable) vm.Ins {
	if len(args) != prim.nargs {
		panic(fmt.Sprintf("[%s] requires exactly [%d] arguments", name, prim.nargs))
	}
	code := NewCodeGen()
	if prim.rev {
		c.compileNodesReverse(code, args, sym)
	} else {
		c.compileNodes(code, args, sym)
	}
	code.Instr(prim.op)
	return code.Emit()
}

func (c *Compiler) compileAdd(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty sum (+) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton sum (+ x) yields x.
		return c.compile(args[0], sym)
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// addition operations except for the first pair. The resulting sequence of
		// instructions is then:
		//
		//   <(+ x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpAdd, <x3>, OpAdd, <x4>, OpAdd, ...
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0], sym))
		for _, arg := range args[1:] {
			code.Append(c.compile(arg, sym))
			code.Instr(vm.OpAdd)
		}
		return code.Emit()
	}
}

func (c *Compiler) compileSub(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty difference (-) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton difference (- x) yields (- 0 x) which if effectively -x.
		code := NewCodeGen()
		code.Instr(vm.OpConst, 0)
		code.Append(c.compile(args[0], sym))
		code.Instr(vm.OpSub)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their difference.
		//
		//   <(- x1 x2)> := <x1>, <x2>, OpSub
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0], sym))
		code.Append(c.compile(args[1], sym))
		code.Instr(vm.OpSub)
		return code.Emit()
	default:
		panic("[-] Too many arguments")
	}
}

func (c *Compiler) compileMul(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty product (*) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton product (* x) yields x.
		return c.compile(args[0], sym)
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// multiplication operations except for the first pair. The resulting
		// sequence of instructions is then:
		//
		//   <(* x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpMul, <x3>, OpMul, <x4>, OpMul, ...
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0], sym))
		for _, arg := range args[1:] {
			code.Append(c.compile(arg, sym))
			code.Instr(vm.OpMul)
		}
		return code.Emit()
	}
}

func (c *Compiler) compileDiv(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty division (/) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton difference (/ x) yields (/ 1 x) which if effectively 1/x.
		code := NewCodeGen()
		code.Instr(vm.OpConst, 1)
		code.Append(c.compile(args[0], sym))
		code.Instr(vm.OpDiv)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their quotient.
		//
		//   <(/ x1 x2)> := <x1>, <x2>, OpDiv
		//
		code := NewCodeGen()
		code.Append(c.compile(args[0], sym))
		code.Append(c.compile(args[1], sym))
		code.Instr(vm.OpDiv)
		return code.Emit()
	default:
		panic("[/] Too many arguments")
	}
}

func (c *Compiler) compileAnd(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty and (and) yields true.
		return vm.Instr(vm.OpTrue)
	case 1:
		// Singleton and (and x) yields x.
		return c.compile(args[0], sym)
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
			code.Prepend(c.compile(args[i], sym))
			code.PrependInstr(vm.OpJumpIfNot, code.Len())
		}
		// Preprend the first argument.
		code.Prepend(c.compile(args[0], sym))
		// Each evaluation except for the last one that yielded False will jump to
		// this istruction that puts False on the stack.
		code.Instr(vm.OpFalse)
		return code.Emit()
	}
}

func (c *Compiler) compileOr(args []Node, sym *SymTable) vm.Ins {
	switch len(args) {
	case 0:
		// Empty or (or) yields false.
		return vm.Instr(vm.OpFalse)
	case 1:
		// Singleton and (or x) yields x.
		return c.compile(args[0], sym)
	default:
		// Analogous to and compilation
		code := NewCodeGen()
		code.Instr(vm.OpJump, 1)
		for i := len(args) - 1; i > 0; i-- {
			code.Prepend(c.compile(args[i], sym))
			code.PrependInstr(vm.OpJumpIf, code.Len())
		}
		code.Prepend(c.compile(args[0], sym))
		code.Instr(vm.OpTrue)
		return code.Emit()
	}
}

func (c *Compiler) compileLet(args []Node, sym *SymTable) vm.Ins {
	if len(args) != 2 {
		panic("[let] requires exactly two arguments")
	}
	if bs, ok := args[0].(*ListNode); ok {
		if len(bs.Items)%2 == 1 {
			panic("[let] reqires an even number of bindings")
		}
		code := NewCodeGen()

		// TODO: separate symbol table would be better.
		// TODO: Problem when shadowing a variable.
		locals := make([]string, 0)
		for i := 0; i < len(bs.Items); i += 2 {
			// for i := len(bs.Items) - 1; i >= 0; i -= 2 {
			if s, ok2 := bs.Items[i].(*SymbolNode); ok2 {
				// Keep track of the let bindings. We will remove them when we fall
				// out of scope.
				locals = append(locals, s.Name)
				// Add the local binding to the symbol table.
				sym.AddVar(s.Name)

				code.Append(c.compile(bs.Items[i+1], sym))
				// Add the let bindings ont at a time so that subsequent bindings
				// will be able to access the privously defined let bindings.
				code.Instr(vm.OpPushArgs, 1)
			} else {
				panic(fmt.Sprintf("[let] cannot bind to [%v]", s))
			}
		}

		code.Append(c.compile(args[1], sym))
		// A let binding does not posess a separate frame. We need to drop all
		// introduced let bindings before we continue.
		code.Instr(vm.OpDropArgs, uint64(len(locals)))

		// Remove the let bindings from the symbol table as well.
		sym.Remove(locals)

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
func (c *Compiler) compileDef(args []Node, sym *SymTable) vm.Ins {
	if len(args) != 2 {
		panic("[def] requires exactly two arguments")
	}
	if s, ok := args[0].(*SymbolNode); ok {
		code := NewCodeGen()
		code.Append(c.compile(args[1], sym))
		hsh := c.hashSymbol(s)
		code.Instr(vm.OpSetGlobal, hsh)
		return code.Emit()
	}
	panic("[def] requires first argument to be a symbol")
}

func (c *Compiler) compileIf(args []Node, sym *SymTable) vm.Ins {
	if len(args) != 2 && len(args) != 3 {
		panic("[if] requires either two or three arguments")
	}
	code := NewCodeGen()
	switch len(args) {
	case 2:
		cnd := c.compile(args[0], sym)
		cns := c.compile(args[1], sym)
		cnsLen := uint64(len(cns))
		code.Append(cnd)
		code.Instr(vm.OpJumpIfNot, cnsLen)
		code.Append(cns)
	case 3:
		cnd := c.compile(args[0], sym)
		cns := c.compile(args[1], sym)
		alt := c.compile(args[2], sym)
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

func (c *Compiler) compileDo(args []Node, sym *SymTable) vm.Ins {
	code := NewCodeGen()
	c.compileNodes(code, args, sym)
	return code.Emit()
}

func (c *Compiler) compileFn(args []Node, sym *SymTable) vm.Ins {
	if len(args) != 2 {
		panic("[fn] expects exactly 2 arguments")
	}
	fd := &fnDef{}
	switch x := args[0].(type) {
	case *ListNode:
		fd.code = c.compileFnBody(x.Items, args[1], sym)
	case []Node:
		fd.code = c.compileFnBody(x, args[1], sym)
	default:
		panic("[fn] first argument must be a list or vector")
	}
	id := len(c.fns)
	c.fns = append(c.fns, fd)
	return vm.Instr(vm.OpRef, uint64(id))
}

func (c *Compiler) compileFnBody(params []Node, body Node, sym *SymTable) vm.Ins {
	sub := sym.NewSymTable()
	code := NewCodeGen()
	switch len(params) {
	case 0:
		// Removes the function argument's end marker from the stack.
		code.Instr(vm.OpPop)
		code.Append(c.compile(body, sub))
	default:
		man, opt := c.compileParams(params)
		// Add the mandatory arguments to the local symbol table.
		sub.Add(man)
		// Push the mandatory arguments to the frames stack.
		code.Instr(vm.OpPushArgs, uint64(len(man)))
		// Check for an optional argument.
		if opt != "" {
			// Add the optional argument to the local symbol table.
			sub.AddVar(opt)
			// The LIST operation will append all remaining arguments to a vector.
			code.Instr(vm.OpList)
			// Then push the vector to the frames stack as well.
			code.Instr(vm.OpPushArgs, 1)
		}
		// Removes the function argument's end marker from the stack.
		code.Instr(vm.OpPop)
		code.Append(c.compile(body, sub))
	}
	code.Instr(vm.OpReturn)
	return code.Emit()
}

func (c *Compiler) compileParams(params []Node) ([]string, string) {
	names := c.verifyParams(params)
	pos := util.IndexOf(names, "&")
	if pos == -1 {
		return names, ""
	}
	if len(names) == pos+1 {
		panic(fmt.Sprintf("[fn] missing optional parameter"))
	}
	if len(names) > pos+2 {
		panic(fmt.Sprintf("[fn] excess optional parameter"))
	}
	return names[:pos], names[pos+1]
}

func (c *Compiler) verifyParams(params []Node) []string {
	names := make([]string, 0)
	for pos, param := range params {
		names = append(names, c.verifyParam(param, pos))
	}
	return names
}

func (c *Compiler) verifyParam(param Node, pos int) string {
	switch sym := param.(type) {
	case *SymbolNode:
		return sym.Name
	default:
		panic(fmt.Sprintf("[fn] parameter [%d] is not a symbol", pos))
	}
}

func (c *Compiler) compileCall(s *SymbolNode, args []Node, sym *SymTable) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEnd)
	c.compileNodesReverse(code, args, sym)
	code.Instr(vm.OpGetGlobal, c.hashSymbol(s))
	code.Instr(vm.OpCall)
	return code.Emit()
}

func (c *Compiler) compileListCall(lst *ListNode, args []Node, sym *SymTable) vm.Ins {
	code := NewCodeGen()
	code.Instr(vm.OpEnd)
	c.compileNodesReverse(code, args, sym)
	code.Append(c.compileList(lst, sym))
	code.Instr(vm.OpCall)
	return code.Emit()
}

func (c *Compiler) compileNodes(code CodeGen, nodes []Node, sym *SymTable) {
	for _, node := range nodes {
		code.Append(c.compile(node, sym))
	}
}

func (c *Compiler) compileNodesReverse(code CodeGen, nodes []Node, sym *SymTable) {
	for i := len(nodes) - 1; i >= 0; i-- {
		code.Append(c.compile(nodes[i], sym))
	}
}

func (c *Compiler) hashSymbol(sym *SymbolNode) uint64 {
	c.hg.Reset()
	c.hg.Write([]byte(sym.Name))
	return c.hg.Sum64()
}
