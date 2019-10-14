package eval

import (
	"fmt"
	"io/ioutil"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
)

type CoreFun func(Evaluator, data.Env, []data.Node) data.Node

type Evaluator interface {
	Eval(node data.Node) data.Node
	Error(format string, args ...interface{}) data.Node
	Errors() []*data.ErrorNode
	AddCoreFun(name string, fun CoreFun)
}

type evaluator struct {
	env     data.Env
	core    map[string]CoreFun
	err     []*data.ErrorNode
	reader  read.Reader
	parser  read.Parser
	printer print.Printer
}

// TODO: PrintErrors

func NewEvaluator(env data.Env) Evaluator {
	e := &evaluator{
		env:     env,
		core:    make(map[string]CoreFun),
		err:     []*data.ErrorNode{},
		reader:  read.NewReader(),
		parser:  read.NewParser(),
		printer: print.NewPrinter(),
	}
	InitCore(e)
	// TODO: Import prelude.
	// TODO: partial evaluation.
	// TODO: rest delimiter | e.g. (fun foobar [x | xs] )
	// TODO: ' (quote) and ` (quasiquote)
	return e
}

func (e *evaluator) debug(format string, env data.Env, args ...data.Node) {
	strs := make([]interface{}, len(args))
	for i, arg := range args {
		strs[i] = e.printer.Print(arg)
	}
	fmt.Printf(format, strs...)
}

func (e *evaluator) Error(format string, args ...interface{}) data.Node {
	err := data.NewError(fmt.Sprintf(format, args...))
	e.err = append(e.err, err)
	return err
}

func (e *evaluator) Errors() []*data.ErrorNode {
	return e.err
}

func (e *evaluator) AddCoreFun(name string, fun CoreFun) {
	e.core[name] = fun
}

func (e *evaluator) Eval(node data.Node) data.Node {
	return e.eval(e.env, node)
}

func (e *evaluator) eval(env data.Env, n data.Node) data.Node {
	for {
		switch x := n.(type) {
		case *data.ListNode:
			if len(x.Items) == 0 {
				return n
			}

			hd := x.Items[0]
			args := x.Items[1:]

			if sym, ok1 := hd.(*data.SymbolNode); ok1 {
				switch sym.Name {
				case "def!":
					return e.evalDef(env, args)
				case "let*":
					env, n = e.evalLet(env, args)
					continue
				case "fn*":
					return e.evalFunDef(env, args)
				case "do":
					n = e.evalDo(env, args)
					continue
				case "if":
					n = e.evalIf(env, args)
					continue
				// TODO: TCO?
				case "parse":
					return e.evalParse(env, args)
				// TODO: TCO?
				case "eval":
					return e.evalEval(env, args)
				case "read-file":
					return e.evalReadFile(env, args)
				case "quote":
					return e.evalQuote(env, args)
				case "quasiquote":
					n = e.evalQuasiquote(env, args)
					continue
				case "defmacro!":
					return e.evalDefMacro(env, args)
				default:
					if fun, ok := e.core[sym.Name]; ok {
						args = e.evalSeq(env, args)
						return fun(e, env, args)
					}
				}
			}
			// e.debug("%s\n", env, hd)
			hd = e.eval(env, hd)

			if fn, ok2 := hd.(*data.FuncNode); ok2 {
				if len(fn.Pars) != len(args) {
					return e.Error("Number of arguments not the same as number of parameters.")
				}
				// Create a new environment for this function.
				fn.Env = data.NewEnv(fn.Env)
				// Evaluate and bind argurments to their parameters in the new function
				// environment.
				for i, par := range fn.Pars {
					arg := e.eval(env, args[i])
					fn.Env.Set(par, arg)
				}
				env = fn.Env
				n = fn.Fun
				continue
			}

			return e.Error("List cannot be evaluated.")
		case *data.SymbolNode:
			return e.evalSymbol(env, x)
		case *data.VectorNode:
			return e.evalVector(env, x)
		case *data.HashMapNode:
			return e.evalHashMap(env, x)
		default:
			// Return unchanged. These are immutable atoms.
			return n
		}
	}
}

