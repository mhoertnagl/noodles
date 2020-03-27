package compiler

type Rewriter interface {
	Rewrite(n Node) Node
}

func RewriteItems(r Rewriter, ns []Node) []Node {
	ms := []Node{}
	for _, n := range ns {
		m := r.Rewrite(n)
		if m != nil {
			ms = append(ms, m)
		}
	}
	return ms
}
