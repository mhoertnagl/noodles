package eval

import (
	"fmt"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/print"
)

type Evaluator interface {
	Eval(node data.Node) data.Node
	Errors() []*data.ErrorNode
}

type evaluator struct {
	env     data.Env
	err     []*data.ErrorNode
	printer print.Printer
}

// TODO: PrintErrors

func NewEvaluator(env data.Env) Evaluator {
	e := &evaluator{
		env:     env,
		err:     []*data.ErrorNode{},
		printer: print.NewPrinter(),
	}
	env.AddSpecialForm("def!", e.evalDef)
	env.AddSpecialForm("let*", e.evalLet)
	env.AddSpecialForm("fn*", e.evalFunDef)
	env.AddSpecialForm("do", e.evalDo)
	env.AddSpecialForm("if", e.evalIf)
	env.AddSpecialForm("+", e.evalSum)
	env.AddSpecialForm("-", e.evalDiff)
	env.AddSpecialForm("<", e.eval2f(func(n0 float64, n1 float64) data.Node { return n0 < n1 }))
	env.AddSpecialForm(">", e.eval2f(func(n0 float64, n1 float64) data.Node { return n0 > n1 }))
	env.AddSpecialForm("<=", e.eval2f(func(n0 float64, n1 float64) data.Node { return n0 <= n1 }))
	env.AddSpecialForm(">=", e.eval2f(func(n0 float64, n1 float64) data.Node { return n0 >= n1 }))
	return e
}

func (e *evaluator) error(format string, args ...interface{}) data.Node {
	err := data.NewError(fmt.Sprintf(format, args...))
	e.err = append(e.err, err)
	return err
}

func (e *evaluator) debug(format string, args ...data.Node) {
	strs := make([]interface{}, len(args))
	for i, arg := range args {
		strs[i] = e.printer.Print(arg)
	}
	fmt.Printf(format, strs...)
}

func (e *evaluator) Eval(node data.Node) data.Node {
	return e.eval(e.env, node)
}

func (e *evaluator) Errors() []*data.ErrorNode {
	return e.err
}

func (e *evaluator) eval(env data.Env, n data.Node) data.Node {
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

	if data.IsSymbol(hd) {
		fn := hd.(*data.SymbolNode).Name
		// TODO: switch of builtins.
		if fun, ok := env.FindSpecialForm(fn); ok {
			// Special Forms get their arguments passed unevaluated. They usually
			// have custom evaluation strategies.
			return fun(env, fn, args)
		}
	}

	// Evaluate the head of the list.
	hd = e.eval(env, hd)

	if data.IsFuncNode(hd) {
		fn := hd.(*data.FuncNode)
		if len(fn.Pars) != len(args) {
			return e.error("Number of arguments not the same as number of parameters.")
		}
		e.debug("Fun call: ")
		for i, par := range fn.Pars {
			arg := e.eval(env, args[i])
			fmt.Printf(par)
			e.debug("=%s ", arg)
			fn.Env.Set(par, arg)
		}
		e.debug("\n")
		return e.eval(fn.Env, fn.Fun)
	}
	return e.error("List cannot be evaluated.")
}

func (e *evaluator) evalVector(env data.Env, n *data.VectorNode) data.Node {
	if len(n.Items) == 0 {
		return n
	}
	items := e.evalSeq(env, n.Items)
	return data.NewVector(items)
}

func (e *evaluator) evalHashMap(env data.Env, n *data.HashMapNode) data.Node {
	c := data.NewHashMap2()
	// TODO: to separate func?
	for key, val := range n.Items {
		k := e.eval(env, key)
		v := e.eval(env, val)
		c.Items[k] = v
	}
	return c
}

func (e *evaluator) evalSeq(env data.Env, items []data.Node) []data.Node {
	res := make([]data.Node, len(items))
	for i, item := range items {
		res[i] = e.eval(env, item)
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
	return e.error("Undefined symbol [%s].", n.Name)
}

// TODO: Should def! be able to overwrite an already defined binding?
// evalDef binds a name to a value. Evaluates the value before it get bound to
// the name and returns it. Redefinitions of the same name in the same
// environment will overwrite the previous value.
func (e *evaluator) evalDef(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.error("def! requires exactly 2 arguments.")
	}
	return e.evalSet(env, ns[0], ns[1])
}

func (e *evaluator) evalFunDef(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.error("fn* requires exactly 2 arguments.")
	}
	if as, ok := ns[0].(*data.ListNode); ok {
		args := make([]string, len(as.Items))
		for i, arg := range as.Items {
			if a, ok2 := arg.(*data.SymbolNode); ok2 {
				args[i] = a.Name
			} else {
				return e.error("Function parameter must be a symbol.")
			}
		}
		return data.NewFuncNode(data.NewEnv(env), args, ns[1])
	}
	return e.error("First argument to fn* must be a list.")
}

