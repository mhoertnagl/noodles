package eval_test

import (
	// "fmt"
	"fmt"
	"testing"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
)

func TestEmptyList(t *testing.T) {
	test(t, "()", "()")
}

func TestNumbers(t *testing.T) {
	test(t, "1", "1")
	test(t, "1.1", "1.1")
}

func TestSum(t *testing.T) {
	test(t, "(+)", "0")
	test(t, "(+ 1)", "1")
	test(t, "(+ 1 1)", "2")
	test(t, "(+ 1 1 1)", "3")
	test(t, "(+ 1 (+ 1 1))", "3")
	test(t, "(+ (+ 1 1) (+ 1 1))", "4")
}

func TestInvalidSum(t *testing.T) {
	test(t, "(+ 1 1 x)", "  [ERROR]  ")
}

func TestDiff(t *testing.T) {
	test(t, "(- 0)", "0")
	test(t, "(- 1)", "-1")
	test(t, "(- 2 1)", "1")
	test(t, "(- 1 2)", "-1")
}

func TestDiff2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def :a (+ 1 1))", "2")
	teste(t, env, "(- :a)", "-2")
}

func TestDiff3(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def :a (+ 1 1))", "2")
	teste(t, env, "(- :a 1)", "1")
	teste(t, env, "(- 1 :a)", "-1")
}

func TestInvalidDiff(t *testing.T) {
	test(t, "(-)", "  [ERROR]  ")
	test(t, "(- :a)", "  [ERROR]  ")
	test(t, "(- false)", "  [ERROR]  ")
	test(t, "(- false 1)", "  [ERROR]  ")
	test(t, "(- 1 false)", "  [ERROR]  ")
	test(t, "(- false false)", "  [ERROR]  ")
	test(t, "(- 1 1 1)", "  [ERROR]  ")
}

func TestSumAndDiff(t *testing.T) {
	test(t, "(+ 1 (- 0))", "1")
	test(t, "(+ (- 0) (- 0))", "0")
}

func TestProd(t *testing.T) {
	test(t, "(*)", "1")
	test(t, "(* 0)", "0")
	test(t, "(* 0 1)", "0")
	test(t, "(* 2 3)", "6")
	test(t, "(* 2 (- 3))", "-6")
}

func TestDiv(t *testing.T) {
	test(t, "(/ 0)", "+Inf")
	test(t, "(/ (- 1) 0)", "-Inf")
	test(t, "(/ 0 1)", "0")
	test(t, "(/ 2)", "0.5")
	test(t, "(/ 2 (- 4))", "-0.5")
}

func TestInvalidDiv(t *testing.T) {
	test(t, "(/)", "  [ERROR]  ")
}

func TestHashMap(t *testing.T) {
	test(t, `{}`, `{}`)
	test(t, `{"a" 1}`, `{"a" 1}`)
	// test(t, `{"a" 1 "b" 2}`, `{"a" 1 "b" 2}`)
	// TODO: Support for keywords (beginning with a :)
	// test(t, `{:a 1}`, `{:a 1}`)
}

func TestDef1(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def :a 1)", "1")
	testenv(t, env, ":a", "1")
	teste(t, env, ":a", "1")
}

func TestDef2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def :a (+ 1 1))", "2")
	testenv(t, env, ":a", "2")
	teste(t, env, ":a", "2")
}

func TestDef3(t *testing.T) {
	env := data.NewRootEnv()
	// Define
	teste(t, env, "(def :a 1)", "1")
	testenv(t, env, ":a", "1")
	teste(t, env, ":a", "1")
	// Redefine
	teste(t, env, "(def :a 2)", "2")
	testenv(t, env, ":a", "2")
	teste(t, env, ":a", "2")
}

func TestInvalidDef(t *testing.T) {
	test(t, "(def)", "  [ERROR]  ")
	test(t, "(def :a)", "  [ERROR]  ")
	test(t, "(def :a 1 :b)", "  [ERROR]  ")
	test(t, "(def 5 2)", "  [ERROR]  ")
}

func TestDefPI(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def *PI* 3.141592653589793)", "3.141592653589793")
	testenv(t, env, "*PI*", "3.141592653589793")
}

func TestLet(t *testing.T) {
	test(t, "(let (:a 1) :a)", "1")
	test(t, "(let (:a (+ 1 1)) :a)", "2")
	test(t, "(let (:a 1) (+ :a :a))", "2")
	test(t, "(let (:a (+ 1 1)) (+ :a :a))", "4")
	test(t, "(let (p (+ 2 3) q (+ 2 p)) (+ p q))", "12")
}

func TestLetVectorBinding(t *testing.T) {
	test(t, "(let [:a 1] :a)", "1")
	test(t, "(let [p (+ 2 3) q (+ 2 p)] (+ p q))", "12")
	test(t, "(let (a 5 b 6) [3 4 a [b 7] 8])", "[3 4 5 [6 7] 8]")
}

// func TestLetHashMapBinding(t *testing.T) {
// 	test(t, "(let {:a 1} :a)", "1")
// }

// TODO: Test outer environment.

func TestInvalidLet(t *testing.T) {
	test(t, "(let)", "  [ERROR]  ")
	test(t, "(let 1 2)", "  [ERROR]  ")
	test(t, "(let (:a 1))", "  [ERROR]  ")
	test(t, "(let (:a 1) :a :b)", "  [ERROR]  ")
}

