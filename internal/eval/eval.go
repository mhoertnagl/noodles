package eval

import (
	"fmt"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/print"
)

type CoreFun func(Evaluator, data.Env, []data.Node) data.Node

type Evaluator interface {
	Eval(node data.Node) data.Node
	EvalEnv(env data.Env, n data.Node) data.Node
	Error(format string, args ...interface{}) data.Node
	Errors() []*data.ErrorNode
	AddCoreFun(name string, fun CoreFun)
}

type evaluator struct {
	env     data.Env
	core    map[string]CoreFun
	err     []*data.ErrorNode
	printer print.Printer
}

// TODO: PrintErrors

func NewEvaluator(env data.Env) Evaluator {
	e := &evaluator{
		env:     env,
		core:    make(map[string]CoreFun),
		err:     []*data.ErrorNode{},
		printer: print.NewPrinter(),
	}
	InitCore(e)
	return e
}

func (e *evaluator) debug(format string, args ...data.Node) {
	strs := make([]interface{}, len(args))
	for i, arg := range args {
		strs[i] = e.printer.Print(arg)
	}
	fmt.Printf(format, strs...)
}

func (e *evaluator) Eval(node data.Node) data.Node {
	return e.EvalEnv(e.env, node)
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

func (e *evaluator) findCoreFun(name string) (CoreFun, bool) {
	if fun, ok := e.core[name]; ok {
		return fun, true
	}
	return nil, false
}

// TODO: Can be private.
func (e *evaluator) EvalEnv(env data.Env, n data.Node) data.Node {
	switch {
	case data.IsList(n):
		return e.evalList(env, n.(*data.ListNode))
	case data.IsVector(n):
		return e.evalVector(env, n.(*data.VectorNode))
	case data.IsHashMap(n):
		return e.evalHashMap(env, n.(*data.HashMapNode))
	case data.IsSymbol(n):
		return e.evalSymbol(env, n.(*data.SymbolNode))
	default:
		// Return unchanged. These are immutable atoms.
		return n
	}
}

func (e *evaluator) evalList(env data.Env, n *data.ListNode) data.Node {
	if len(n.Items) == 0 {
		return n
	}

	hd := n.Items[0]
	args := n.Items[1:]

	if sym, ok1 := hd.(*data.SymbolNode); ok1 {
		switch sym.Name {
		case "def!":
			return e.evalDef(env, sym.Name, args)
		case "let*":
			return e.evalLet(env, sym.Name, args)
		case "fn*":
			return e.evalFunDef(env, sym.Name, args)
		case "do":
			return e.evalDo(env, sym.Name, args)
		case "if":
			return e.evalIf(env, sym.Name, args)
		default:
			if fun, ok := e.findCoreFun(sym.Name); ok {
				args = e.evalSeq(env, args)
				return fun(e, env, args)
			}
		}
	}

	hd = e.EvalEnv(env, hd)

	if fn, ok2 := hd.(*data.FuncNode); ok2 {
		if len(fn.Pars) != len(args) {
			return e.Error("Number of arguments not the same as number of parameters.")
		}
		// Create a new environment for this function.
		fn.Env = data.NewEnv(fn.Env)
		// Evaluate and bind argurments to their parameters in the new function
		// environment.
		for i, par := range fn.Pars {
			arg := e.EvalEnv(env, args[i])
			fn.Env.Set(par, arg)
		}
		return e.EvalEnv(fn.Env, fn.Fun)
	}
	return e.Error("List cannot be evaluated.")
}

func (e *evaluator) evalVector(env data.Env, n *data.VectorNode) data.Node {
	if len(n.Items) == 0 {
		return n
	}
	items := e.evalSeq(env, n.Items)
	return data.NewVector(items)
}

func (e *evaluator) evalHashMap(env data.Env, n *data.HashMapNode) data.Node {
	c := data.NewEmptyHashMap()
	for key, val := range n.Items {
		k := e.EvalEnv(env, key)
		v := e.EvalEnv(env, val)
		if sk, ok := k.(string); ok {
			c.Items[sk] = v
		} else {
			e.Error("HashMap key must be string.")
		}
	}
	return c
}

func (e *evaluator) evalSeq(env data.Env, items []data.Node) []data.Node {
	res := make([]data.Node, len(items))
	for i, item := range items {
		res[i] = e.EvalEnv(env, item)
	}
	return res
}

// evalSymbol searches for a symbol in the environment (and all parent
// environments) and returns its value.
// Returns an error when no such element exists.
func (e *evaluator) evalSymbol(env data.Env, n *data.SymbolNode) data.Node {
	if v, ok := env.Lookup(n.Name); ok {
		return v
	}
	return e.Error("Undefined symbol [%s].", n.Name)
}

// TODO: Should def! be able to overwrite an already defined binding?
// evalDef binds a name to a value. Evaluates the value before it get bound to
// the name and returns it. Redefinitions of the same name in the same
// environment will overwrite the previous value.
func (e *evaluator) evalDef(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("def! requires exactly 2 arguments.")
	}
	return e.evalSet(env, ns[0], ns[1])
}

