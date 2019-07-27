package eval

import (
	"github.com/mhoertnagl/splis2/internal/read"
)

type Evaluator interface {
	Eval(env Env, node read.Node) read.Node
}

type evaluator struct {
}

func NewEvaluator() Evaluator {
	return &evaluator{}
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
	default:
		// Return unchanged. These are immutable atoms.
		return n
	}
}

func (e *evaluator) evalList(env Env, n *read.ListNode) read.Node {
	return &read.ListNode{Items: e.evalSeq(env, n.Items)}
}

func (e *evaluator) evalVector(env Env, n *read.VectorNode) read.Node {
	return &read.VectorNode{Items: e.evalSeq(env, n.Items)}
}

func (e *evaluator) evalHashMap(env Env, n *read.HashMapNode) read.Node {
	c := &read.HashMapNode{Items: make(map[read.Node]read.Node)}
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
