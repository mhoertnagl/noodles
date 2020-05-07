package cmp

import (
	"fmt"
	"math"

	"github.com/mhoertnagl/splis2/internal/asm"
	"github.com/mhoertnagl/splis2/internal/util"
	"github.com/mhoertnagl/splis2/internal/vm"
)

// TODO: prelude docs and unit tests.

// TODO: Push everything in reverse order
// TODO: Explicit GT and GE

// TODO: Variadic list, ...

// TODO: *STDIN*
// TODO: read
// TODO: parse?

// TODO: Alternative syntax for :+ +: ? .+ == cons, +. == append, ++ == concat

// TODO: Should we keep let bindings?
// TODO: let binding recursive functions - see TestCompileLet5.
// TODO: let bindings should have their own symbol table.
//       Fix Indexof to account for this fact.

// TODO: Feed global names to disassembler and any place where they make sense.

// TODO: dissolve in various special forms like cond, if, let, def and, or will not work and should be disallowed.

// TODO: FRAMES Debug funzt nicht für beliebige Funktionsaufrufe da jeder FRAME
//       unterschiedliche Größe haben kann.

type Compiler struct {
	specs    specDefs
	prims    primDefs
	varPrims varPrimDefs
	fns      fnDefs
	defs     *defMap
	code     asm.AsmCode
	lblId    int
}

func NewCompiler() *Compiler {
	c := &Compiler{
		fns:   make(fnDefs, 0),
		defs:  newDefMap(),
		lblId: 0,
	}

	c.specs = specDefs{}
	c.specs.add("debug", c.compileDebug)
	c.specs.add("-", c.compileSub)
	c.specs.add("/", c.compileDiv)
	c.specs.add("set", c.compileSet)
	c.specs.add("let", c.compileLet)
	c.specs.add("def", c.compileDef)
	c.specs.add("if", c.compileIf)
	c.specs.add("cond", c.compileCond)
	c.specs.add("do", c.compileNodes)
	c.specs.add("fn", c.compileFn)
	c.specs.add("and", c.compileAnd)
	c.specs.add("or", c.compileOr)
	c.specs.add("rec", c.compileRec)

	c.prims = primDefs{}
	c.prims.add("nth", vm.OpNth, 2, false)
	c.prims.add("drop", vm.OpDrop, 2, false)
	c.prims.add("len", vm.OpLength, 1, false)
	c.prims.add("<", vm.OpLT, 2, false)
	c.prims.add("<=", vm.OpLE, 2, false)
	c.prims.add(">", vm.OpLT, 2, true)
	c.prims.add(">=", vm.OpLE, 2, true)
	c.prims.add("=", vm.OpEQ, 2, false)
	c.prims.add("!=", vm.OpNE, 2, false)
	c.prims.add("mod", vm.OpMod, 2, false)
	c.prims.add("not", vm.OpNot, 1, false)
	c.prims.add(".+", vm.OpCons, 2, true)
	c.prims.add("+.", vm.OpAppend, 2, false)
	c.prims.add("dissolve", vm.OpDissolve, 1, false)
	c.prims.add("halt", vm.OpHalt, 0, false)
	c.prims.add("explode", vm.OpExplode, 1, false)

	c.varPrims = varPrimDefs{}
	c.varPrims.add("+", vm.OpAdd, 0)
	c.varPrims.add("*", vm.OpMul, 0)
	c.varPrims.add("write", vm.OpWrite, 1)
	c.varPrims.add("++", vm.OpConcat, 0)
	c.varPrims.add("join", vm.OpJoin, 0)

	return c
}

func (c *Compiler) Compile(node Node) asm.AsmCode {
	sym := NewSymTable()
	ctx := NewCtx()
	c.code = make(asm.AsmCode, 0)
	c.compile(node, sym, ctx)
	return c.code
}

func (c *Compiler) compile(node Node, sym *SymTable, ctx *Ctx) {
	switch n := node.(type) {
	case bool:
		if n {
			c.instr(vm.OpTrue)
		} else {
			c.instr(vm.OpFalse)
		}
	case int64:
		c.instr(vm.OpConst, uint64(n))
	case float64:
		c.instr(vm.OpConstF, math.Float64bits(n))
	case string:
		c.str(n)
	case *SymbolNode:
		c.compileSymbol(n, sym, ctx)
	case []Node:
		c.compileVector(n, sym, ctx)
	case *ListNode:
		c.compileList(n, sym, ctx)
	default:
		panic(fmt.Sprintf("unsupported node [%v:%T]", node, node))
	}
}

