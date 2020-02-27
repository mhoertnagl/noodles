package compiler

import "fmt"

type argsMap map[string]Node

func NewArgsMap(pars []string, args []Node) argsMap {
	ams := argsMap{}
	for i := 0; i < len(pars); i++ {
		ams[pars[i]] = args[i]
	}
	return ams
}

func NewArgsRewriter(pars []string, args []Node) *argsRewriter {
	return &argsRewriter{
		ams: NewArgsMap(pars, args),
	}
}

type argsRewriter struct {
	ams argsMap
}

func (r *argsRewriter) Rewrite(n Node) Node {
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
	panic(fmt.Sprintf("Args-Rewriter: Unsupported node [%v:%T]", n, n))
}

func (r *argsRewriter) rewriteBoolean(n bool) Node {
	return n
}

func (r *argsRewriter) rewriteInteger(n int64) Node {
	return n
}

func (r *argsRewriter) rewriteString(n string) Node {
	return n
}

func (r *argsRewriter) rewriteSymbol(n *SymbolNode) Node {
	if a, ok := r.ams[n.Name]; ok {
		return a
	}
	return n
}

func (r *argsRewriter) rewriteVector(n *VectorNode) Node {
	return n
}

func (r *argsRewriter) rewriteList(n *ListNode) Node {
	if len(n.Items) == 0 {
		return n
	}
	return NewList(r.rewriteItems(n.Items))
}

func (r *argsRewriter) rewriteItems(ns []Node) []Node {
	ms := make([]Node, len(ns))
	for i, n := range ns {
		ms[i] = r.Rewrite(n)
	}
	return ms
}
