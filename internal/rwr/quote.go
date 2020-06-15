package rwr

import (
	"reflect"

	"github.com/mhoertnagl/noodles/internal/cmp"
)

type QuoteRewriter struct {
}

func NewQuoteRewriter() *QuoteRewriter {
	return &QuoteRewriter{}
}

func (r *QuoteRewriter) Rewrite(n cmp.Node) cmp.Node {
	_, m := r.rewrite(n)
	return m
}

func (r *QuoteRewriter) rewrite(n cmp.Node) ([]cmp.Node, cmp.Node) {
	switch x := n.(type) {
	case []cmp.Node:
		return r.rewriteItems(x)
	case *cmp.ListNode:
		return r.rewriteList(x)
	default:
		return r.empty(), n
	}
}

func (r *QuoteRewriter) rewriteList(n *cmp.ListNode) ([]cmp.Node, cmp.Node) {
	syms := r.empty()
	if len(n.Items) == 0 {
		return syms, n
	}
	if cmp.IsCall(n, "quote") {
		ss, m := r.rewrite(n.Items[1])
		return r.empty(), cmp.Fn(ss, m)
	}
	if cmp.IsCall(n, "unquote") {
		switch y := n.Items[1].(type) {
		case *cmp.SymbolNode:
			syms = append(syms, y)
			return syms, y
		case *cmp.ListNode:
			if cmp.IsCall(y, "dissolve") {
				syms = append(syms, y.Items[1])
				return syms, y
			}
		}
	}
	ss, ms := r.rewriteItems(n.Items)
	return ss, cmp.NewList(ms)
}

func (r *QuoteRewriter) rewriteItems(ns []cmp.Node) ([]cmp.Node, []cmp.Node) {
	ss := r.empty()
	ms := r.empty()
	for _, n := range ns {
		s, m := r.rewrite(n)
		ss = join(ss, s)
		ms = append(ms, m)
	}
	return ss, ms
}

func (r *QuoteRewriter) empty() []cmp.Node {
	return make([]cmp.Node, 0)
}

func join(a []cmp.Node, b []cmp.Node) []cmp.Node {
	for _, x := range b {
		if !contains(a, x) {
			a = append(a, x)
		}
	}
	return a
}

func contains(a []cmp.Node, x cmp.Node) bool {
	for _, y := range a {
		if reflect.DeepEqual(y, x) {
			return true
		}
	}
	return false
}
