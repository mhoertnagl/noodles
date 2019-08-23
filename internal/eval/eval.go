package eval

import (
	"fmt"
	"github.com/mhoertnagl/splis2/internal/read"
)

type Evaluator interface {
	Eval(node read.Node) read.Node
}

type evaluator struct {
	env Env
	err []*read.ErrorNode
}

// TODO: PrintErrors

func NewEvaluator(env Env) Evaluator {
	e := &evaluator{env: env}
	env.AddSpecialForm("def!", e.evalDef)
	env.AddSpecialForm("let*", e.evalLet)
	env.AddSpecialForm("do", e.evalDo)
	env.AddSpecialForm("if", e.evalIf)
	env.AddSpecialForm("+", e.evalSum)
	env.AddSpecialForm("<", e.eval2f(func(n0 float64, n1 float64) read.Node { return n0 < n1 }))
	env.AddSpecialForm(">", e.eval2f(func(n0 float64, n1 float64) read.Node { return n0 > n1 }))
	env.AddSpecialForm("<=", e.eval2f(func(n0 float64, n1 float64) read.Node { return n0 <= n1 }))
	env.AddSpecialForm(">=", e.eval2f(func(n0 float64, n1 float64) read.Node { return n0 >= n1 }))
	return e
}

func (e *evaluator) error(format string, args ...interface{}) read.Node {
	err := read.NewError(fmt.Sprintf(format, args...))
	e.err = append(e.err, err)
	return err
}

func (e *evaluator) Eval(node read.Node) read.Node {
	return e.eval(e.env, node)
}

func (e *evaluator) eval(env Env, node read.Node) read.Node {
	switch n := node.(type) {
	case *read.ListNode:
		return e.evalList(env, n)
	case *read.VectorNode:
		return e.evalVector(env, n)
	case *read.HashMapNode:
		return e.evalHashMap(env, n)
	case *read.SymbolNode:
		return e.evalSymbol(env, n)
	default:
		// Return unchanged. These are immutable atoms.
		return n
	}
}

func (e *evaluator) evalList(env Env, n *read.ListNode) read.Node {
	if len(n.Items) == 0 {
		return n
	}

	switch x := n.Items[0].(type) {
	case *read.SymbolNode:
		if fun, ok := env.FindSpecialForm(x.Name); ok {
			// Special Forms get their arguments passed unevaluated. They usually
			// have custom evaluation strategies.
			return fun(env, x.Name, n.Items[1:])
		}
		// TODO: else if // user defined fun.
	}

	// Evaluate all items of the list and return a new list with the evaluated
	// items.
	items := e.evalSeq(env, n.Items)
	return read.NewList(items)
}

func (e *evaluator) evalVector(env Env, n *read.VectorNode) read.Node {
	items := e.evalSeq(env, n.Items)
	return read.NewVector(items)
}

func (e *evaluator) evalHashMap(env Env, n *read.HashMapNode) read.Node {
	c := read.NewHashMap2()
	// TODO: to separate func?
	for key, val := range n.Items {
		k := e.eval(env, key)
		v := e.eval(env, val)
		c.Items[k] = v
	}
	return c
}

func (e *evaluator) evalSeq(env Env, items []read.Node) []read.Node {
	res := make([]read.Node, len(items))
	for i, item := range items {
		res[i] = e.eval(env, item)
	}
	return res
}

// evalSymbol searches for a symbol in the environment (and all parent
// environments) and returns its value.
// Returns an error when no such element exists.
func (e *evaluator) evalSymbol(env Env, n *read.SymbolNode) read.Node {
	if v, ok := env.Lookup(n.Name); ok {
		return v
	}
	return e.error("Undefined symbol [%s].", n.Name)
}

// TODO: Should def! be able to overwrite an already defined binding?
// evalDef binds a name to a value. Evaluates the value before it get bound to
// the name and returns it. Redefinitions of the same name in the same
// environment will overwrite the previous value.
func (e *evaluator) evalDef(env Env, name string, ns []read.Node) read.Node {
	if len(ns) != 2 {
		return e.error("def! requires exactly 2 arguments.")
	}
	return e.evalSet(env, ns[0], ns[1])
}