// evalSymbol returns the symbol if it is a core function. If not, serches the
// current environment (and all parent environments) and returns its value.
// Returns an error if no such element exists.
func (e *evaluator) evalSymbol(env data.Env, n *data.SymbolNode) data.Node {
	// TODO: core functions should be defined in the environment.
	// First check if the symbol defines a core function. Return the symbol
	// unchanged if this is true.
	if _, ok1 := e.core[n.Name]; ok1 {
		return n
	}
	// See if a value is bound to the symbol is defined in the environment.
	// Return the value.
	if v, ok2 := env.Lookup(n.Name); ok2 {
		return v
	}
	// Else the symbol is undefined. Report an error.
	return e.Error("Undefined symbol [%s].", n.Name)
}

// evalVector evaluates all arguments of the vector and returns a new vector
// with the results in the same order as the original arguments.
func (e *evaluator) evalVector(env data.Env, n *data.VectorNode) data.Node {
	// Return the original vector if it is empty.
	// TODO: Always return a new vector?
	if len(n.Items) == 0 {
		return n
	}
	items := e.evalSeq(env, n.Items)
	return data.NewVector(items)
}

func (e *evaluator) evalHashMap(env data.Env, n *data.HashMapNode) data.Node {
	// Return the original hash map if it is empty.
	// TODO: Always return a new map?
	if len(n.Items) == 0 {
		return n
	}
	c := data.NewEmptyHashMap()
	for key, val := range n.Items {
		k := e.eval(env, key)
		v := e.eval(env, val)
		if sk, ok := k.(string); ok {
			c.Items[sk] = v
		} else {
			e.Error("HashMap key must be string.")
		}
	}
	return c
}

// evalSeq evaluates a list of nodes sequentially and returns a list of
// evaluated nodes in the original order.
func (e *evaluator) evalSeq(env data.Env, items []data.Node) []data.Node {
	res := make([]data.Node, len(items))
	for i, item := range items {
		res[i] = e.eval(env, item)
	}
	return res
}

// TODO: Should def! be able to overwrite an already defined binding?
// evalDef binds a name to a value. Evaluates the value before it get bound to
// the name and returns it. Redefinitions of the same name in the same
// environment will overwrite the previous value.
func (e *evaluator) evalDef(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("[def!] requires exactly 2 arguments.")
	}
	return e.evalSet(env, ns[0], ns[1])
}

// evalSet evaluates the name and the val argument and binds name to val in the
// environment.
func (e *evaluator) evalSet(env data.Env, name data.Node, val data.Node) data.Node {
	v := e.eval(env, val)
	if x, ok := name.(*data.SymbolNode); ok {
		return env.Set(x.Name, v)
	}
	// TODO: Perhaps we can evaluate a node if it is not a symbol node.
	// n := e.eval(env, name)
	// TODO: StringNode. We should append an obscure unicode character to the
	// string to make it different from other symbols. Or we add "".
	// This would make debugging easier.
	return e.Error("Cannot bind to [%s].", name)
}

// evalLet binds a list, vector or hash-map of pairs to a new local environment
// and evaluates it's body with respect to this new environment.
// If the second argument is neiher a list, vector or hash-map this function
// yields a runtime error. The list of arguments has to be a sequence of name-
// value pairs. The values will be evaluated before they get bound to their
// respective names.
func (e *evaluator) evalLet(env data.Env, ns []data.Node) (data.Env, data.Node) {
	if len(ns) != 2 {
		return env, e.Error("[let*] requires exactly 2 arguments.")
	}
	sub := data.NewEnv(env)
	switch b := ns[0].(type) {
	case *data.ListNode:
		e.evalSeqBindings(sub, b.Items)
	case *data.VectorNode:
		e.evalSeqBindings(sub, b.Items)
	case *data.HashMapNode:
		fmt.Printf("%v\n", b.Items)
		e.evalHashMapBindings(sub, b.Items)
	default:
		return env, e.Error("Cannot [let*]-bind non-sequence.")
	}
	// Return the new environment and the unevaluated body of let* for TCO.
	return sub, ns[1]
}

