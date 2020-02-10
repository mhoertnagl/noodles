package compiler

import "fmt"

type Rewriter interface {
	Rewrite(n Node) Node
	// RewriteBoolean(n bool) Node
	// RewriteInteger(n int64) Node
	// RewriteSymbol(n *SymbolNode) Node
	// RewriteVector(n *VectorNode) Node
	// RewriteList(n *ListNode) Node
}

func NewQuoteRewriter() *quoteRewriter {
	return &quoteRewriter{}
}

type quoteRewriter struct {
}

func (r *quoteRewriter) Rewrite(n Node) Node {
	switch x := n.(type) {
	case bool:
		return r.rewriteBoolean(x)
	case int64:
		return r.rewriteInteger(x)
	case *SymbolNode:
		return r.rewriteSymbol(x)
	case *VectorNode:
		return r.rewriteVector(x)
	case *ListNode:
		return r.rewriteList(x)
	}
	panic(fmt.Sprintf("Unsupported node [%v]", n))
}

func (r *quoteRewriter) rewriteBoolean(n bool) Node {
	return n
}

func (r *quoteRewriter) rewriteInteger(n int64) Node {
	return n
}

func (r *quoteRewriter) rewriteSymbol(n *SymbolNode) Node {
	return n
}

func (r *quoteRewriter) rewriteVector(n *VectorNode) Node {
	return n
}

func (r *quoteRewriter) rewriteList(n *ListNode) Node {
	if len(n.Items) == 0 {
		return n
	}
	items := n.Items
	args := items[1:]
	switch x := items[0].(type) {
	case *SymbolNode:
		switch x.Name {
		case "quote":
			return NewList2(NewSymbol("fn"), NewVector2(), args[0])
		default:
			return n
		}
	default:
		return n
	}
}
