package rwr_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/cmp"
	"github.com/mhoertnagl/splis2/internal/rwr"
	"github.com/mhoertnagl/splis2/internal/util"
)

func TestRewriteBoolean(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `true`, `true`)
}

func TestRewriteInteger(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `1`, `1`)
}

func TestRewriteSymbol(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `a`, `a`)
}

func TestRewriteVector(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `[1 2 3]`, `[1 2 3]`)
}

func TestRewriteList(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `(1 2 3)`, `(1 2 3)`)
}

func TestRewriteSimpleQuote(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ 1 1)`, `(fn [] (+ 1 1))`)
}

func TestRewriteReplacementQuote(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ ~a ~b)`, `(fn [a b] (+ a b))`)
}

func TestRewriteSpliceQuote(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `'(+ ~a ~@b)`, `(fn [a b] (+ a @b))`)
}

func TestRewriteNestedQuote(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw,
		`'(+ ~a ('(+ ~a ~b) a 1))`,
		`(fn [a] (+ a ((fn [a b] (+ a b)) a 1)))`,
	)
}

func TestRewriteQuoteWithMultiOccurenceOfSingelVariable(t *testing.T) {
	rw := rwr.NewQuoteRewriter()
	testRewriter(t, rw, `'(* ~n ~n ~n)`, `(fn [n] (* n n n))`)
}

func TestRewriteArgsSimple(t *testing.T) {
	pars := []string{"a"}
	args := []cmp.Node{parse("(+ 1 1)")}
	rw := rwr.NewArgsRewriter(pars, "", args)
	testRewriter(t, rw, `(* a a)`, `(* (+ 1 1) (+ 1 1))`)
}

func TestRewriteArgsVector(t *testing.T) {
	pars := []string{"a", "b", "c"}
	args := []cmp.Node{parse("1"), parse("2"), parse("3")}
	rw := rwr.NewArgsRewriter(pars, "", args)
	testRewriter(t, rw, `(nth 1 [a b c])`, `(nth 1 [1 2 3])`)
}

func TestRewriteArgsDeep(t *testing.T) {
	pars := []string{"a", "b"}
	args := []cmp.Node{parse("(+ 1 1)"), parse("(- 2)")}
	rw := rwr.NewArgsRewriter(pars, "", args)
	testRewriter(t, rw, `(* (* a b) a)`, `(* (* (+ 1 1) (- 2)) (+ 1 1))`)
}

func TestRewriteVarArgs(t *testing.T) {
	pars := []string{"a"}
	args := []cmp.Node{parse("1"), parse("2"), parse("3"), parse("4")}
	rw := rwr.NewArgsRewriter(pars, "b", args)
	testRewriter(t, rw, `(:: a b)`, `(:: 1 [2 3 4])`)
}

func TestRewriteArgDissolve1(t *testing.T) {
	pars := []string{"a"}
	args := []cmp.Node{parse("b")}
	rw := rwr.NewArgsRewriter(pars, "", args)
	testRewriter(t, rw, `(+ @a)`, `(+ (dissolve b))`)
}

func TestRewriteArgDissolve2(t *testing.T) {
	pars := []string{"a"}
	args := []cmp.Node{parse("[1 2 3]")}
	rw := rwr.NewArgsRewriter(pars, "", args)
	testRewriter(t, rw, `(+ @a)`, `(+ 1 2 3)`)
}

func TestRewriteVarArgsDissolve(t *testing.T) {
	args := []cmp.Node{parse("(+ 1 1)"), parse("(- 2 2)")}
	rw := rwr.NewArgsRewriter([]string{}, "a", args)
	testRewriter(t, rw, `(do @a)`, `(do (+ 1 1) (- 2 2))`)
}