func (e *evaluator) evalSeqBindings(env data.Env, b []data.Node) {
	for i := 0; i < len(b); i += 2 {
		e.evalSet(env, b[i], b[i+1])
	}
}

func (e *evaluator) evalHashMapBindings(env data.Env, b data.Map) {
	fmt.Printf("%v\n", b)
	for k, v := range b {
		e.evalSet(env, k, v)
	}
}

// evalFunDef creates a new Function Node with references to the current
// environment, the lsit of parameter names and the function body.
func (e *evaluator) evalFunDef(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("[fn*] requires exactly 2 arguments.")
	}
	switch ps := ns[0].(type) {
	case *data.ListNode:
		if params, ok2 := paramNames(ps.Items); ok2 {
			return data.NewFuncNode(env, params, ns[1])
		}
		return e.Error("Function parameter must be a symbol.")
	case *data.VectorNode:
		if params, ok2 := paramNames(ps.Items); ok2 {
			return data.NewFuncNode(env, params, ns[1])
		}
		return e.Error("Function parameter must be a symbol.")
	default:
		return e.Error("First argument to [fn*] must be a list or vector.")
	}
}

// paramNames assumes that all nodes in [ns] are symbol nodes and returns the
// string list of the names of these symbols preserving order.
func paramNames(ns []data.Node) ([]string, bool) {
	params := make([]string, len(ns))
	for i, n := range ns {
		if p, ok := n.(*data.SymbolNode); ok {
			params[i] = p.Name
		} else {
			return nil, false
		}
	}
	return params, true
}

// evalDo evaluates a list of items and returns the final evaluated result.
// Returns nil when the list is empty.
func (e *evaluator) evalDo(env data.Env, ns []data.Node) data.Node {
	z := len(ns) - 1
	if z <= 0 {
		return nil
	}
	for _, n := range ns[:z] {
		e.eval(env, n)
	}
	// Return the last item unevaluated for TCO.
	return ns[z]
}

// evalIf evaluates its first argument. If it is true?? evaluates the second
// argument else evaluates the third argument. If the first argument is not true
// and no third argument is given then it returns nil.
func (e *evaluator) evalIf(env data.Env, ns []data.Node) data.Node {
	len := len(ns)
	if len != 2 && len != 3 {
		return e.Error("[if] requires either 2 or 3 arguments.")
	}
	// Evaluate the condition.
	cond := e.eval(env, ns[0])
	// Return immediatly if the condition evaluated to an error.
	if data.IsError(cond) {
		return cond
	}
	if e.isTrue(env, cond) {
		// Return unevaluated true branch for TCO.
		return ns[1]
	} else if len == 3 {
		// Return unevaluated false branch for TCO.
		return ns[2]
	}
	return nil
}

// TODO: Create a core function. We need it for (true? ...) and (false? ...)
// isTrue returns false when the node is nil, false, 0, "", (), [] and {}.
// It returns true in all remaining cases.
func (e *evaluator) isTrue(env data.Env, n data.Node) bool {
	switch x := n.(type) {
	case nil:
		return false
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		return len(x) != 0
	case *data.SymbolNode:
		// TODO: Add a unit test. What if the symbol is not defined?
		v := e.evalSymbol(env, x)
		return e.isTrue(env, v)
	case *data.ListNode:
		return len(x.Items) != 0
	case *data.VectorNode:
		return len(x.Items) != 0
	case *data.HashMapNode:
		return len(x.Items) != 0
	default:
		return true
	}
}