// evalLet binds a list, vector or hash-map of pairs to a new local environment
// and evaluates it's body with respect to this new environment.
// If the second argument is neiher a list, vector or hash-map this function
// yields a runtime error. The list of arguments has to be a sequence of name-
// value pairs. The values will be evaluated before they get bound to their
// respective names.
func (e *evaluator) evalLet(env data.Env, name string, ns []data.Node) data.Node {
	if len(ns) != 2 {
		return e.error("let* requires exactly 2 arguments.")
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
		return e.error("Cannot let*-bind non-sequence.")
	}
	// Evaluate the body with the new local environment.
	return e.eval(sub, ns[1])
}

// evalSum computes the sum of all arguments.
// Non-numeric arguments will be ignored.
func (e *evaluator) evalSum(env data.Env, name string, ns []data.Node) data.Node {
	var sum float64
	for _, n := range ns {
		m := e.eval(env, n)
		if v, ok := m.(float64); ok {
			sum += v
		} else {
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.error("[%s] is not a number.", "")
		}
	}
	return sum
}

func (e *evaluator) evalDiff(env data.Env, name string, ns []data.Node) data.Node {
	len := len(ns)
	if len == 1 {
		v := e.eval(env, ns[0])
		if n, ok := v.(float64); ok {
			return -n
		} else {
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.error("[%s] is not a number.", "")
		}
	}
	if len == 2 {
		v1 := e.eval(env, ns[0])
		v2 := e.eval(env, ns[1])
		if n1, ok1 := v1.(float64); ok1 {
			if n2, ok2 := v2.(float64); ok2 {
				return n1 - n2
			}
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.error("[%s] is not a number.", "")
		}
		// TODO: Add a printer instance to the evaluator to print expressions.
		return e.error("[%s] is not a number.", "")
	}
	return e.error("- requires either 1 or 2 arguments.")
}

// evalDo evaluates a list of items and returns the final evaluated result.
// Returns nil when the list is empty.
func (e *evaluator) evalDo(env data.Env, name string, ns []data.Node) data.Node {
	var r data.Node
	for _, n := range ns {
		r = e.eval(env, n)
	}
	return r
}

// evalIf evaluates its first argument. If it is true?? evaluates the second
// argument else evaluates the third argument. If the first argument is not true
// and no third argument is given then it returns nil.
func (e *evaluator) evalIf(env data.Env, name string, ns []data.Node) data.Node {
	len := len(ns)
	if len != 2 && len != 3 {
		return e.error("if requires either 2 or 3 arguments.")
	}
	cond := e.eval(env, ns[0])
	if data.IsError(cond) {
		return cond
	}
	e.debug("Condition: %s\n", cond)
	if e.evalTrue(env, cond) {
		e.debug("True\n")
		return e.eval(env, ns[1])
	} else if len == 3 {
		return e.eval(env, ns[2])
	}
	return nil
}

// TODO: evalAnd, evalOr

func (e *evaluator) evalTrue(env data.Env, n data.Node) bool {
	switch {
	// case data.IsError(n):
	// 	return false
	case data.IsNil(n):
		return false
	case data.IsBool(n):
		return n.(bool)
	case data.IsNumber(n):
		return n.(float64) != 0
	case data.IsString(n):
		return len(n.(string)) != 0
	case data.IsList(n):
		return len(n.(*data.ListNode).Items) != 0
	case data.IsVector(n):
		return len(n.(*data.VectorNode).Items) != 0
	case data.IsHashMap(n):
		return len(n.(*data.HashMapNode).Items) != 0
	default:
		return true
	}
}

func (e *evaluator) eval2f(f func(float64, float64) data.Node) data.SpecialForm {
	return func(env data.Env, name string, ns []data.Node) data.Node {
		if len(ns) != 2 {
			return e.error("%s requires exactly 2 numeric arguments.", name)
		}
		v0 := e.eval(env, ns[0])
		v1 := e.eval(env, ns[1])
		if n0, ok0 := v0.(float64); ok0 {
			if n1, ok1 := v1.(float64); ok1 {
				return f(n0, n1)
			}
		}
		// TODO: Angeben welches argument kein float ist.
		return e.error("")
	}
}

func (e *evaluator) evalSeqBindings(env data.Env, b []data.Node) {
	for i := 0; i < len(b); i += 2 {
		e.evalSet(env, b[i], b[i+1])
	}
}

func (e *evaluator) evalHashMapBindings(env data.Env, b map[data.Node]data.Node) {
	for k, v := range b {
		e.evalSet(env, k, v)
	}
}

// evalSet evaluates the name and the val argument and binds name to val in the
// environment.
func (e *evaluator) evalSet(env data.Env, name data.Node, val data.Node) data.Node {
	v := e.eval(env, val)
	if x, ok := name.(*data.SymbolNode); ok {
		return env.Set(x.Name, v)
	}
	// Perhaps we can evaluate a node if it is not a symbol node.
	//n := e.eval(env, name)
	// TODO: StringNode. We should append an obscure unicode character to the string to make it different from other symbols.
	// Or we add "". This would make debugging easier.
	return e.error("Cannot bind to [%s].", name)
}
