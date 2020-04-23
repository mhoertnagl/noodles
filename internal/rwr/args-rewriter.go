package rwr

import (
	"github.com/mhoertnagl/splis2/internal/cmp"
)

type argsMap map[string]cmp.Node

type ArgsRewriter struct {
	ams argsMap
}

func NewArgsRewriter(man []string, opt string, args []cmp.Node) *ArgsRewriter {
	ams := argsMap{}

	for i := 0; i < len(man); i++ {
		ams[man[i]] = args[i]
	}

	if opt != "" {
		ams[opt] = args[len(man):]
	}

	return &ArgsRewriter{ams: ams}
}

func (r *ArgsRewriter) Rewrite(n cmp.Node) cmp.Node {
	switch x := n.(type) {
	case *cmp.SymbolNode:
		return r.rewriteSymbol(x)
	case []cmp.Node:
		return RewriteItems(r, x)
	case *cmp.ListNode:
		return r.rewriteList(x)
	default:
		return n
	}
}

func (r *ArgsRewriter) rewriteSymbol(n *cmp.SymbolNode) cmp.Node {
	if a, ok := r.ams[n.Name]; ok {
		return a
	}
	return n
}

func (r *ArgsRewriter) rewriteList(n *cmp.ListNode) cmp.Node {
	l := make([]cmp.Node, 0)
	for _, a := range n.Items {
		if x, ok := a.(*cmp.ListNode); ok {
			l = append(l, r.rewriteListArg(x)...)
		} else {
			l = append(l, r.Rewrite(a))
		}
	}
	return cmp.NewList(l)
}

func (r *ArgsRewriter) rewriteListArg(a *cmp.ListNode) []cmp.Node {
	if cmp.IsCall(a, "dissolve") {
		k := r.Rewrite(a.Items[1])
		if kl, ok := k.([]cmp.Node); ok {
			return RewriteItems(r, kl)
		}
	}
	return []cmp.Node{r.Rewrite(a)}
}