func TestDo(t *testing.T) {
	test(t, "(do)", "nil")
	test(t, "(do (+ 1 1))", "2")
	test(t, "(do (+ 1 1) (+ 2 2))", "4")
	test(t, "(do (def a 3) (def b 7) (+ a b))", "10")
}

func TestCond(t *testing.T) {
	test(t, "(cond true 1 false 0)", "1")
	test(t, "(cond false 1 true 0)", "0")
	test(t, "(cond (> 2 1) 2 (> 2 1) 1 (> 2 1) 0)", "2")
	test(t, "(cond (> 2 2) 2 (> 2 1) 1 (> 2 1) 0)", "1")
	test(t, "(cond (> 2 2) 2 (> 2 2) 1 (> 2 1) 0)", "0")
}

func TestIf(t *testing.T) {
	test(t, "(if nil 1 0)", "0") // Error?
	test(t, "(if false 1 0)", "0")
	test(t, "(if true 1 0)", "1")
	test(t, "(if 42 1 0)", "1")           // Error?
	test(t, `(if "" 1 0)`, "0")           // Error?
	test(t, `(if "x" 1 0)`, "1")          // Error?
	test(t, "(if :x 1 0)", "  [ERROR]  ") // Error? // :x is undefined.
	// // test(t, "(if () 1 0)", "0")      // Error?
	// // test(t, "(if (42) 1 0)", "1")    // Error?
	test(t, "(if [] 1 0)", "0")       // Error?
	test(t, "(if [42] 1 0)", "1")     // Error?
	test(t, "(if {} 1 0)", "0")       // Error?
	test(t, `(if {"a" 42} 1 0)`, "1") // Error?

	test(t, "(if false (+ 2 2) (+ 1 1))", "2")
	test(t, "(if true (+ 2 2) (+ 1 1))", "4")

	test(t, "(if (< 0 1) 1 0)", "1")
	test(t, "(if (< 0 0) 1 0)", "0")
	test(t, "(if (< 1 0) 1 0)", "0")

	test(t, "(if (> 0 1) 1 0)", "0")
	test(t, "(if (> 0 0) 1 0)", "0")
	test(t, "(if (> 1 0) 1 0)", "1")
}

func TestIfWithoutElse(t *testing.T) {
	test(t, "(if false 1)", "nil")
	test(t, "(if true 1)", "1")
}

func TestInvalidIf(t *testing.T) {
	test(t, "(if)", "  [ERROR]  ")
	test(t, "(if true)", "  [ERROR]  ")
	test(t, "(if true 1 0 2)", "  [ERROR]  ")
}

func TestList(t *testing.T) {
	test(t, "(list)", "()")
	test(t, "(list 0)", "(0)")
	test(t, "(list 0 1)", "(0 1)")
	test(t, "(list (+ 4 4) 8)", "(8 8)")
	test(t, "(list + 4 4)", "(+ 4 4)")
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
	test(t, `(count {"a" 1})`, "1")
	test(t, `(count {"a" 1 "b" 2 "c" 3 "d" 4})`, "4")
}

func TestInvalidCount(t *testing.T) {
	test(t, "(count)", "  [ERROR]  ")
	test(t, "(count nil)", "  [ERROR]  ")
	test(t, "(count 0)", "  [ERROR]  ")
	test(t, `(count "a")`, "  [ERROR]  ")
	test(t, `(count () ())`, "  [ERROR]  ")
}

func TestEmpty(t *testing.T) {
	test(t, "(empty? (list))", "true")
	test(t, "(empty? (list 1))", "false")
	test(t, "(empty? (list 1 2 3))", "false")
	test(t, "(empty? [])", "true")
	test(t, "(empty? [1])", "false")
	test(t, "(empty? [1 2 3 4 5])", "false")
	test(t, "(empty? {})", "true")
	test(t, `(empty? {"a" 1})`, "false")
	test(t, `(empty? {"a" 1 "b" 2 "c" 3 "d" 4})`, "false")
}

func TestInvalidEmpty(t *testing.T) {
	test(t, "(empty?)", "  [ERROR]  ")
	test(t, "(empty? nil)", "  [ERROR]  ")
	test(t, "(empty? 0)", "  [ERROR]  ")
	test(t, `(empty? "a")`, "  [ERROR]  ")
	test(t, `(empty? () ())`, "  [ERROR]  ")
}

// TODO: rawstr -> escape sequenzen werden nicht übersetzt
// TODO: str    -> escape sequencen werden übersetzt