// compileSymbol compiles a symbol. The symbol can refer to a local variable or
// a global definition.
func (c *Compiler) compileSymbol(n *SymbolNode, sym *SymTable, ctx *Ctx) {
	// The symbol is locally bound. Load the bound value from the FRAMES stack.
	if idx, ok := sym.IndexOf(n.Name); ok {
		// fmt.Printf("GET %s @ %d\n", n.Name, idx)
		c.instr(vm.OpGetArg, uint64(idx))
		return
	}
	// The symbol refers to a globally defined value. Load the value from the
	// DEFS stack.
	if id, ok := c.defs.get(n.Name); ok {
		c.instr(vm.OpGetGlobal, id)
		return
	}
	// The symbol is neither a local argument nor a global value.
	panic(fmt.Sprintf("unknown symbol [%s]", n.Name))
}

// compileVector compiles a vector. If it is empty it will compile to a single
// EmptyVector instruction. If there are elements they will be compiled in
// reverse order and bracketed in End and List instructions.
func (c *Compiler) compileVector(n []Node, sym *SymTable, ctx *Ctx) {
	switch len(n) {
	case 0:
		c.instr(vm.OpEmptyVector)
	default:
		c.instr(vm.OpEnd)
		c.compileNodesReverse(n, sym, ctx)
		c.instr(vm.OpList)
	}
}

// compileList compiles a function invocation. An empty list will compile to an
// EmptyVector instruction. If there is at leas a single element in the list
// and the first element of the list is a symbol, it will be matched with the
// available special forms and primitive functions. The corresponding functions
// will be compiled if there is a match. Else a function call will be compiled.
//   If the first element is itself a list we compile ths list beforehand. The
// result of that list call is expected to yield a function reference.
func (c *Compiler) compileList(n *ListNode, sym *SymTable, ctx *Ctx) {
	if n.Empty() {
		c.instr(vm.OpEmptyVector)
	}
	switch x := n.First().(type) {
	case *SymbolNode:
		// Special forms handle their arguments in various ways. The arguments
		// may not get compiled in sequence.
		if spec, ok := c.specs[x.Name]; ok {
			spec(n.Rest(), sym, ctx)
			return
		}
		// Primitive functions follow the same pattern: first compile all the
		// arguments (either ascending or descending), then append a single VM
		// instruction.
		if prim, ok := c.prims[x.Name]; ok {
			c.compilePrim(prim, n.Rest(), sym, ctx)
			return
		}
		// Variable primitive functions are like primitives but with a variable
		// number of arguments.
		if prim, ok := c.varPrims[x.Name]; ok {
			c.compileVarPrim(prim, n.Rest(), sym, ctx)
			return
		}
		// Compile a call to the global function.
		c.compileCall(x, n.Rest(), sym, ctx)
	case *ListNode:
		// Compile the list. We expect the result of the computation to be a
		// REF value which we can then call.
		c.compileListCall(x, n.Rest(), sym, ctx)
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v:%T]", x, x))
	}
}

// compilePrim compiles primitive functions with an exact number of arguments.
func (c *Compiler) compilePrim(prim primDef, args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != prim.nargs {
		panic(fmt.Sprintf("[%s] requires exactly [%d] arguments", prim.name, prim.nargs))
	}
	if prim.rev {
		c.compileNodesReverse(args, sym, ctx)
	} else {
		c.compileNodes(args, sym, ctx)
	}
	c.instr(prim.op)
}

// compileVarPrim compiles primitive functions with a variable number of
// arguments. Optionally ther can be a lower limit on the number of arguments.
func (c *Compiler) compileVarPrim(prim varPrimDef, args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) < prim.argsMin {
		panic(fmt.Sprintf("[%s] requires at least [%d] arguments", prim.name, prim.argsMin))
	}

	c.instr(vm.OpEnd)
	c.compileNodesReverse(args, sym, ctx)
	c.instr(prim.op)
}

