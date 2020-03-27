package compiler

import (
	"reflect"
)

type QuoteRewriter struct {
}

func NewQuoteRewriter() *QuoteRewriter {
	return &QuoteRewriter{}
}

func (r *QuoteRewriter) Rewrite(n Node) Node {
	_, m := r.rewrite(n)
	return m
}

func (r *QuoteRewriter) rewrite(n Node) ([]Node, Node) {
	switch x := n.(type) {
	case *VectorNode:
		return r.rewriteVector(x)
	case *ListNode:
		return r.rewriteList(x)
	default:
		return r.empty(), n
	}
}

func (r *QuoteRewriter) rewriteVector(n *VectorNode) ([]Node, Node) {
	s, m := r.rewriteItems(n.Items)
	return s, NewVector(m)
}

func (r *QuoteRewriter) rewriteList(n *ListNode) ([]Node, Node) {
	syms := r.empty()
	if len(n.Items) == 0 {
		return syms, n
	}
	if IsCall(n, "quote") {
		ss, m := r.rewrite(n.Items[1])
		return r.empty(), Fn(ss, m)
	}
	if IsCall(n, "unquote") {
		switch y := n.Items[1].(type) {
		case *SymbolNode:
			syms = append(syms, y)
			return syms, y
		case *ListNode:
			if IsCall(y, "dissolve") {
				syms = append(syms, y.Items[1])
				return syms, y
			}
		}
	}
	ss, ms := r.rewriteItems(n.Items)
	return ss, NewList(ms)
}

func (r *QuoteRewriter) rewriteItems(ns []Node) ([]Node, []Node) {
	ss := r.empty()
	ms := r.empty()
	for _, n := range ns {
		s, m := r.rewrite(n)
		ss = join(ss, s)
		ms = append(ms, m)
	}
	return ss, ms
}

func (r *QuoteRewriter) empty() []Node {
	return make([]Node, 0)
}

func join(a []Node, b []Node) []Node {
	for _, x := range b {
		if !contains(a, x) {
			a = append(a, x)
		}
	}
	return a
}

func contains(a []Node, x Node) bool {
	for _, y := range a {
		if reflect.DeepEqual(y, x) {
			return true
		}
	}
	return false
}
