package rwr

import "github.com/mhoertnagl/splis2/internal/cmp"

type Rewriter interface {
	Rewrite(n cmp.Node) cmp.Node
}

func RewriteItems(r Rewriter, ns []cmp.Node) []cmp.Node {
	ms := []cmp.Node{}
	for _, n := range ns {
		m := r.Rewrite(n)
		if m != nil {
			ms = append(ms, m)
		}
	}
	return ms
}
