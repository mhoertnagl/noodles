package eval

import (
  "fmt"
	"github.com/mhoertnagl/splis2/internal/read"
)

type Evaluator interface {
	Eval(env Env, node read.Node) read.Node
}

type evaluator struct {
  err []*read.ErrorNode
}

func NewEvaluator() Evaluator {
	return &evaluator{}
}

func (e *evaluator) error(format string, args ...interface{}) read.Node {
	err := &read.ErrorNode{Msg: fmt.Sprintf(format, args...)}
	e.err = append(e.err, err)
	return err
}

func (e *evaluator) Eval(env Env, node read.Node) read.Node {
	return e.eval(env, node)
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
  // TODO: type SpecialForm func(Env, []read.Node)read.Node
  // TODO: map[string]SpecialForm
  
	return &read.ListNode{Items: e.evalSeq(env, n.Items)}
}

func (e *evaluator) evalVector(env Env, n *read.VectorNode) read.Node {
	return &read.VectorNode{Items: e.evalSeq(env, n.Items)}
}

func (e *evaluator) evalHashMap(env Env, n *read.HashMapNode) read.Node {
	c := &read.HashMapNode{Items: make(map[read.Node]read.Node)}
  // TODO: to separate func?
	for key, val := range n.Items {
		k := e.eval(env, key)
		v := e.eval(env, val)
		c.Items[k] = v
	}
	return c
}

func (e *evaluator) evalSeq(env Env, items []read.Node) []read.Node {
	res := []read.Node{}
	for _, item := range items {
		res = append(res, e.eval(env, item))
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
func (e *evaluator) evalDef(env Env, n *read.ListNode) read.Node {
  return e.evalSet(env, n.Items[1], n.Items[2])
}

// evalSet evaluates the name and the val argument and binds name to val in the 
// environment.
func (e *evaluator) evalSet(env Env, name read.Node, val read.Node) read.Node {
  n := e.eval(env, name)
  v := e.eval(env, val)
  switch x := n.(type) {
  case *read.SymbolNode:
    env.Set(x.Name, v)
    return v
  default:
    return e.error("Cannot bind to [%s].", name)
  }
}

// evalLet binds a list, vector or hash-map of pairs to a noe local environment
// and evaluates it's body in it.
// If the second argument is neiher a list, vector or hash-map this it yields a 
// runtime error.
func (e *evaluator) evalLet(env Env, n *read.ListNode) read.Node {
  sub := NewEnv(env)
  bindings := e.eval(env, n.Items[1])
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
  return e.eval(sub, n.Items[2])
}

func (e *evaluator) evalSeqBindings(env Env, b []read.Node) {
  for i := 0; i < len(b); i+=2 {
    e.evalSet(env, b[i], b[i+1]) 
  }  
}

func (e *evaluator) evalHashMapBindings(env Env, b map[read.Node]read.Node) {
  for k, v := range b {
    e.evalSet(env, k, v)  
  }  
}