func (c *Compiler) compileDebug(args []Node, sym *SymTable, ctx *Ctx) {
	c.instr(vm.OpDebug, uint64(args[0].(int64)))
}

func (c *Compiler) compileSub(args []Node, sym *SymTable, ctx *Ctx) {
	switch len(args) {
	case 0:
		// Empty difference (-) yields 0.
		c.instr(vm.OpConst, 0)
	case 1:
		// Singleton difference (- x) yields (- 0 x) which if effectively -x.
		c.instr(vm.OpConst, 0)
		c.compile(args[0], sym, ctx)
		c.instr(vm.OpSub)
	case 2:
		// Only supports at most two operands and computes their difference.
		//
		//   <(- x1 x2)> :=
		//        <x1>
		//        <x2>
		//        Sub
		//
		c.compile(args[0], sym, ctx)
		c.compile(args[1], sym, ctx)
		c.instr(vm.OpSub)
	default:
		panic("[-] Too many arguments")
	}
}

func (c *Compiler) compileDiv(args []Node, sym *SymTable, ctx *Ctx) {
	switch len(args) {
	case 0:
		// Empty division (/) yields 1.
		c.instr(vm.OpConst, 1)
	case 1:
		// Singleton difference (/ x) yields (/ 1 x) which if effectively 1/x.
		c.instr(vm.OpConst, 1)
		c.compile(args[0], sym, ctx)
		c.instr(vm.OpDiv)
	case 2:
		// Only supports at most two operands and computes their quotient.
		//
		//   <(/ x1 x2)> :=
		//        <x1>
		//        <x2>
		//        Div
		//
		c.compile(args[0], sym, ctx)
		c.compile(args[1], sym, ctx)
		c.instr(vm.OpDiv)
	default:
		panic("[/] Too many arguments")
	}
}

func (c *Compiler) compileAnd(args []Node, sym *SymTable, ctx *Ctx) {
	switch len(args) {
	case 0:
		// Empty and (and) yields true.
		c.instr(vm.OpTrue)
	case 1:
		// Singleton and (and x) yields x.
		c.compile(args[0], sym, ctx)
	default:
		//   <(and x1 x2 ... xn)> :=
		//        <x1>
		//        JumpIfNot L0
		//        <x2>
		//        JumpIfNot L0
		//        ...
		//        <xn>
		//        Jump L1
		//    L0: False
		//    L1: ...
		//
		lbl := c.newLbl()
		end := c.newLbl()
		for i := 0; i < len(args)-1; i++ {
			c.compile(args[i], sym, ctx)
			c.labeled(vm.OpJumpIfNot, lbl)
		}
		c.compile(args[len(args)-1], sym, ctx)
		c.labeled(vm.OpJump, end)
		c.label(lbl)
		c.instr(vm.OpFalse)
		c.label(end)
	}
}

func (c *Compiler) compileOr(args []Node, sym *SymTable, ctx *Ctx) {
	switch len(args) {
	case 0:
		// Empty or (or) yields false.
		c.instr(vm.OpFalse)
	case 1:
		// Singleton and (or x) yields x.
		c.compile(args[0], sym, ctx)
	default:
		lbl := c.newLbl()
		end := c.newLbl()
		for i := 0; i < len(args)-1; i++ {
			c.compile(args[i], sym, ctx)
			c.labeled(vm.OpJumpIf, lbl)
		}
		c.compile(args[len(args)-1], sym, ctx)
		c.labeled(vm.OpJump, end)
		c.label(lbl)
		c.instr(vm.OpTrue)
		c.label(end)
	}
}

func (c *Compiler) compileSet(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 2 {
		panic("[set] requires exactly two arguments")
	}
	s, ok := args[0].(*SymbolNode)
	if !ok {
		panic("[set] requires first argument to be a symbol")
	}

	// Add the local binding to the symbol table. We do this before we compile
	// the body. This permits recursive definitions.
	sym.AddVar(s.Name)
	// n, _ := sym.IndexOf(s.Name)
	// fmt.Printf("SET %s @ %d\n", s.Name, n)

	c.compile(args[1], sym, ctx)
	c.instr(vm.OpPushArgs, 1)
}