// func TestStr(t *testing.T) {
// 	test(t, `(str)`, `""`)
// 	test(t, `(str "")`, `""`)
// 	test(t, `(str "abc")`, `"abc"`)
// 	test(t, `(str "\"")`, `"""`)
// 	test(t, `(str 1 "abc" 3)`, `"1abc3"`)
// 	test(t, `(str "abc  def" "ghi jkl")`, `"abc  defghi jkl"`)
// 	test(t, `(str "abc\ndef\nghi")`, `"`+"abc\ndef\nghi"+`"`)
// 	test(t, `(str "abc\\def\\ghi")`, `"abc\def\ghi"`)
// 	test(t, `(str (list 1 2 "abc" "\"") "def")`, `"(1 2 abc ")def"`)
// 	test(t, `(str (list))`, `"()"`)
// }
//
// func TestRawstr(t *testing.T) {
// 	test(t, `(rawstr)`, `""`)
// 	test(t, `(rawstr "")`, `""`)
// 	test(t, `(rawstr "abc")`, `"abc"`)
// 	test(t, `(rawstr "\"")`, `"\""`)
// 	test(t, `(rawstr 1 "abc" 3)`, `"1abc3"`)
// 	test(t, `(rawstr "abc  def" "ghi jkl")`, `"abc  defghi jkl"`)
// 	test(t, `(rawstr "abc\ndef\nghi")`, `"abc\ndef\nghi"`)
// 	test(t, `(rawstr "abc\\def\\ghi")`, `"abc\\def\\ghi"`)
// 	test(t, `(rawstr (list 1 2 "abc" "\"") "def")`, `"(1 2 abc \")def"`)
// 	test(t, `(rawstr (list))`, `"()"`)
// }

func TestLT(t *testing.T) {
	test(t, "(< 0 1)", "true")
	test(t, "(< 0 0)", "false")
	test(t, "(< 1 0)", "false")
}

func TestLT2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def N 0)", "0")
	teste(t, env, "(< N 1)", "true")
	teste(t, env, "(< N 0)", "false")
}

func TestLT3(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def N 1)", "1")
	teste(t, env, "(< N 1)", "false")
	teste(t, env, "(< N 0)", "false")
}

// TODO: variable second argument.

func TestInvalidLT(t *testing.T) {
	test(t, "(<)", "  [ERROR]  ")
	test(t, "(< 0)", "  [ERROR]  ")
	test(t, "(< 1 0 1)", "  [ERROR]  ")
	test(t, "(< x 0)", "  [ERROR]  ")
	test(t, "(< 0 x)", "  [ERROR]  ")
	test(t, "(< x x)", "  [ERROR]  ")
}

func TestGT(t *testing.T) {
	test(t, "(> 0 1)", "false")
	test(t, "(> 0 0)", "false")
	test(t, "(> 1 0)", "true")
}

func TestInvalidGT(t *testing.T) {
	test(t, "(>)", "  [ERROR]  ")
	test(t, "(> 0)", "  [ERROR]  ")
	test(t, "(> 1 0 1)", "  [ERROR]  ")
	test(t, "(> x 0)", "  [ERROR]  ")
	test(t, "(> 0 x)", "  [ERROR]  ")
	test(t, "(> x x)", "  [ERROR]  ")
}

func TestLE(t *testing.T) {
	test(t, "(<= 0 1)", "true")
	test(t, "(<= 0 0)", "true")
	test(t, "(<= 1 0)", "false")
	test(t, "(<= x 0)", "  [ERROR]  ")
	test(t, "(<= 0 x)", "  [ERROR]  ")
	test(t, "(<= x x)", "  [ERROR]  ")
}

func TestInvalidLE(t *testing.T) {
	test(t, "(<=)", "  [ERROR]  ")
	test(t, "(<= 0)", "  [ERROR]  ")
	test(t, "(<= 1 0 1)", "  [ERROR]  ")
}

func TestGE(t *testing.T) {
	test(t, "(>= 0 1)", "false")
	test(t, "(>= 0 0)", "true")
	test(t, "(>= 1 0)", "true")
}

func TestInvalidGE(t *testing.T) {
	test(t, "(>=)", "  [ERROR]  ")
	test(t, "(>= 0)", "  [ERROR]  ")
	test(t, "(>= 1 0 1)", "  [ERROR]  ")
	test(t, "(>= x 0)", "  [ERROR]  ")
	test(t, "(>= 0 x)", "  [ERROR]  ")
	test(t, "(>= x x)", "  [ERROR]  ")
}

func TestNilEquivalence(t *testing.T) {
	test(t, "(= nil nil)", "true")
	test(t, "(= 0 nil)", "false")
	test(t, "(= nil 0)", "false")
}

func TestBooleanEquivalence(t *testing.T) {
	test(t, "(= false false)", "true")
	test(t, "(= false true)", "false")
	test(t, "(= true false)", "false")
	test(t, "(= true true)", "true")
}

func TestNumberEquivalence(t *testing.T) {
	test(t, "(= 1 1)", "true")
	test(t, "(= 1 0)", "false")
	test(t, "(= 0 1)", "false")
	test(t, "(= (+ 1 1) 2)", "true")
	test(t, "(= 2 (+ 1 1))", "true")
	test(t, "(= 1.1 1.1)", "true")
}

func TestEnvVarEquivalence(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def x 1", "1")
	teste(t, env, "(= x 1)", "true")
	teste(t, env, "(= 1 x)", "true")
	teste(t, env, "(= x x)", "true")
}

func TestStringEquivalence(t *testing.T) {
	test(t, `(= "" "")`, "true")
	test(t, `(= "x" "")`, "false")
	test(t, `(= "" "x")`, "false")
	test(t, `(= "xyz" "xyz")`, "true")
}

