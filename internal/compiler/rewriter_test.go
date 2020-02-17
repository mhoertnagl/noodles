package compiler_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/compiler"
)

func TestRewriteBoolean(t *testing.T) {
	testcrw(t, `true`, `true`)
}

func TestRewriteInteger(t *testing.T) {
	testcrw(t, `1`, `1`)
}

func TestRewriteSymbol(t *testing.T) {
	testcrw(t, `a`, `a`)
}

func TestRewriteVector(t *testing.T) {
	testcrw(t, `[1 2 3]`, `[1 2 3]`)
}

func TestRewriteList(t *testing.T) {
	testcrw(t, `(1 2 3)`, `(1 2 3)`)
}

func TestRewriteSimpleQuote(t *testing.T) {
	testcrw(t, `'(+ 1 1)`, `(fn [] (+ 1 1))`)
}

func TestRewriteReplacementQuote(t *testing.T) {
	testcrw(t, `'(+ ~a ~b)`, `(fn [a b] (+ a b))`)
}

func TestRewriteSpliceQuote(t *testing.T) {
	testcrw(t, `'(+ ~a ~@b)`, `(fn [a b] (+ a @b))`)
}

func testcrw(t *testing.T, i string, e string) {
	t.Helper()
	r := compiler.NewReader()
	p := compiler.NewParser()
	r.Load(i)
	in := p.Parse(r)
	r.Load(e)
	en := p.Parse(r)
	testrw(t, in, en)
}

func testrw(t *testing.T, i compiler.Node, e compiler.Node) {
	t.Helper()
	qr := compiler.NewQuoteRewriter()
	a := qr.Rewrite(i)
	pr := compiler.NewPrinter()
	as := pr.Print(a)
	es := pr.Print(e)
	if equalNode(a, e) == false {
		t.Errorf("Mismatch Expecting \n  [%s]\n but got \n  [%s].", es, as)
	}
}

func equalNode(l compiler.Node, r compiler.Node) bool {
	switch lx := l.(type) {
	case bool:
		if rx, ok := r.(bool); ok {
			return lx == rx
		}
	case int64:
		if rx, ok := r.(int64); ok {
			return lx == rx
		}
	case *compiler.SymbolNode:
		if rx, ok := r.(*compiler.SymbolNode); ok {
			return lx.Name == rx.Name
		}
	case *compiler.VectorNode:
		if rx, ok := r.(*compiler.VectorNode); ok {
			return equalList(lx.Items, rx.Items)
		}
	case *compiler.ListNode:
		if rx, ok := r.(*compiler.ListNode); ok {
			return equalList(lx.Items, rx.Items)
		}
	}
	return false
}

func equalList(l []compiler.Node, r []compiler.Node) bool {
	if len(l) != len(r) {
		return false
	}
	for i := 0; i < len(l); i++ {
		if equalNode(l[i], r[i]) == false {
			return false
		}
	}
	return true
}
