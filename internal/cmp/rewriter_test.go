package cmp_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/cmp"
	"github.com/mhoertnagl/splis2/internal/util"
)

func TestRewriteBoolean(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `true`, `true`)
}

func TestRewriteInteger(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `1`, `1`)
}

func TestRewriteSymbol(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `a`, `a`)
}

func TestRewriteVector(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `[1 2 3]`, `[1 2 3]`)
}

func TestRewriteList(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `(1 2 3)`, `(1 2 3)`)
}

func TestRewriteSimpleQuote(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ 1 1)`, `(fn [] (+ 1 1))`)
}

func TestRewriteReplacementQuote(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ ~a ~b)`, `(fn [a b] (+ a b))`)
}

func TestRewriteSpliceQuote(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ ~a ~@b)`, `(fn [a b] (+ a @b))`)
}

func TestRewriteNestedQuote(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw,
		`'(+ ~a ('(+ ~a ~b) a 1))`,
		`(fn [a] (+ a ((fn [a b] (+ a b)) a 1)))`,
	)
}

func TestRewriteQuoteWithMultiOccuranceOfSingelVariable(t *testing.T) {
	rw := cmp.NewQuoteRewriter()
	testRewriter(t, rw, `'(* ~n ~n ~n)`, `(fn [n] (* n n n))`)
}

func TestRewriteArgsSimple(t *testing.T) {
	pars := []string{"a"}
	args := []cmp.Node{parse("(+ 1 1)")}
	rw := cmp.NewArgsRewriter(pars, args)
	testRewriter(t, rw, `(* a a)`, `(* (+ 1 1) (+ 1 1))`)
}

func TestRewriteArgsDeep(t *testing.T) {
	pars := []string{"a", "b"}
	args := []cmp.Node{parse("(+ 1 1)"), parse("(- 2)")}
	rw := cmp.NewArgsRewriter(pars, args)
	testRewriter(t, rw, `(* (* a b) a)`, `(* (* (+ 1 1) (- 2)) (+ 1 1))`)
}

func TestRewriteDefmacroSimple(t *testing.T) {
	is := `(do
    (defmacro defn [name args body] (def name (fn args body)))
    (defn inc [x] (+ x 1))
  )`
	es := `(do (def inc (fn [x] (+ x 1))))`
	rw := cmp.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroNested(t *testing.T) {
	is := `(do
    (defmacro m1 [a b] (m2 b a))
    (defmacro m2 [a b] (- a b))
    (m1 1 2)
  )`
	es := `(do (- 2 1))`
	rw := cmp.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteUse(t *testing.T) {
	paths := []string{util.SplisHomePath()}
	is := `(do
		(use "test/prelude")
		(inc 41)
	)`
	es := `(do
		(do
			(def inc (fn [x] (+ x 1)))
			(def dec (fn [x] (- x 1)))
		)
		(inc 41)
	)`
	rw := cmp.NewUseRewriter(paths)
	testRewriter(t, rw, is, es)
}

func testRewriter(t *testing.T, rw cmp.Rewriter, i string, e string) {
	t.Helper()
	in := parse(i)
	en := parse(e)
	an := rw.Rewrite(in)
	as := cmp.PrintAst(an)
	es := cmp.PrintAst(en)
	if equalNode(an, en) == false {
		t.Errorf("Mismatch Expecting \n  [%s]\n but got \n  [%s].", es, as)
	}
}

func parse(i string) cmp.Node {
	r := cmp.NewReader()
	p := cmp.NewParser()
	r.Load(i)
	return p.Parse(r)
}

func equalNode(l cmp.Node, r cmp.Node) bool {
	switch lx := l.(type) {
	case bool:
		if rx, ok := r.(bool); ok {
			return lx == rx
		}
	case int64:
		if rx, ok := r.(int64); ok {
			return lx == rx
		}
	case *cmp.SymbolNode:
		if rx, ok := r.(*cmp.SymbolNode); ok {
			return lx.Name == rx.Name
		}
	case []cmp.Node:
		if rx, ok := r.([]cmp.Node); ok {
			return equalList(lx, rx)
		}
	case *cmp.ListNode:
		if rx, ok := r.(*cmp.ListNode); ok {
			return equalList(lx.Items, rx.Items)
		}
	}
	return false
}

func equalList(l []cmp.Node, r []cmp.Node) bool {
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