func TestListEquivalence(t *testing.T) {
	test(t, `(= (list) nil)`, "false")
	test(t, `(= nil (list))`, "false")
	test(t, `(= (list) (list))`, "true")
	test(t, `(= (list 1) (list))`, "false")
	test(t, `(= (list 1) (list 1))`, "true")
	test(t, `(= (list 1) (list 1 2))`, "false")
	test(t, `(= (list 1 2) (list 1 2))`, "true")
	test(t, `(= (list 1 2) (list 2 1))`, "false")
}

func TestVectorEquivalence(t *testing.T) {
	test(t, `(= [] nil)`, "false")
	test(t, `(= nil [])`, "false")
	test(t, `(= [] [])`, "true")
	test(t, `(= [1] [])`, "false")
	test(t, `(= [] [1])`, "false")
	test(t, `(= [1] [1])`, "true")
	test(t, `(= [1] [1 2])`, "false")
	test(t, `(= [1 2] [1 2])`, "true")
	test(t, `(= [1 2] [2 1])`, "false")
}

func TestHashMapEquivalence(t *testing.T) {
	test(t, `(= {} nil)`, "false")
	test(t, `(= nil {})`, "false")
	test(t, `(= {} {})`, "true")
	test(t, `(= {"a" nil} {"b" 1})`, "false")
	test(t, `(= {"a" 1} {"b" nil})`, "false")
	test(t, `(= {"a" 1} {})`, "false")
	test(t, `(= {} {"a" 1})`, "false")
	test(t, `(= {"a" 1} {"a" 1})`, "true")
	test(t, `(= {"a" 1} {"a" 2})`, "false")
	test(t, `(= {"a" 1} {"b" 1})`, "false")
	test(t, `(= {"a" 1} {"a" 1 "b" 2})`, "false")
	test(t, `(= {"a" 1 "b" 2} {"a" 1})`, "false")
	test(t, `(= {"a" 1 "b" 2} {"a" 1 "b" 2})`, "true")
	test(t, `(= {"a" 1 "b" 2} {"b" 2 "a" 1})`, "true")
	test(t, `(= {"a" 1 "b" 2} {"c" 1 "d" 2})`, "false")
}

func TestInvalidEquivalence(t *testing.T) {
	test(t, "(=)", "  [ERROR]  ")
	test(t, "(= 1 1 1)", "  [ERROR]  ")
}

func TestFun(t *testing.T) {
	//test(t, "(fn () 42)", "#<fn>")
	test(t, "((fn () 40))", "40")
	test(t, "((fn (a) a) 41)", "41")
	test(t, "((fn (a b) b) 0 42)", "42")
	test(t, "((fn (a b c) (+ a b c)) 1 2 3)", "6")
	test(t, "((fn (f x) (f x)) (fn (a) (+ 1 a)) 7)", "8")
	test(t, "(((fn (a) (fn (b) (+ a b))) 5) 7)", "12")
	test(t, "((fn [x] (if x false true)) true)", "false")
	test(t, "((fn [x] (if x false true)) false)", "true")
	test(t, "((fn [f x] (f x)) (fn [a] (+ 1 a)) 7)", "8")
	// test(t, "(((fn [a b] b)) 0 43)", "43")
}

// func TestPartialFun(t *testing.T) {
// 	test(t, "(((fn [a b] b) 0) 43)", "43")
// }

func TestFun2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def gen-plus5 (fn () (fn (b) (+ 5 b))))", "")
	teste(t, env, "(def plus5 (gen-plus5))", "")
	teste(t, env, "(plus5 7)", "12")
}

func TestFun3(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def gen-plusX (fn (x) (fn (b) (+ x b))))", "")
	teste(t, env, "(def plus7 (gen-plusX 7))", "")
	teste(t, env, "(plus7 8)", "15")
}

func TestFun4(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def iffun (fn (N) (if (> N 0) 33 22)))", "")
	teste(t, env, "(iffun 0)", "22")
	teste(t, env, "(iffun 1)", "33")
}

func TestFun5(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def sumdown (fn (N) (if (> N 0) (+ N (sumdown (- N 1))) 0)))", "")
	teste(t, env, "(sumdown 1)", "1")
	teste(t, env, "(sumdown 2)", "3")
	teste(t, env, "(sumdown 6)", "21")
}

func TestFib(t *testing.T) {
	env := data.NewRootEnv()
	fib := `
		(def fib (fn (N)
			(if (= N 0)
				1
				(if (= N 1)
					1
					(+ (fib (- N 1)) (fib (- N 2)))))))
	`
	teste(t, env, fib, "")
	teste(t, env, "(fib 0)", "1")
	teste(t, env, "(fib 1)", "1")
	teste(t, env, "(fib 2)", "2")
	teste(t, env, "(fib 4)", "5")
	teste(t, env, "(fib 5)", "8")
}

func TestSum2(t *testing.T) {
	env := data.NewRootEnv()
	sum2 := `
		(def sum2 (fn (n acc)
			(if (= n 0)
				acc
				(sum2 (- n 1) (+ n acc)))))
	`
	teste(t, env, sum2, "")
	teste(t, env, "(sum2 10 0)", "55")
	// NOTE: Slow test.
	// teste(t, env, "(def res2 nil)", "nil")
	// teste(t, env, "(def res2 (sum2 10000 0))", "50005000")
}