// evalLet binds a list, vector or hash-map of pairs to a new local environment
// and evaluates it's body with respect to this new environment.
// If the second argument is neiher a list, vector or hash-map this function
// yields a runtime error. The list of arguments has to be a sequence of name-
// value pairs. The values will be evaluated before they get bound to their
// respective names.
func (e *evaluator) evalLet(env Env, name string, ns []read.Node) read.Node {
	if len(ns) != 2 {
		return e.error("let* requires exactly 2 arguments.")
	}
	sub := NewEnv(env)
	switch b := ns[0].(type) {
	case *read.ListNode:
		e.evalSeqBindings(sub, b.Items)
	case *read.VectorNode:
		e.evalSeqBindings(sub, b.Items)
	case *read.HashMapNode:
		e.evalHashMapBindings(sub, b.Items)
	default:
		return e.error("Cannot let*-bind non-sequence.")
	}
	// Evaluate the body with the new local environment.
	return e.eval(sub, ns[1])
}

// evalSum computes the sum of all arguments.
// Non-numeric arguments will be ignored.
func (e *evaluator) evalSum(env Env, name string, ns []read.Node) read.Node {
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

// evalDo evaluates a list of items and returns the final evaluated result.
// Returns nil when the list is empty.
func (e *evaluator) evalDo(env Env, name string, ns []read.Node) read.Node {
	var r read.Node
	for _, n := range ns {
		r = e.eval(env, n)
	}
	return r
}

// evalIf evaluates its first argument. If it is true?? evaluates the second
// argument else evaluates the third argument. If the first argument is not true
// and no third argument is given then it returns nil.
func (e *evaluator) evalIf(env Env, name string, ns []read.Node) read.Node {
	len := len(ns)
	if len != 2 && len != 3 {
		return e.error("if requires either 2 or 3 arguments.")
	}
	cond := e.eval(env, ns[0])
	if e.evalTrue(env, cond) {
		return e.eval(env, ns[1])
	} else if len == 3 {
		return e.eval(env, ns[2])
	}
	return nil
}

// TODO: evalAnd, evalOr

func (e *evaluator) evalTrue(env Env, n read.Node) bool {
	switch {
	case read.IsError(n):
		return false
	case read.IsNil(n):
		return false
	case read.IsBool(n):
		return n.(bool)
	case read.IsNumber(n):
		return n.(float64) != 0
	case read.IsString(n):
		return len(n.(string)) != 0
	case read.IsList(n):
		return len(n.(*read.ListNode).Items) != 0
	case read.IsVector(n):
		return len(n.(*read.VectorNode).Items) != 0
	case read.IsHashMap(n):
		return len(n.(*read.HashMapNode).Items) != 0
	default:
		return true
	}
}

func (e *evaluator) eval2f(f func(float64, float64) read.Node) SpecialForm {
	return func(env Env, name string, ns []read.Node) read.Node {
		if len(ns) != 2 {
			return e.error("%s requires exactly 2 numeric arguments.", name)
		}
		if n0, ok0 := ns[0].(float64); ok0 {
			if n1, ok1 := ns[1].(float64); ok1 {
				return f(n0, n1)
			}
		}
		// TODO: Angeben welches argument kein float ist.
		return e.error("")
	}
}

func (e *evaluator) evalSeqBindings(env Env, b []read.Node) {
	for i := 0; i < len(b); i += 2 {
		e.evalSet(env, b[i], b[i+1])
	}
}

func (e *evaluator) evalHashMapBindings(env Env, b map[read.Node]read.Node) {
	for k, v := range b {
		e.evalSet(env, k, v)
	}
}

// evalSet evaluates the name and the val argument and binds name to val in the
// environment.
func (e *evaluator) evalSet(env Env, name read.Node, val read.Node) read.Node {
	v := e.eval(env, val)
	if x, ok := name.(*read.SymbolNode); ok {
		return env.Set(x.Name, v)
	}
	// Perhaps we can evaluate a node if it is not a symbol node.
	//n := e.eval(env, name)
	// TODO: StringNode. We should append an obscure unicode character to the string to make it different from other symbols.
	// Or we add "". This would make debugging easier.
	return e.error("Cannot bind to [%s].", name)
}