func (c *Compiler) compileLet(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 2 {
		panic("[let] requires exactly two arguments")
	}

	bs, ok := args[0].(*ListNode)
	if !ok {
		panic("[let] requires first argument to be a list of bindings")
	}
	if len(bs.Items)%2 == 1 {
		panic("[let] reqires an even number of bindings")
	}

	// TODO: separate symbol table would be better.
	// TODO: Problem when shadowing a variable.
	locals := make([]string, 0)
	for i := 0; i < len(bs.Items); i += 2 {
		s, ok := bs.Items[i].(*SymbolNode)
		if !ok {
			panic(fmt.Sprintf("[let] cannot bind to [%v]", s))
		}
		// Keep track of the let bindings. We will remove them when we fall
		// out of scope.
		locals = append(locals, s.Name)
		// Add the local binding to the symbol table.
		sym.AddVar(s.Name)

		c.compile(bs.Items[i+1], sym, ctx)
		// Add the let bindings one at a time so that subsequent bindings
		// will be able to access the privously defined let bindings.
		c.instr(vm.OpPushArgs, 1)
	}

	c.compile(args[1], sym, ctx)
	// A let binding does not posess a separate frame. We need to drop all
	// introduced let bindings before we continue.
	c.instr(vm.OpDropArgs, uint64(len(locals)))

	// Remove the let bindings from the symbol table as well.
	sym.Remove(locals)
}

// compileDef compiles a global definition. Global definitions will be bound in
// the root environment and are available in the entire codebase for the entire
// lifetime of the program.
//
//    <(def x y)> :=
//        <y>
//        SetGlobal #x
//
func (c *Compiler) compileDef(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 2 {
		panic("[def] requires exactly two arguments")
	}

	s, ok := args[0].(*SymbolNode)
	if !ok {
		panic("[def] requires first argument to be a symbol")
	}
	// Assing a new ID to the definition name. It's required to do this before
	// compiling the body of the definition in order to make the symbol available
	// to recursive function calls.
	id := c.defs.getOrAdd(s.Name)

	c.compile(args[1], sym, ctx)
	c.instr(vm.OpSetGlobal, id)
}

//
//   <(if cond cons)> :=
//       <cond>
//       JumpIfNot L0
//       <cons>
//   L0: ...
//
//   <(if cond cons alt)> :=
//       <cons>
//       JumpIfNot L0
//       <cons>
//       Jump L1
//   L0: <alt>
//   L1: ...
//
func (c *Compiler) compileIf(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 2 && len(args) != 3 {
		panic("[if] requires either two or three arguments")
	}
	switch len(args) {
	case 2:
		end := c.newLbl()
		c.compile(args[0], sym, ctx)
		c.labeled(vm.OpJumpIfNot, end)
		c.compile(args[1], sym, ctx)
		c.label(end)
	case 3:
		alt := c.newLbl()
		end := c.newLbl()
		c.compile(args[0], sym, ctx)
		c.labeled(vm.OpJumpIfNot, alt)
		c.compile(args[1], sym, ctx)
		c.labeled(vm.OpJump, end)
		c.label(alt)
		c.compile(args[2], sym, ctx)
		c.label(end)
	}
}

// compileCond compiles a condition expression. It is expected to have an even
// number of alternating condition and code block arguments.
//
//   <(cond cond1 block1 cond2 block2 ... condN blockN)> :=
//       <cond1>
//       JumpIfNot L0
//       <block1>
//       Jump LX
//   L0: <cond2>
//       JumpIfNot L1
//       <block2>
//       Jump LX
//   L1: <cond3>
//       ...
//   LN: <condN>
//       JumpIfNot LX
//       <blockN>
//   LX: ...
//
func (c *Compiler) compileCond(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args)%2 == 1 {
		panic("[cond] reqires an even number of case-block pairs")
	}
	len := len(args)
	end := c.newLbl()
	for i := 0; i < len-2; i += 2 {
		nxt := c.newLbl()
		c.compile(args[i], sym, ctx)
		c.labeled(vm.OpJumpIfNot, nxt)
		c.compile(args[i+1], sym, ctx)
		c.labeled(vm.OpJump, end)
		c.label(nxt)
	}
	if len >= 2 {
		c.compile(args[len-2], sym, ctx)
		c.labeled(vm.OpJumpIfNot, end)
		c.compile(args[len-1], sym, ctx)
		c.label(end)
	}
}