func TestNotDef(t *testing.T) {
	env := data.NewRootEnv()
	not := `
    (def not (fn [x]
      (if x false true)))
	`
	teste(t, env, not, "")
	teste(t, env, "(not true)", "false")
	teste(t, env, "(not false)", "true")
}

func TestFunShorthand(t *testing.T) {
	env := data.NewRootEnv()
	not := `
    (do
      (defmacro fun
        (fn [name args body]
          (quasiquote (def ~name (fn ~args ~body)) )))

      (fun plus7 [x] (+ 7 x))
    )
	`
	teste(t, env, not, "")
	teste(t, env, "(plus7 1)", "8")
	teste(t, env, "(plus7 11)", "18")
}

//
// NOTE: Slow test.
// func TestMutualRecursive(t *testing.T) {
// 	env := data.NewRootEnv()
// 	teste(t, env, "(def foo (fn (n) (if (= n 0) 0 (bar (- n 1)))))", "")
// 	teste(t, env, "(def bar (fn (n) (if (= n 0) 0 (foo (- n 1)))))", "")
// 	teste(t, env, "(foo 10000)", "0")
// }
//
func TestParse(t *testing.T) {
	test(t, `(parse "(1 2 (3 4) nil)")`, "(1 2 (3 4) nil)")
	test(t, `(parse "(+ 2 3)")`, "(+ 2 3)")
	test(t, `(parse "7 ;; comment")`, "7")
	test(t, `(parse ";; comment")`, "")
}

func TestRead2(t *testing.T) {
	env := data.NewRootEnv()
	src := `
  ;; Returns the negation of x.
  (def not (fn [x]
    (if x false true)))
	`
	teste(t, env, src, "")
	teste(t, env, "(not true)", "false")
	teste(t, env, "(not false)", "true")
}

func TestEval(t *testing.T) {
	test(t, "(eval 4)", "4")
	test(t, "(eval (+ 4 4))", "8")
	test(t, `(eval (parse "(+ 2 3)"))`, "5")
}

func TestEval2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def mal-prog (list + 1 2))", "(+ 1 2)")
	teste(t, env, "(eval mal-prog)", "3")
}

func TestReadFile(t *testing.T) {
	exp := `"This is a Test.\n` + "\n" + `"`
	test(t, `(read-file "../../test/test.txt")`, exp)
}

func TestInvalidReadFile(t *testing.T) {
	test(t, `(read-file)`, "  [ERROR]  ")
	test(t, `(read-file true)`, "  [ERROR]  ")
	test(t, `(read-file "foo.bar")`, "  [ERROR]  ")
}

func TestExecuteFile(t *testing.T) {
	src := `
    (eval
      (parse
        (read-file "../../test/test.splis")))`
	env := data.NewRootEnv()
	teste(t, env, src, "")
	teste(t, env, "(plus5 1)", "6")
	teste(t, env, "(plus5 9)", "14")
}

func TestUseModule(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, `(use "test/test")`, "")
	teste(t, env, "(plus5 (plus5 0))", "10")
}

func TestCons(t *testing.T) {
	test(t, "(:: 42 (list))", "(42)")
	test(t, "(:: 1 (list 2 3 4))", "(1 2 3 4)")
	test(t, "(:: (list 1 2) (list))", "((1 2))")
	test(t, "(:: (list 1 2) (list (list 3 4)))", "((1 2) (3 4))")
}

// func TestConsVectors(t *testing.T) {
// 	test(t, "(:: 42 [])", "(42)")
// 	test(t, "(:: 1 [2 3 4])", "(1 2 3 4)")
// 	test(t, "(:: [1 2] []", "([1 2])")
// 	test(t, "(:: [1 2] [[3 4]])", "([1 2] [3 4])")
// }

func TestInvalidCons(t *testing.T) {
	test(t, "(::)", "  [ERROR]  ")
	test(t, "(:: 1 (2 3 4) 5)", "  [ERROR]  ")
	test(t, "(:: (1 2) 2)", "  [ERROR]  ")
}

func TestConcat(t *testing.T) {
	test(t, "(:::)", "()")
	test(t, "(::: (list))", "()")
	test(t, "(::: (list) (list))", "()")
	test(t, "(::: (list) (list) (list))", "()")
	test(t, "(::: (list 1 2) (list))", "(1 2)")
	test(t, "(::: (list) (list 3 4))", "(3 4)")
	test(t, "(::: (list 1 2) (list 3 4))", "(1 2 3 4)")
	test(t, "(::: (list 1 2) (list 3 4) (list 5 6))", "(1 2 3 4 5 6)")
	test(t, "(::: (list (list 1 2)) (list (list 3 4)))", "((1 2) (3 4))")
}

// func TestConcatVectors(t *testing.T) {
// 	test(t, "(::: [])", "()")
// 	test(t, "(::: [] [])", "()")
// 	test(t, "(::: [] [] [])", "()")
// 	test(t, "(::: [1 2] [])", "(1 2)")
// 	test(t, "(::: [] [3 4])", "(3 4)")
// 	test(t, "(::: [1 2] [3 4])", "(1 2 3 4)")
// 	test(t, "(::: [1 2] [3 4] [5 6])", "(1 2 3 4 5 6)")
// 	test(t, "(::: [[1 2]] [[3 4]])", "([1 2] [3 4])")
// }

