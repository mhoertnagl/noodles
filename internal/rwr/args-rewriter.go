package rwr

import "github.com/mhoertnagl/splis2/internal/cmp"

// TODO: Macro tests mit compiler in rewriter tests verschieben.

type argsMap map[string]cmp.Node

type ArgsRewriter struct {
	ams argsMap
}

func NewArgsRewriter(pars []string, args []cmp.Node) *ArgsRewriter {
	ams := argsMap{}
	for i := 0; i < len(pars); i++ {
		ams[pars[i]] = args[i]
	}
	return &ArgsRewriter{
		ams: ams,
	}
}

func (r *ArgsRewriter) Rewrite(n cmp.Node) cmp.Node {
	switch x := n.(type) {
	case *cmp.SymbolNode:
		if a, ok := r.ams[x.Name]; ok {
			return a
		}
		return x
	case []cmp.Node:
		return RewriteItems(r, x)
	case *cmp.ListNode:
		return cmp.NewList(RewriteItems(r, x.Items))
	default:
		return n
	}
}