func (c *Compiler) compileFn(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 2 {
		panic("[fn] expects exactly 2 arguments")
	}
	// Accept parameter lists either as (a1 a2 ..) or as [a1 a2 ...] though it
	// is customary to use the second notational form through out.
	switch x := args[0].(type) {
	case *ListNode:
		c.compileFn2(x.Items, args[1], sym, ctx)
	case []Node:
		c.compileFn2(x, args[1], sym, ctx)
	default:
		panic("[fn] first argument must be a list or vector")
	}
}

func (c *Compiler) compileFn2(params []Node, body Node, sym *SymTable, ctx *Ctx) {
	// fmt.Printf("OLD PARAMS %v\n", params)
	// fmt.Println(sym)
	// Find all closure parameters.
	eps := c.listClosureParamsForSub(params, body, sym)
	// Prepend closure paramesters.
	newParams := []Node{}
	newParams = append(newParams, eps...)
	newParams = append(newParams, params...)

	// fmt.Printf("NEW PARAMS %v\n", newParams)

	sub := sym.NewSymTable()

	skp := c.newLbl()
	fen := c.newLbl()
	// Compiles the function body in-place.
	// Jump over the function implementation.
	c.labeled(vm.OpJump, skp)
	// The function entry point.
	// fmt.Println("BODY")
	c.label(fen)
	// Compile the acutal function code.
	c.compileFnBody(newParams, body, sub, ctx)
	// This marks the end of the function.
	// fmt.Println("BODY END")
	c.label(skp)
	// Push the extern arguments on the stack for the closure.
	// fmt.Println("CLOSURE CALL")
	for _, ep := range eps {
		c.compileSymbol(ep.(*SymbolNode), sym, ctx)
	}
	// fmt.Println("CLOSURE CALL END")
	// Return a closure to the function.
	c.ref(len(eps), fen)
}

func (c *Compiler) compileFnBody(params []Node, body Node, sym *SymTable, ctx *Ctx) {
	switch len(params) {
	case 0:
		// Removes the function argument's end marker from the stack.
		c.instr(vm.OpPop)
		c.compile(body, sym, ctx)
		c.instr(vm.OpReturn)
	default:
		man, opt := c.extractParams(params)
		// Add the mandatory arguments to the local symbol table.
		sym.Add(man)
		// Push the mandatory arguments to the frames stack.
		c.instr(vm.OpPushArgs, uint64(len(man)))
		// Check for an optional argument.
		if opt != "" {
			// Add the optional argument to the local symbol table.
			sym.AddVar(opt)
			// The LIST operation will append all remaining arguments to a vector.
			c.instr(vm.OpList)
			// Then push the vector to the frames stack as well.
			c.instr(vm.OpPushArgs, 1)
		} else {
			// Removes the function argument's end marker from the stack.
			c.instr(vm.OpPop)
		}
		// Compile the body with this closure context.
		c.compile(body, sym, ctx)
		c.instr(vm.OpReturn)
	}
}

func (c *Compiler) listClosureParamsForSub(params []Node, node Node, sym *SymTable) []Node {
	sub := sym.NewSymTable()
	man, opt := c.extractParams(params)
	// fmt.Printf("PARAMS %v, %s\n", man, opt)
	// Add the mandatory arguments to the local symbol table.
	sub.Add(man)
	// Check for an optional argument.
	if opt != "" {
		// Add the optional argument to the local symbol table.
		sub.AddVar(opt)
	}

	// fmt.Println(sub)

	return c.listClosureParams(node, sub)
}