func TestConcatEnv(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def a (list 1 2))", "(1 2)")
	teste(t, env, "(def b (list 3 4))", "(3 4)")
	teste(t, env, "(::: a b (list 5 6))", "(1 2 3 4 5 6)")
}

func TestInvalidConcat(t *testing.T) {
	test(t, "(::: 2)", "  [ERROR]  ")
	test(t, "(::: (1 2) 2)", "  [ERROR]  ")
	test(t, "(::: 1 (2 3 4) 5)", "  [ERROR]  ")
}

func TestQuote(t *testing.T) {
	test(t, "(quote 42)", "42")
	test(t, "(quote (1 2 3))", "(1 2 3)")
	test(t, "(quote (1 (2 (3))))", "(1 (2 (3)))")
	test(t, "(quote [6 5 4])", "[6 5 4]")
}

func TestQuote2(t *testing.T) {
	test(t, "(= (quote abc) (quote abc))", "true")
	test(t, "(= (quote abc) (quote abcd))", "false")
	test(t, `(= (quote abc) "abc")`, "false")
	test(t, `(= "abc" (quote abc))`, "false")
	// test(t, `(= "abc" (str (quote abc)))`, "true")
	test(t, "(= (quote abc) nil)", "false")
	test(t, "(= nil (quote abc))", "false")
}

func TestInvalidQuote(t *testing.T) {
	test(t, "(quote)", "  [ERROR]  ")
	test(t, "(quote 1 2)", "  [ERROR]  ")
}

func TestQuasiquote(t *testing.T) {
	test(t, "(quasiquote ())", "()")
	// test(t, "(quasiquote (()))", "(())")
	test(t, "(quasiquote 43)", "43")
	test(t, "(quasiquote (1 2 3))", "(1 2 3)")
	test(t, "(quasiquote (1 (2 (3))))", "(1 (2 (3)))")
	test(t, "(quasiquote [6 5 4])", "[6 5 4]")
	test(t, "(quasiquote (nil))", "(nil)")
	test(t, "(quasiquote (unquote 7))", "7")
}

func TestQuasiquote2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def a 8)", "8")
	teste(t, env, "(quasiquote a)", "a")
	teste(t, env, "(quasiquote (unquote a))", "8")
	teste(t, env, "(quasiquote (1 a 3))", "(1 a 3)")
	teste(t, env, "(quasiquote (1 (unquote a) 3))", "(1 8 3)")

	teste(t, env, `(def b (quote (1 "b" "d")))`, `(1 "b" "d")`)
	teste(t, env, "(quasiquote (1 b 3))", "(1 b 3)")
	teste(t, env, `(quasiquote (1 (unquote b) 3))`, `(1 (1 "b" "d") 3)`)
	teste(t, env, "(quasiquote ((unquote 1) (unquote 2)))", "(1 2)")

	teste(t, env, `(def c (quote (1 "b" "d")))`, `(1 "b" "d")`)
	teste(t, env, `(quasiquote (1 c 3))`, `(1 c 3)`)
	teste(t, env, `(quasiquote (1 (splice-unquote c) 3))`, `(1 1 "b" "d" 3)`)
}

func TestQuasiquoteQuine(t *testing.T) {
	test(t,
		"((fn [q] (quasiquote ((unquote q) (quote (unquote q))))) (quote (fn [q] (quasiquote ((unquote q) (quote (unquote q)))))))",
		"((fn [q] (quasiquote ((unquote q) (quote (unquote q))))) (quote (fn [q] (quasiquote ((unquote q) (quote (unquote q)))))))")
}

func TestInvalidQuasiquote(t *testing.T) {
	test(t, "(quasiquote)", "  [ERROR]  ")
	test(t, "(quasiquote 1 2)", "  [ERROR]  ")
}

func TestTrivialMacros(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(defmacro one (fn [] 1))", "")
	teste(t, env, "(one)", "1")
	teste(t, env, "(defmacro two (fn [] 2))", "")
	teste(t, env, "(two)", "2")
}

