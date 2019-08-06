package eval_test

import (
	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"testing"
)

func TestEvalNumbers(t *testing.T) {
	test(t, "1", "1")
	test(t, "1.1", "1.1")
}

func TestEvalSum(t *testing.T) {
	test(t, "(+)", "0")
	test(t, "(+ 1)", "1")
	test(t, "(+ 1 1)", "2")
	test(t, "(+ 1 1 1)", "3")
	test(t, "(+ 1 (+ 1 1))", "3")
	test(t, "(+ (+ 1 1) (+ 1 1))", "4")
}

func TestEvalDef1(t *testing.T) {
	env := eval.NewEnv(nil)
	teste(t, env, "(def! :a 1)", "1")
	testenv(t, env, ":a", "1")
	teste(t, env, ":a", "1")
}

func TestEvalDef2(t *testing.T) {
	env := eval.NewEnv(nil)
	teste(t, env, "(def! :a (+ 1 1))", "2")
	testenv(t, env, ":a", "2")
	teste(t, env, ":a", "2")
}

func TestEvalDef3(t *testing.T) {
	env := eval.NewEnv(nil)
	// Define
	teste(t, env, "(def! :a 1)", "1")
	testenv(t, env, ":a", "1")
	teste(t, env, ":a", "1")
	// Redefine
	teste(t, env, "(def! :a 2)", "2")
	testenv(t, env, ":a", "2")
	teste(t, env, ":a", "2")
}

func TestEvalInvalidDef(t *testing.T) {
	test(t, "(def!)", "  [ERROR]  ")
	test(t, "(def! :a)", "  [ERROR]  ")
	test(t, "(def! :a 1 :b)", "  [ERROR]  ")
  test(t, "(def! 5 2)", "  [ERROR]  ")
}

func TestEvalLet(t *testing.T) {
	test(t, "(let* (:a 1) :a)", "1")
	test(t, "(let* (:a (+ 1 1)) :a)", "2")
	test(t, "(let* (:a 1) (+ :a :a))", "2")
	test(t, "(let* (:a (+ 1 1)) (+ :a :a))", "4")
	test(t, "(let* (p (+ 2 3) q (+ 2 p)) (+ p q))", "12")
}

func TestEvalLetVectorBinding(t *testing.T) {
	test(t, "(let* [:a 1] :a)", "1")
	test(t, "(let* [p (+ 2 3) q (+ 2 p)] (+ p q))", "12")
	test(t, "(let* (a 5 b 6) [3 4 a [b 7] 8])", "[3 4 5 [6 7] 8]")
}

// TODO: Test outer environment.

func TestEvalInvalidLet(t *testing.T) {
	test(t, "(let*)", "  [ERROR]  ")
	test(t, "(let* (:a 1))", "  [ERROR]  ")
	test(t, "(let* (:a 1) :a :b)", "  [ERROR]  ")
}

func TestEvalDo(t *testing.T) {
  test(t, "(do)", "nil")
  test(t, "(do (+ 1 1) (+ 2 2))", "4")
	test(t, "(do (def! a 3) (def! b 7) (+ a b))", "10")
}

func test(t *testing.T, i string, e string) {
	teste(t, eval.NewEnv(nil), i, e)
}

func teste(t *testing.T, env eval.Env, i string, e string) {
	r := read.NewReader()
	r.Load(i)
	p := read.NewParser()
	n := p.Parse(r)
	w := print.NewPrinter()
	v := eval.NewEvaluator(env)
	m := v.Eval(n)
	a := w.Print(m)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}

func testenv(t *testing.T, env eval.Env, name string, e string) {
	n := env.Lookup(name)
	w := print.NewPrinter()
	a := w.Print(n)
	if a != e {
		t.Errorf("Env variable [%s] should be [%s] but is [%s]", name, e, a)
	}
}