func (c *Compiler) listClosureParams(node Node, sym *SymTable) []Node {
	switch n := node.(type) {
	case *SymbolNode:
		idx, ok := sym.IndexOf(n.Name)
		// fmt.Printf("POT SYM %s @ %d\n", n.Name, idx)
		if ok && idx < 0 {
			return []Node{n}
		}
		// Check if it is a local symbol that is external to the current scope.
		// if idx, ok := sym.IndexOf(n.Name); ok && idx < 0 {
		// 	fmt.Printf("CLOS SYM %s\n", n.Name)
		// 	return []Node{n}
		// }
	case []Node:
		return c.listClosureParamsList(n, sym)
	case *ListNode:
		return c.listClosureParamsList(n.Items, sym)
	}
	return []Node{}
}

func (c *Compiler) listClosureParamsList(nodes []Node, sym *SymTable) []Node {
	res := []Node{}
	for _, node := range nodes {
		for _, param := range c.listClosureParams(node, sym) {
			if notContainsSymbol(res, param.(*SymbolNode)) {
				res = append(res, param)
			}
		}
	}
	return res
}

func notContainsSymbol(syms []Node, s *SymbolNode) bool {
	for _, sym := range syms {
		if symx, ok := sym.(*SymbolNode); ok {
			if symx.Name == s.Name {
				return false
			}
		}
	}
	return true
}

func (c *Compiler) extractParams(params []Node) ([]string, string) {
	names := c.verifyParams(params)
	pos := util.IndexOf(names, "&")
	if pos == -1 {
		return names, ""
	}
	if len(names) == pos+1 {
		panic("[fn] missing optional parameter")
	}
	if len(names) > pos+2 {
		panic("[fn] excess optional parameter")
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

// compileRec compiles a manual recursive call. In fact it will only set the
// Recursive flag to true. Subsequent compilations of call and list call will
// consider the set flag and compile a recursive call instead of a regular one.
func (c *Compiler) compileRec(args []Node, sym *SymTable, ctx *Ctx) {
	if len(args) != 1 {
		panic("[rec] expects exactly 1 argument")
	}

	if x, ok := args[0].(*ListNode); ok {
		c.compileList(x, sym, ctx.NewRecCtx(true))
	} else {
		panic("[rec] argument must be a list call")
	}
}

func (c *Compiler) compileCall(s *SymbolNode, args []Node, sym *SymTable, ctx *Ctx) {
	c.instr(vm.OpEnd)
	// Reset recursive invocation for all arguments.
	c.compileNodesReverse(args, sym, ctx.NewRecCtx(false))
	// The calling function can be implemented either as a global definintion
	// or passed as a local argument from a let binding.
	c.compileSymbol(s, sym, ctx)
	c.compileCallByCtx(ctx)
}

func (c *Compiler) compileListCall(lst *ListNode, args []Node, sym *SymTable, ctx *Ctx) {
	c.instr(vm.OpEnd)
	c.compileNodesReverse(args, sym, ctx.NewRecCtx(false))
	c.compileList(lst, sym, ctx)
	c.compileCallByCtx(ctx)
}

// compileCallByCtx compiles either a Call or a RecCall depending on the
// Recursive flag in the context object.
func (c *Compiler) compileCallByCtx(ctx *Ctx) {
	if ctx.Recurse {
		c.instr(vm.OpRecCall)
	} else {
		c.instr(vm.OpCall)
	}
}

func (c *Compiler) compileNodes(nodes []Node, sym *SymTable, ctx *Ctx) {
	for _, node := range nodes {
		c.compile(node, sym, ctx)
	}
}

func (c *Compiler) compileNodesReverse(nodes []Node, sym *SymTable, ctx *Ctx) {
	for i := len(nodes) - 1; i >= 0; i-- {
		c.compile(nodes[i], sym, ctx)
	}
}

func (c *Compiler) label(name string) {
	c.code = append(c.code, asm.Label(name))
}

func (c *Compiler) labeled(op vm.Op, name string) {
	c.code = append(c.code, asm.Labeled(op, name))
}

func (c *Compiler) ref(cargs int, name string) {
	c.code = append(c.code, asm.Ref(cargs, name))
}

func (c *Compiler) instr(op vm.Op, args ...uint64) {
	c.code = append(c.code, asm.Instr(op, args...))
}

func (c *Compiler) str(str string) {
	c.code = append(c.code, asm.Str(str))
}

func (c *Compiler) newLbl() string {
	lbl := fmt.Sprintf("L%d", c.lblId)
	c.lblId++
	return lbl
}