// evalSet evaluates the name and the val argument and binds name to val in the
// environment.
func (e *evaluator) evalSet(env data.Env, name data.Node, val data.Node) data.Node {
	v := e.EvalEnv(env, val)
	if x, ok := name.(*data.SymbolNode); ok {
		return env.Set(x.Name, v)
	}
	// Perhaps we can evaluate a node if it is not a symbol node.
	//n := e.EvalEnv(env, name)
	// TODO: StringNode. We should append an obscure unicode character to the string to make it different from other symbols.
	// Or we add "". This would make debugging easier.
	return e.Error("Cannot bind to [%s].", name)
}

// evalLet binds a list, vector or hash-map of pairs to a new local environment
// and evaluates it's body with respect to this new environment.
// If the second argument is neiher a list, vector or hash-map this function
// yields a runtime error. The list of arguments has to be a sequence of name-
// value pairs. The values will be evaluated before they get bound to their
// respective names.
func (e *evaluator) evalLet(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("let* requires exactly 2 arguments.")
	}
	sub := data.NewEnv(env)
	switch b := ns[0].(type) {
	case *data.ListNode:
		e.evalSeqBindings(sub, b.Items)
	case *data.VectorNode:
		e.evalSeqBindings(sub, b.Items)
	case *data.HashMapNode:
		e.evalHashMapBindings(sub, b.Items)
	default:
		return e.Error("Cannot let*-bind non-sequence.")
	}
	// Evaluate the body with the new local environment.
	return e.EvalEnv(sub, ns[1])
}

func (e *evaluator) evalSeqBindings(env data.Env, b []data.Node) {
	for i := 0; i < len(b); i += 2 {
		e.evalSet(env, b[i], b[i+1])
	}
}

func (e *evaluator) evalHashMapBindings(env data.Env, b data.Map) {
	for k, v := range b {
		e.evalSet(env, k, v)
	}
}

func (e *evaluator) evalFunDef(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.Error("fn* requires exactly 2 arguments.")
	}
	if ps, ok := ns[0].(*data.ListNode); ok {
		params := make([]string, len(ps.Items))
		for i, param := range ps.Items {
			if p, ok2 := param.(*data.SymbolNode); ok2 {
				params[i] = p.Name
			} else {
				return e.Error("Function parameter must be a symbol.")
			}
		}
		return data.NewFuncNode(env, params, ns[1])
	}
	return e.Error("First argument to fn* must be a list.")
}

// evalDo evaluates a list of items and returns the final evaluated result.
// Returns nil when the list is empty.
func (e *evaluator) evalDo(env data.Env, name string, ns []data.Node) data.Node {
	var r data.Node
	for _, n := range ns {
		r = e.EvalEnv(env, n)
	}
	return r
}

// evalIf evaluates its first argument. If it is true?? evaluates the second
// argument else evaluates the third argument. If the first argument is not true
// and no third argument is given then it returns nil.
func (e *evaluator) evalIf(env data.Env, name string, ns []data.Node) data.Node {
	len := len(ns)
	if len != 2 && len != 3 {
		return e.Error("if requires either 2 or 3 arguments.")
	}
	cond := e.EvalEnv(env, ns[0])
	if data.IsError(cond) {
		return cond
	}
	if isTrue(env, cond) {
		return e.EvalEnv(env, ns[1])
	} else if len == 3 {
		return e.EvalEnv(env, ns[2])
	}
	return nil
}

func isTrue(env data.Env, n data.Node) bool {
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
		return isTrue(env, x)
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
