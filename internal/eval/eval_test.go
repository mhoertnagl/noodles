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

func TestEvalInvalidSum(t *testing.T) {
	test(t, "(+ 1 1 x)", "  [ERROR]  ")
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

// TODO: Test HashMap let binding.

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

func TestEvalIf(t *testing.T) {
	test(t, "(if nil 1 0)", "0") // Error?
	test(t, "(if false 1 0)", "0")
	test(t, "(if true 1 0)", "1")
	test(t, "(if 42 1 0)", "1")      // Error?
	test(t, `(if "" 1 0)`, "0")      // Error?
	test(t, `(if "x" 1 0)`, "1")     // Error?
	test(t, "(if :x 1 0)", "0")      // Error? // :x is undefined.
	test(t, "(if () 1 0)", "0")      // Error?
	test(t, "(if (42) 1 0)", "1")    // Error?
	test(t, "(if [] 1 0)", "0")      // Error?
	test(t, "(if [42] 1 0)", "1")    // Error?
	test(t, "(if {} 1 0)", "0")      // Error?
	test(t, "(if {:a 42} 1 0)", "1") // Error?

	test(t, "(if false (+ 2 2) (+ 1 1))", "2")
	test(t, "(if true (+ 2 2) (+ 1 1))", "4")

	test(t, "(if (< 0 1) 1 0)", "1")
	test(t, "(if (< 0 0) 1 0)", "0")
	test(t, "(if (< 1 0) 1 0)", "0")
}

func TestEvalIfWithoutElse(t *testing.T) {
	test(t, "(if false 1)", "nil")
	test(t, "(if true 1)", "1")
}

func TestEvalInvalidIf(t *testing.T) {
	test(t, "(if)", "  [ERROR]  ")
	test(t, "(if true)", "  [ERROR]  ")
	test(t, "(if true 1 0 2)", "  [ERROR]  ")
}

func TestEvalLT(t *testing.T) {
	test(t, "(< 0 1)", "true")
	test(t, "(< 0 0)", "false")
	test(t, "(< 1 0)", "false")
}

func TestEvalInvalidLT(t *testing.T) {
	test(t, "(<)", "  [ERROR]  ")
	test(t, "(< 0)", "  [ERROR]  ")
	test(t, "(< 1 0 1)", "  [ERROR]  ")
	test(t, "(< x 0)", "  [ERROR]  ")
	test(t, "(< 0 x)", "  [ERROR]  ")
	test(t, "(< x x)", "  [ERROR]  ")
}

func TestEvalGT(t *testing.T) {
	test(t, "(> 0 1)", "false")
	test(t, "(> 0 0)", "false")
	test(t, "(> 1 0)", "true")
}

func TestEvalInvalidGT(t *testing.T) {
	test(t, "(>)", "  [ERROR]  ")
	test(t, "(> 0)", "  [ERROR]  ")
	test(t, "(> 1 0 1)", "  [ERROR]  ")
	test(t, "(> x 0)", "  [ERROR]  ")
	test(t, "(> 0 x)", "  [ERROR]  ")
	test(t, "(> x x)", "  [ERROR]  ")
}

func TestEvalLE(t *testing.T) {
	test(t, "(<= 0 1)", "true")
	test(t, "(<= 0 0)", "true")
	test(t, "(<= 1 0)", "false")
	test(t, "(<= x 0)", "  [ERROR]  ")
	test(t, "(<= 0 x)", "  [ERROR]  ")
	test(t, "(<= x x)", "  [ERROR]  ")
}

func TestEvalInvalidLE(t *testing.T) {
	test(t, "(<=)", "  [ERROR]  ")
	test(t, "(<= 0)", "  [ERROR]  ")
	test(t, "(<= 1 0 1)", "  [ERROR]  ")
}

func TestEvalGE(t *testing.T) {
	test(t, "(>= 0 1)", "false")
	test(t, "(>= 0 0)", "true")
	test(t, "(>= 1 0)", "true")
}

func TestEvalInvalidGE(t *testing.T) {
	test(t, "(>=)", "  [ERROR]  ")
	test(t, "(>= 0)", "  [ERROR]  ")
	test(t, "(>= 1 0 1)", "  [ERROR]  ")
	test(t, "(>= x 0)", "  [ERROR]  ")
	test(t, "(>= 0 x)", "  [ERROR]  ")
	test(t, "(>= x x)", "  [ERROR]  ")
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
	if n, ok := env.Lookup(name); ok {
		w := print.NewPrinter()
		a := w.Print(n)
		if a != e {
			t.Errorf("Env variable [%s] should be [%s] but is [%s]", name, e, a)
		}
	} else {
		t.Errorf("Env variable [%s] undefined", name)
	}
}