// evalParse parses the string input and returns the parsed AST.
func (e *evaluator) evalParse(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("[parse] requires exactly 1 argument.")
	}
	n := e.eval(env, ns[0])
	if s, ok := n.(string); ok {
		e.reader.Load(s)
		return e.parser.Parse(e.reader)
	}
	return e.Error("[parse] argument must be a string.")
}

// evalEval evaluates the first argument twice. This is useful in combination
// with [parse] or [quote]. For instance consider (eval (read "(+ 1 1)")). The
// evaluation of the argument yields (eval (+ 1 1)) and evaluation again gives
// 2 as expected. On the other hand an expression like (eval (+ 1 1)) where the
// argument can be evaluated at once is equivalent to (+ 1 1).
// Evaluates the node in the local environment.
func (e *evaluator) evalEval(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("[eval] requires exactly 1 argument.")
	}
	n := e.eval(env, ns[0])
	return e.eval(env, n)
}

// evalReadFile reads the contents of a file into a string and returns the
// result. The only argument to this function is the file path string. Returns
// an error if the file could not be found or read.
func (e *evaluator) evalReadFile(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("[read-file] requires exactly 1 argument.")
	}
	n := e.eval(env, ns[0])
	if s, ok := n.(string); ok {
		f, err := ioutil.ReadFile(s)
		if err != nil {
			return e.Error("[read-file] %s", err)
		}
		return string(f)
	}
	return e.Error("[read-file] argument must be a string.")
}

// evalQuote returns its only argument unevaluated.
func (e *evaluator) evalQuote(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("[quote] requires exactly 1 argument.")
	}
	return ns[0]
}

// evalQuasiquote is a quoting operations where one can specify unquoted or
// splice-unquoted holes in the quoted expression.
// An unquoted subexpression will be evaluated first and the result will
// replace the original expression in the quoted expression.
// A splice-unquoted list will be evaluated and the results will be spliced
// into the quoted parent list.
func (e *evaluator) evalQuasiquote(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("[quasiquote] requires exactly 1 argument.")
	}
	return e.quasiquote1(env, ns[0])
}

func (e *evaluator) quasiquote1(env data.Env, n data.Node) data.Node {
	if x, ok := n.(*data.ListNode); ok && len(x.Items) > 0 {
		return e.quasiquoteList(env, x)
	}
	// If it is not a list, quote the element.
	return data.Quote(n)
}

func (e *evaluator) quasiquoteList(env data.Env, n *data.ListNode) data.Node {
	switch y := n.Items[0].(type) {
	case *data.SymbolNode:
		if y.Name == "unquote" {
			// Return the element without an enclosing quote.
			return n.Items[1]
		}
	case *data.ListNode:
		if len(y.Items) == 0 {
			return data.Quote(y)
		}
		if z, ok := y.Items[0].(*data.SymbolNode); ok {
			if z.Name == "splice-unquote" {
				return data.Concat(
					y.Items[1],
					e.quasiquote1(env, data.NewList(n.Items[1:])),
				)
			}
		}
	}
	return data.Cons(
		e.quasiquote1(env, n.Items[0]),
		e.quasiquote1(env, data.NewList(n.Items[1:])),
	)
}

// evalDefMacro binds the second argument to the first one and declares the
// second argument a macro. The first argument has to be a symbol and the
// second argument a function.
func (e *evaluator) evalDefMacro(env data.Env, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("[macrodef!] requires exactly 2 arguments.")
	}
	if sym, ok1 := ns[0].(*data.SymbolNode); ok1 {
		v := e.eval(env, ns[1])
		if fun, ok2 := v.(*data.FuncNode); ok2 {
			// Delcare this function a macro.
			fun.IsMacro = true
			return env.Set(sym.Name, fun)
		}
		return e.Error("[macrodef!] requires second argument to be a function.")
	}
	return e.Error("[macrodef!] Cannot bind to [%s].", ns[0])
}
