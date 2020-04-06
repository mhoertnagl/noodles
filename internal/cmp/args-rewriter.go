package cmp

type argsMap map[string]Node

type ArgsRewriter struct {
	ams argsMap
}

func NewArgsRewriter(pars []string, args []Node) *ArgsRewriter {
	ams := argsMap{}
	for i := 0; i < len(pars); i++ {
		ams[pars[i]] = args[i]
	}
	return &ArgsRewriter{
		ams: ams,
	}
}

func (r *ArgsRewriter) Rewrite(n Node) Node {
	switch x := n.(type) {
	case *SymbolNode:
		if a, ok := r.ams[x.Name]; ok {
			return a
		}
		return x
	case []Node:
		return RewriteItems(r, x)
	case *ListNode:
		return NewList(RewriteItems(r, x.Items))
	default:
		return n
	}
}
