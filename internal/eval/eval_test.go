package eval_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
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

func TestEvalDiff(t *testing.T) {
	test(t, "(- 1)", "-1")
	test(t, "(- 2 1)", "1")
	test(t, "(- 1 2)", "-1")
}

func TestEvalDiff2(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! :a (+ 1 1))", "2")
	teste(t, env, "(- :a)", "-2")
}

func TestEvalDiff3(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! :a (+ 1 1))", "2")
	teste(t, env, "(- :a 1)", "1")
	teste(t, env, "(- 1 :a)", "-1")
}

func TestEvalInvalidDiff(t *testing.T) {
	test(t, "(-)", "  [ERROR]  ")
	test(t, "(- 1 1 1)", "  [ERROR]  ")
}

func TestEvalDef1(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! :a 1)", "1")
	testenv(t, env, ":a", "1")
	teste(t, env, ":a", "1")
}

func TestEvalDef2(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! :a (+ 1 1))", "2")
	testenv(t, env, ":a", "2")
	teste(t, env, ":a", "2")
}

func TestEvalDef3(t *testing.T) {
	env := data.NewEnv(nil)
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
	test(t, "(if 42 1 0)", "1")           // Error?
	test(t, `(if "" 1 0)`, "0")           // Error?
	test(t, `(if "x" 1 0)`, "1")          // Error?
	test(t, "(if :x 1 0)", "  [ERROR]  ") // Error? // :x is undefined.
	// // test(t, "(if () 1 0)", "0")      // Error?
	// // test(t, "(if (42) 1 0)", "1")    // Error?
	test(t, "(if [] 1 0)", "0")      // Error?
	test(t, "(if [42] 1 0)", "1")    // Error?
	test(t, "(if {} 1 0)", "0")      // Error?
	test(t, "(if {:a 42} 1 0)", "1") // Error?

	test(t, "(if false (+ 2 2) (+ 1 1))", "2")
	test(t, "(if true (+ 2 2) (+ 1 1))", "4")

	test(t, "(if (< 0 1) 1 0)", "1")
	test(t, "(if (< 0 0) 1 0)", "0")
	test(t, "(if (< 1 0) 1 0)", "0")

	test(t, "(if (> 0 1) 1 0)", "0")
	test(t, "(if (> 0 0) 1 0)", "0")
	test(t, "(if (> 1 0) 1 0)", "1")
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

func TestList(t *testing.T) {
	test(t, "(list)", "()")
	test(t, "(list 0)", "(0)")
	test(t, "(list 0 1)", "(0 1)")
}

func TestIsList(t *testing.T) {
	test(t, "(list? nil)", "false")
	test(t, "(list? 0)", "false")
	test(t, `(list? "a")`, "false")
	test(t, "(list? [])", "false")
	test(t, "(list? {})", "false")
	test(t, "(list? (list))", "true")
	test(t, "(list? (list 1))", "true")
}

func TestInvalidIsList(t *testing.T) {
	test(t, "(list?)", "  [ERROR]  ")
	test(t, "(list? () ())", "  [ERROR]  ")
}

func TestCount(t *testing.T) {
	test(t, "(count (list))", "0")
	test(t, "(count (list 1))", "1")
	test(t, "(count (list 1 2 3))", "3")
	test(t, "(count [])", "0")
	test(t, "(count [1])", "1")
	test(t, "(count [1 2 3 4 5])", "5")
	test(t, "(count {})", "0")
	test(t, "(count {a 1})", "1")
	test(t, "(count {a 1 b 2 c 3 d 4})", "4")
}

func TestInvalidCount(t *testing.T) {
	test(t, "(count nil)", "  [ERROR]  ")
	test(t, "(count 0)", "  [ERROR]  ")
	test(t, `(count "a")`, "  [ERROR]  ")
}

func TestEvalLT(t *testing.T) {
	test(t, "(< 0 1)", "true")
	test(t, "(< 0 0)", "false")
	test(t, "(< 1 0)", "false")
}

func TestEvalLT2(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! N 0)", "0")
	teste(t, env, "(< N 1)", "true")
	teste(t, env, "(< N 0)", "false")
}

func TestEvalLT3(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! N 1)", "1")
	teste(t, env, "(< N 1)", "false")
	teste(t, env, "(< N 0)", "false")
}

// TODO: variable second argument.

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

func TestEvalNilEquivalence(t *testing.T) {
	test(t, "(= nil nil)", "true")
	test(t, "(= 0 nil)", "false")
	test(t, "(= nil 0)", "false")
}

func TestEvalBooleanEquivalence(t *testing.T) {
	test(t, "(= false false)", "true")
	test(t, "(= false true)", "false")
	test(t, "(= true false)", "false")
	test(t, "(= true true)", "true")
}

func TestEvalNumberEquivalence(t *testing.T) {
	test(t, "(= 1 1)", "true")
	test(t, "(= 1 0)", "false")
	test(t, "(= 0 1)", "false")
	test(t, "(= (+ 1 1) 2)", "true")
	test(t, "(= 2 (+ 1 1))", "true")
}