func TestRewriteDefmacroSimple0(t *testing.T) {
	is := `(do
    (defmacro else [] true)
    (cond else 1)
  )`
	es := `(do (cond true 1))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroSimple1(t *testing.T) {
	is := `(do
    (defmacro defn [name args & body]
      (def name (fn args (do @body))))
    (defn inc [x] (+ x 1) (- x 1))
  )`
	es := `(do (def inc (fn [x] (do (+ x 1) (- x 1)))))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroSimple2(t *testing.T) {
	is := `(do
    (defmacro defn [name args & body]
      (def name (fn args (do @body))))
    (defn add [x y] (+ x y) (- x 1))
  )`
	es := `(do (def add (fn [x y] (do (+ x y) (- x 1)))))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroNested(t *testing.T) {
	is := `(do
    (defmacro m1 [a b] (m2 b a))
    (defmacro m2 [a b] (- a b))
    (m1 1 2)
  )`
	es := `(do (- 2 1))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroNameClash(t *testing.T) {
	is := `(do
    (defmacro m1 [& x] (+ @x 1))
    (m1 x y)
  )`
	es := `(do (+ x y 1))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroVarArg1(t *testing.T) {
	is := `(do
    (defmacro defn [name args & body]
      (def name (fn args (do @body))))
    (defn vec [& args] args)
  )`
	es := `(do (def vec (fn [& args] (do args))))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroVarArg2(t *testing.T) {
	is := `(do
    (defmacro m1 [a & b] (:: a b))
    (m1 1 2 3 4)
  )`
	es := `(do (:: 1 [2 3 4]))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroPrint1(t *testing.T) {
	is := `(do
    (defmacro print [& args] (write *STD-OUT* @args))
    (print "Hello" ", " "World" "!")
  )`
	es := `(do (write *STD-OUT* "Hello" ", " "World" "!"))`
	rw := rwr.NewMacroRewriter()
	testRewriter(t, rw, is, es)
}

func TestRewriteDefmacroPrint2(t *testing.T) {
	is := `(do
    (defmacro print [& args] (write *STD-OUT* @args))
    (defmacro println [& args] (print @args "\n"))
    (println "Hello, World!")
  )`
	es := `(do (write *STD-OUT* "Hello, World!" "\n"))`
	rw := rwr.NewMacroRewriter()
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
	rw := rwr.NewUseRewriter(paths)
	testRewriter(t, rw, is, es)
}

// TODO: Move to rewriter_test
//
// func TestCompileSimpleQuote(t *testing.T) {
// 	testc(t, `'(+ 1 1)`,
// 		asm.Instr(vm.OpRef, 10),
//
// 		// (fn [] (+ 1 1))
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpConst, 1),
// 		asm.Instr(vm.OpConst, 1),
// 		asm.Instr(vm.OpAdd),
// 		asm.Instr(vm.OpReturn),
// 	)
// }
//
// func TestCompileReplacementQuote(t *testing.T) {
// 	testc(t, `'(+ ~a ~b)`,
// 		asm.Instr(vm.OpRef, 10),
//
// 		// (fn [a b] (+ a b))
// 		asm.Instr(vm.OpPushArgs, 2),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpAdd),
// 		asm.Instr(vm.OpReturn),
// 	)
// }
//
// func TestCompileSpliceQuote(t *testing.T) {
// 	testc(t, `'(+ ~a ~@b)`,
// 		asm.Instr(vm.OpRef, 10),
//
// 		// (fn [a b] (+ a @b))
// 		asm.Instr(vm.OpPushArgs, 2),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpDissolve),
// 		asm.Instr(vm.OpAdd),
// 		asm.Instr(vm.OpReturn),
// 	)
// }
//
// func TestCompileSpliceQuote2(t *testing.T) {
// 	testc(t, `'(+ ~@a ~@b)`,
// 		asm.Instr(vm.OpRef, 10),
//
// 		// (fn [a b] (+ @a @b))
// 		asm.Instr(vm.OpPushArgs, 2),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpDissolve),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpDissolve),
// 		asm.Instr(vm.OpAdd),
// 		asm.Instr(vm.OpReturn),
// 	)
// }
//
// func TestCompileQuote3(t *testing.T) {
// 	code := `
//   (do
//     (def cube '(* ~n ~n ~n))
//     (cube 3)
//   )
//   `
// 	testc(t, code,
// 		asm.Instr(vm.OpRef, 39),
// 		asm.Instr(vm.OpSetGlobal, 0),
// 		asm.Instr(vm.OpEnd),
// 		asm.Instr(vm.OpConst, 3),
// 		asm.Instr(vm.OpGetGlobal, 0),
// 		asm.Instr(vm.OpCall),
//
// 		asm.Instr(vm.OpPushArgs, 1),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpMul),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpMul),
// 		asm.Instr(vm.OpReturn),
// 	)
// }
//

func testRewriter(t *testing.T, rw rwr.Rewriter, i string, e string) {
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
	case string:
		if rx, ok := r.(string); ok {
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
