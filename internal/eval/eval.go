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
			// Evaluate the function arguments and apply function.
			args := e.evalSeq(env, n.Items[1:])
			return fun(env, args)
		}
	}

	// Evaluate all items of the list and return a new list with the evaluated items.
	items := e.evalSeq(env, n.Items)
	return read.NewList(items)
}

func (e *evaluator) evalVector(env Env, n *read.VectorNode) read.Node {
	return read.NewVector(e.evalSeq(env, n.Items))
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

func (e *evaluator) evalSymbol(env Env, n *read.SymbolNode) read.Node {
	if v := env.Lookup(n.Name); v != nil {
		return v
	}
	return e.error("Undefined variable [%s].", n.Name)
}

// evalDef binds a name to a value. Redefinitions of the same name in the same
// environment will be ignored silently.
// (def! a 42) will bind a to 42 in the current environment. Returns the bound
// value 42.
func (e *evaluator) evalDef(env Env, ns []read.Node) read.Node {
	return e.evalSet(env, ns[0], ns[1])
}

// evalLet binds a list, vector or hash-map of pairs to a noe local environment
// and evaluates it's body in it.
// If the second argument is neiher a list, vector or hash-map this it yields a
// runtime error.
func (e *evaluator) evalLet(env Env, ns []read.Node) read.Node {
	sub := NewEnv(env)
	bindings := e.eval(env, ns[0])
	switch b := bindings.(type) {
	case *read.ListNode:
		e.evalSeqBindings(sub, b.Items)
	case *read.VectorNode:
		e.evalSeqBindings(sub, b.Items)
	case *read.HashMapNode:
		e.evalHashMapBindings(sub, b.Items)
	default:
		return e.error("Cannot bind non-sequence.")
	}
	return e.eval(sub, ns[1])
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
	// TODO: Evaluating a symbol node is not a good idea.
	// Perhaps we can evaluate a node if it is not a symbol node.
	//n := e.eval(env, name)
	v := e.eval(env, val)
	switch x := name.(type) {
	case *read.SymbolNode:
		env.Set(x.Name, v)
		return v
	default:
		return e.error("Cannot bind to [%s].", name)
	}
}
