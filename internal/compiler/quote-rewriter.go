package compiler

import (
	"fmt"
	"reflect"
)

type quoteRewriter struct {
}

func NewQuoteRewriter() *quoteRewriter {
	return &quoteRewriter{}
}

func (r *quoteRewriter) Rewrite(n Node) Node {
	_, m := r.rewrite(n)
	return m
}

func (r *quoteRewriter) rewrite(n Node) ([]Node, Node) {
	switch x := n.(type) {
	case bool:
		return r.rewriteBoolean(x)
	case int64:
		return r.rewriteInteger(x)
	case string:
		return r.rewriteString(x)
	case *SymbolNode:
		return r.rewriteSymbol(x)
	case *VectorNode:
		return r.rewriteVector(x)
	case *ListNode:
		return r.rewriteList(x)
	}
	panic(fmt.Sprintf("Quote-Rewriter: Unsupported node [%v:%T]", n, n))
}

func (r *quoteRewriter) rewriteBoolean(n bool) ([]Node, Node) {
	return r.empty(), n
}

func (r *quoteRewriter) rewriteInteger(n int64) ([]Node, Node) {
	return r.empty(), n
}

func (r *quoteRewriter) rewriteString(n string) ([]Node, Node) {
	return r.empty(), n
}

func (r *quoteRewriter) rewriteSymbol(n *SymbolNode) ([]Node, Node) {
	return r.empty(), n
}

func (r *quoteRewriter) rewriteVector(n *VectorNode) ([]Node, Node) {
	return r.empty(), n
}

func (r *quoteRewriter) rewriteList(n *ListNode) ([]Node, Node) {
	syms := r.empty()
	if len(n.Items) == 0 {
		return syms, n
	}
	switch x := n.Items[0].(type) {
	case *SymbolNode:
		switch x.Name {
		case "quote":
			ss, m := r.rewrite(n.Items[1])
			return r.empty(), Fn(ss, m)
		case "unquote":
			switch y := n.Items[1].(type) {
			case *SymbolNode:
				syms = append(syms, y)
				return syms, y
			case *ListNode:
				switch z := y.Items[0].(type) {
				case *SymbolNode:
					switch z.Name {
					case "dissolve":
						syms = append(syms, y.Items[1])
						return syms, y
					}
				}
			}
		}
	}
	ss, ms := r.rewriteItems(n.Items)
	return ss, NewList(ms)
}

func (r *quoteRewriter) rewriteItems(ns []Node) ([]Node, []Node) {
	ss := r.empty()
	ms := r.empty()
	for _, n := range ns {
		s, m := r.rewrite(n)
		ss = join(ss, s)
		ms = append(ms, m)
	}
	return ss, ms
}

func (r *quoteRewriter) empty() []Node {
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