func TestUnlessMacros(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(defmacro unless (fn [pred a b] `(if ~pred ~b ~a)))", "")
	teste(t, env, "(unless false 7 8)", "7")
	teste(t, env, "(unless true 7 8)", "8")
}

func TestUnlessMacros2(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(defmacro unless2 (fn [pred a b] (list 'if (list 'not pred) a b)))", "")
	teste(t, env, "(unless2 false 7 8)", "7")
	teste(t, env, "(unless2 true 7 8)", "8")
	teste(t, env, "(macroexpand (unless2 2 3 4))", "(if (not 2) 3 4)")
}

func TestMacroResultEvaluation(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(identity 1)", "1")
	teste(t, env, "(let (a 123) (identity a))", "123")
}

func TestMacroUsesClosures(t *testing.T) {
	env := data.NewRootEnv()
	teste(t, env, "(def x 2)", "2")
	teste(t, env, "(defmacro a (fn [] x))", "")
	teste(t, env, "(a)", "2")
	teste(t, env, "(let (x 3) (a))", "2")
}

func TestInvalidDefMacro(t *testing.T) {
	test(t, "(defmacro)", "  [ERROR]  ")
	test(t, "(defmacro a)", "  [ERROR]  ")
	test(t, "(defmacro 2 (fn [] 0))", "  [ERROR]  ")
	test(t, "(defmacro a 2)", "  [ERROR]  ")
}

func TestHead(t *testing.T) {
	test(t, "(head '(true)", "true")
	test(t, "(head '(1 2 3))", "1")
	test(t, "(head '(a b))", "a")

	test(t, "(head '[false]", "false")
	test(t, "(head '[3 2 1])", "3")
	test(t, "(head '[b a])", "b")
}

func TestInvalidHead(t *testing.T) {
	test(t, "(head x", "  [ERROR]  ")
	test(t, "(head '()", "  [ERROR]  ")
	test(t, "(head '[]", "  [ERROR]  ")
}

func TestTail(t *testing.T) {
	test(t, "(tail '(true)", "()")
	test(t, "(tail '(1 2 3))", "(2 3)")
	test(t, "(tail '(a b))", "(b)")

	test(t, "(tail '[false]", "[]")
	test(t, "(tail '[3 2 1])", "[2 1]")
	test(t, "(tail '[b a])", "[a]")
}

func TestInvalidTail(t *testing.T) {
	test(t, "(tail x", "  [ERROR]  ")
	test(t, "(tail '()", "  [ERROR]  ")
	test(t, "(tail '[]", "  [ERROR]  ")
}

func TestPreludeNot(t *testing.T) {
	test(t, "(not true)", "false")
	test(t, "(not false)", "true")
}

func TestPreludeInc(t *testing.T) {
	test(t, "(inc 0)", "1")
	test(t, "(inc 6)", "7")

	test(t, "(dec 0)", "-1")
	test(t, "(dec 6)", "5")
}

func TestPreludeIsPositive(t *testing.T) {
	test(t, "(pos? (- 1))", "false")
	test(t, "(pos? 0)", "false")
	test(t, "(pos? 1)", "true")
}

func TestPreludeAbs(t *testing.T) {
	test(t, "(abs (- 1))", "1")
	test(t, "(abs 0)", "0")
	test(t, "(abs 1)", "1")
}

func TestPreludeIsNil(t *testing.T) {
	test(t, "(nil? nil)", "true")
	test(t, "(nil? 0)", "false")
	test(t, "(nil? ())", "false")
	test(t, "(nil? []", "false")
	test(t, "(nil? {}", "false")
	test(t, "(nil? x)", "false")
}

func TestPreludeInvalidIsNil(t *testing.T) {
	test(t, "(nil?)", "  [ERROR]  ")
	test(t, "(nil? () [])", "  [ERROR]  ")
}

func TestPreludeIsTrue(t *testing.T) {
	test(t, "(true? true)", "true")
	test(t, "(true? false)", "false")
	test(t, "(true? nil)", "false")
	test(t, "(true? 0)", "false")
	test(t, "(true? ())", "false")
	test(t, "(true? []", "false")
	test(t, "(true? {}", "false")
	test(t, "(true? x)", "false")
}

func TestPreludeIsFalse(t *testing.T) {
	test(t, "(false? true)", "false")
	test(t, "(false? false)", "true")
	test(t, "(false? nil)", "false")
	test(t, "(false? 0)", "false")
	test(t, "(false? ())", "false")
	test(t, "(false? []", "false")
	test(t, "(false? {}", "false")
	test(t, "(false? x)", "false")
}

func TestPreludeRange(t *testing.T) {
	test(t, "(range 0 5 1)", "(0 1 2 3 4)")
	test(t, "(range 5 9 1)", "(5 6 7 8)")
	test(t, "(range 0 0 1)", "()")
	test(t, "(range 1 0 1)", "()")

	test(t, "(range 0 3 (/ 2))", "(0 0.5 1 1.5 2 2.5)")
}

func TestPreludeMap(t *testing.T) {
	test(t, "(map inc '())", "()")
	test(t, "(map inc '(1 2 3))", "(2 3 4)")
	test(t, "(map inc '(3 2 1))", "(4 3 2)")

	test(t, "(map dec []", "[]")
	test(t, "(map dec [1 2 3])", "[0 1 2]")
	test(t, "(map dec [3 2 1])", "[2 1 0]")

	test(t, "(map (fn [x] (+ x 10)) [1 2 3])", "[11 12 13]")
}

func TestPreludeReduce(t *testing.T) {
	test(t, `(reduce + 0 '())`, "0")
	test(t, `(reduce + 0 '(3 2 1))`, "6")

	test(t, `(reduce + 0 [])`, "0")
	test(t, `(reduce + 0 [3 2 1])`, "6")
}

func TestPreludeSum(t *testing.T) {
	test(t, `(sum [])`, "0")
	test(t, `(sum [4 5 6])`, "15")
	test(t, `(sum (irange 1 5))`, "10")
}

func TestPreludeProd(t *testing.T) {
	test(t, `(prod [])`, "1")
	test(t, `(prod [4 5 6])`, "120")
	test(t, `(prod (irange 1 5))`, "24")
}

func TestJoin(t *testing.T) {
	test(t, `(join)`, `""`)
	test(t, `(join "a")`, `"a"`)
	test(t, `(join "a" "b" "c")`, `"abc"`)
	// test(t, `(join *SPLIS_HOME* "lib/core/math" ".splis")`, `"/home/mathias/go/src/github.com/mhoertnagl/splis2/lib/core/math.splis"`)
}

func TestPreludeFilter(t *testing.T) {
	test(t, `(filter (fn [x] (> x 0)) [])`, "[]")
	test(t, `(filter (fn [x] (> x 0)) [1 2 3])`, "[1 2 3]")
	test(t, `(filter (fn [x] (> x 0)) [0 0 0])`, "[]")
	test(t, `(filter (fn [x] (> x 0)) [(- 5) 3 (- 2) 4 0 7])`, "[3 4 7]")
}

func TestPreludeRemove(t *testing.T) {
	test(t, `(remove (fn [x] (> x 0)) [])`, "[]")
	test(t, `(remove (fn [x] (> x 0)) [1 2 3])`, "[]")
	test(t, `(remove (fn [x] (> x 0)) [(- 1) (- 4) 0])`, "[-1 -4 0]")
	test(t, `(remove (fn [x] (> x 0)) [(- 5) 3 (- 2) 4 0 7])`, "[-5 -2 0]")
}

func TestStr(t *testing.T) {
	test(t, `(str)`, `""`)
	test(t, `(str 42)`, `"42"`)
	test(t, `(str 4 2)`, `"42"`)
	test(t, `(str (+ 4 (+ 5 6)) [1 2] "x")`, `"15[1 2]"x""`)
}

func TestWrite(t *testing.T) {
	test(t, `(write *STDOUT* "Test\n")`, "nil")
	test(t, `(write *STDOUT* "Test\\n")`, "nil")
	test(t, `(write *STDOUT* (str [1 2 3]) (str [4 5 6]) "\n")`, "nil")
}

func TestPreludePrint(t *testing.T) {
	test(t, `(print)`, "nil")
	test(t, `(print "www")`, "nil")
	test(t, `(print "xxx" "\n")`, "nil")
	test(t, `(println "yyy")`, "nil")
	test(t, `(println "zzz" "$$$")`, "nil")
}

func TestPreludeApply(t *testing.T) {
	test(t, `(apply print ["a" "b" "c"])`, "nil")
	test(t, `(apply and [true false true])`, "false")
	test(t, `(apply and [true true true])`, "true")
	test(t, `(apply and [(< 0 1) (= 0 0)])`, "true")
}

func TestAnd(t *testing.T) {
	test(t, `(and)`, "true")
	test(t, `(and false)`, "false")
	test(t, `(and true)`, "true")
	test(t, `(and false false)`, "false")
	test(t, `(and false true)`, "false")
	test(t, `(and true false)`, "false")
	test(t, `(and true true)`, "true")
	test(t, `(and (< 0 1) (= 0 0))`, "true")
}

func TestOr(t *testing.T) {
	test(t, `(or)`, "false")
	test(t, `(or false)`, "false")
	test(t, `(or true)`, "true")
	test(t, `(or false false)`, "false")
	test(t, `(or false true)`, "true")
	test(t, `(or true false)`, "true")
	test(t, `(or true true)`, "true")
	test(t, `(or (< 0 1) (= 0 0))`, "true")
}

func TestQuot(t *testing.T) {
	test(t, `(quot 9 3)`, "3")
	test(t, `(quot 8 3)`, "2")
	test(t, `(quot 7 3)`, "2")
	test(t, `(quot 6 3)`, "2")
	test(t, `(quot 5 3)`, "1")
}

func TestMod(t *testing.T) {
	test(t, `(mod 9 3)`, "0")
	test(t, `(mod 8 3)`, "2")
	test(t, `(mod 7 3)`, "1")
	test(t, `(mod 6 3)`, "0")
	test(t, `(mod 5 3)`, "2")
}

func TestPreludeAll(t *testing.T) {
	test(t, `(all true? [(< 0 1) (= 0 0)])`, "true")
	test(t, `(all (fn [x] (< 0 x)) [1 2 3])`, "true")
	test(t, `(all (fn [x] (> 0 x)) [1 2 3])`, "false")
}

func TestPreludeAny(t *testing.T) {
	test(t, `(any true? [(< 0 1) (= 0 0)])`, "true")
	test(t, `(any (fn [x] (< 2 x)) [1 2 3])`, "true")
	test(t, `(any (fn [x] (> 1 x)) [1 2 3])`, "false")
}

func TestPreludeAllTrue(t *testing.T) {
	test(t, `(any-true [(< 0 1) (= 0 0)])`, "true")
	test(t, `(all-true [(< 0 1) (= 0 0)])`, "true")
}

func test(t *testing.T, i string, e string) {
	t.Helper()
	// fmt.Printf("\n\nTesting %s ...\n", i)
	teste(t, data.NewRootEnv(), i, e)
}

func teste(t *testing.T, env data.Env, i string, e string) {
	t.Helper()
	r := read.NewReader()
	p := read.NewParser()
	w := print.NewPrinter()
	v := eval.NewEvaluator(env)
	fmt.Print(w.PrintErrors(v.Errors()...))
	r.Load(i)
	n := p.Parse(r)
	m := v.Eval(n)
	a := w.Print(m)
	fmt.Print(w.PrintErrors(v.Errors()...))
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}

func testenv(t *testing.T, env data.Env, name string, e string) {
	t.Helper()
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
