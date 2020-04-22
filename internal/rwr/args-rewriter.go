package rwr

import (
	"fmt"

	"github.com/mhoertnagl/splis2/internal/cmp"
)

type argsMap map[string]cmp.Node

type ArgsRewriter struct {
	ams argsMap
}

func NewArgsRewriter(man []string, opt string, args []cmp.Node) *ArgsRewriter {
	ams := argsMap{}

	for i := 0; i < len(man); i++ {
		switch y := args[i].(type) {
		case []cmp.Node:
			ams[man[i]] = singletonList(args[i])
		default:
			ams[man[i]] = y
		}
		// ams[man[i]] = args[i]
	}

	if opt != "" {
		ams[opt] = args[len(man):]
	}

	return &ArgsRewriter{ams: ams}
}

func (r *ArgsRewriter) Rewrite(n cmp.Node) cmp.Node {
	return rewrite(r, n)[0]
	// switch x := n.(type) {
	// case *cmp.SymbolNode:
	// 	if a, ok := r.ams[x.Name]; ok {
	// 		return a
	// 	}
	// 	return x
	// case []cmp.Node:
	// 	return RewriteItems(r, x)
	// case *cmp.ListNode:
	//
	// 	if cmp.IsCall(x, "dissolve") {
	// 		return cmp.Do(r.Rewrite(x.Items[1]))
	// 	}
	// 	return cmp.NewList(RewriteItems(r, x.Items))
	// default:
	// 	return n
	// }
}

func rewrite(r *ArgsRewriter, n cmp.Node) []cmp.Node {
	switch x := n.(type) {
	case *cmp.SymbolNode:
		if a, ok := r.ams[x.Name]; ok {
			switch y := a.(type) {
			case []cmp.Node:
				return y
			default:
				return singletonList(a)
			}
		}
		return singletonList(x)
	case []cmp.Node:
		fmt.Printf("%v\n", rewriteItems(r, x))
		fmt.Printf("%v\n", singletonList(rewriteItems(r, x)))
		return singletonList(rewriteItems(r, x))
	case *cmp.ListNode:
		if cmp.IsCall(x, "dissolve") {
			return rewrite(r, x.Items[1])
		}
		return singletonList(cmp.NewList(rewriteItems(r, x.Items)))
	default:
		return singletonList(n)
	}
}

func singletonList(n cmp.Node) []cmp.Node {
	return []cmp.Node{n}
}

func rewriteItems(r *ArgsRewriter, ns []cmp.Node) []cmp.Node {
	ms := []cmp.Node{}
	for _, n := range ns {
		m := rewrite(r, n)
		if m != nil {
			ms = append(ms, m...)
		}
	}
	return ms
}
