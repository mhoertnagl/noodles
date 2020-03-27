package compiler

type usingsSet = map[string]struct{}

type UseRewriter struct {
	usings usingsSet
	rdr    Reader
	prs    Parser
}

// TODO: pass dir paths to files.
func NewUseRewriter() *UseRewriter {
	return &UseRewriter{
		usings: usingsSet{},
		rdr:    NewReader(),
		prs:    NewParser(),
	}
}

func (r *UseRewriter) Rewrite(n Node) Node {
	switch x := n.(type) {
	case *VectorNode:
		return NewVector(RewriteItems(r, x.Items))
	case *ListNode:
		return r.rewriteList(x)
	default:
		return n
	}
}

func (r *UseRewriter) rewriteList(n *ListNode) Node {
	if len(n.Items) == 0 {
		return n
	}
	// TODO: Length of items should be 2 (use "...")
	if IsCall(n, "use") {
		if s, ok := n.Items[1].(string); ok {
			if _, ok2 := r.usings[s]; ok2 {
				// File has already been included. Skip.
				return nil
			}
			// TODO: Create file path.
			// TODO: Read file.
			r.rdr.Load("...")
			c := r.prs.Parse(r.rdr)
			return r.Rewrite(c)
		}
	}
	return NewList(RewriteItems(r, n.Items))
}