func TestEvalEnvVarEquivalence(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! x 1", "")
	teste(t, env, "(= x 1)", "true")
	teste(t, env, "(= 1 x)", "true")
	teste(t, env, "(= x x)", "true")
}

func TestEvalStringEquivalence(t *testing.T) {
	test(t, `(= "" "")`, "true")
	test(t, `(= "x" "")`, "false")
	test(t, `(= "" "x")`, "false")
	test(t, `(= "xyz" "xyz")`, "true")
}

func TestEvalListEquivalence(t *testing.T) {
	test(t, `(= (list) nil)`, "false")
	test(t, `(= nil (list))`, "false")
	test(t, `(= (list) (list))`, "true")
	test(t, `(= (list 1) (list))`, "false")
	test(t, `(= (list 1) (list 1))`, "true")
	test(t, `(= (list 1) (list 1 2))`, "false")
	test(t, `(= (list 1 2) (list 1 2))`, "true")
	test(t, `(= (list 1 2) (list 2 1))`, "false")
}

func TestEvalVectorEquivalence(t *testing.T) {
	test(t, `(= [] nil)`, "false")
	test(t, `(= nil [])`, "false")
	test(t, `(= [] [])`, "true")
	test(t, `(= [1] []])`, "false")
	test(t, `(= [1] [1])`, "true")
	test(t, `(= [1] [1 2])`, "false")
	test(t, `(= [1 2] [1 2])`, "true")
	test(t, `(= [1 2] [2 1])`, "false")
}

func TestEvalHashMapEquivalence(t *testing.T) {
	test(t, `(= {} nil)`, "false")
	test(t, `(= nil {})`, "false")
	test(t, `(= {} {})`, "true")
	test(t, `(= {a 1} {})`, "false")
	test(t, `(= {a 1} {a 1})`, "true")
	test(t, `(= {a 1} {a 2})`, "false")
	test(t, `(= {a 1} {b 1})`, "false")
	test(t, `(= {a 1} {a 1 b 2})`, "false")
	test(t, `(= {a 1 b 2} {a 1 b 2})`, "true")
	test(t, `(= {a 1 b 2} {b 1 a 2})`, "true")
}

func TestEvalInvalidEquivalence(t *testing.T) {
	test(t, "(=)", "  [ERROR]  ")
	test(t, "(= 1 1 1)", "  [ERROR]  ")
}

func TestEvalFun(t *testing.T) {
	//test(t, "(fn* () 42)", "#<fn>")
	test(t, "((fn* () 42))", "42")
	test(t, "((fn* (a) a) 42)", "42")
	test(t, "((fn* (a b) b) 0 42)", "42")
	test(t, "((fn* (a b c) (+ a b c)) 1 2 3)", "6")
	test(t, "((fn* (f x) (f x)) (fn* (a) (+ 1 a)) 7)", "8")
	test(t, "(((fn* (a) (fn* (b) (+ a b))) 5) 7)", "12")
}

func TestEvalFun2(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! gen-plus5 (fn* () (fn* (b) (+ 5 b))))", "")
	teste(t, env, "(def! plus5 (gen-plus5))", "")
	teste(t, env, "(plus5 7)", "12")
}

func TestEvalFun3(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! gen-plusX (fn* (x) (fn* (b) (+ x b))))", "")
	teste(t, env, "(def! plus7 (gen-plusX 7))", "")
	teste(t, env, "(plus7 8)", "15")
}

func TestEvalFun4(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! iffun (fn* (N) (if (> N 0) 33 22)))", "")
	teste(t, env, "(iffun 0)", "22")
	teste(t, env, "(iffun 1)", "33")
}

func TestEvalFun5(t *testing.T) {
	env := data.NewEnv(nil)
	teste(t, env, "(def! sumdown (fn* (N) (if (> N 0) (+ N (sumdown (- N 1))) 0)))", "")
	teste(t, env, "(sumdown 1)", "1")
	teste(t, env, "(sumdown 2)", "3")
	teste(t, env, "(sumdown 6)", "21")
}

// TODO: implement =
// func TestEvalFun6(t *testing.T) {
// 	env := data.NewEnv(nil)
// 	teste(t, env, "(def! fib (fn* (N) (if (= N 0) 1 (if (= N 1) 1 (+ (fib (- N 1)) (fib (- N 2)))))))", "")
// 	teste(t, env, "(fib 1)", "1")
// 	teste(t, env, "(fib 2)", "2")
// 	teste(t, env, "(fib 4)", "5")
// }

func test(t *testing.T, i string, e string) {
	teste(t, data.NewEnv(nil), i, e)
}

func teste(t *testing.T, env data.Env, i string, e string) {
	r := read.NewReader()
	p := read.NewParser()
	w := print.NewPrinter()
	v := eval.NewEvaluator(env)
	r.Load(i)
	n := p.Parse(r)
	m := v.Eval(n)
	a := w.Print(m)
	for _, err := range v.Errors() {
		w.PrintError(err)
	}
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}

func testenv(t *testing.T, env data.Env, name string, e string) {
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
