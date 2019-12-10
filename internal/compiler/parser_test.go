package compiler_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/compiler"
)

func TestParseNil(t *testing.T) {
	testpw(t, " nil ", "nil")
}

func TestParseBool(t *testing.T) {
	testpw(t, " true ", "true")
	testpw(t, " false ", "false")
}

func TestParseNumbers(t *testing.T) {
	testpw(t, " 0 ", "0")
	testpw(t, " 1 ", "1")
	// testpw(t, " 0.3 ", "0.3")
	// testpw(t, " -.3 ", "-0.3")
}

func TestParseInvalidNumbers(t *testing.T) {
	testpw(t, " 1# ", "  [ERROR]  ")
}

func TestParseStrings(t *testing.T) {
	testpw(t, ` "" `, `""`)
	testpw(t, ` "x" `, `"x"`)
}

func TestParseIncomleteStrings(t *testing.T) {
	testpw(t, ` "x `, `  [ERROR]  `)
}

func TestParseSymbol(t *testing.T) {
	testpw(t, " + ", "+")
}

func TestParseLists(t *testing.T) {
	testpw(t, " (  ) ", "()")
	testpw(t, " (+  ) ", "(+)")
	testpw(t, " (+ 1  2  3  ) ", "(+ 1 2 3)")
	testpw(t, " (flatten   ( 1 2 )  ( 3)  ) ", "(flatten (1 2) (3))")
}

func TestParseIncompleteLists(t *testing.T) {
	testpw(t, " ( ", "()")
	testpw(t, " (+ ", "(+)")
	testpw(t, " (+ 1 ", "(+ 1)")
}

func TestParseVectors(t *testing.T) {
	testpw(t, " [    ] ", "[]")
	testpw(t, " [+  ] ", "[+]")
	testpw(t, " [+  -    *  /  ] ", "[+ - * /]")
	testpw(t, ` [   [ "a""b" ]  [ "c"]  ] `, `[["a" "b"] ["c"]]`)
}

func TestParseIncompleteVectors(t *testing.T) {
	testpw(t, " [ ", "[]")
	testpw(t, " [+ ", "[+]")
	testpw(t, " [+ 1 ", "[+ 1]")
}

func TestParseUnexpectedSeqClosingTag(t *testing.T) {
	testpw(t, " ) ", "  [ERROR]  ")
	testpw(t, " ] ", "  [ERROR]  ")
	testpw(t, " } ", "  [ERROR]  ")
}

func TestParseHashMaps(t *testing.T) {
	testpw(t, " {   } ", "{}")
	testpw(t, ` { "a" 1 } `, `{"a" 1}`)
	// TODO: Order is not guaranteed.
	// testpw(t, ` { "a" 1"b"  2 } `, `{"a" 1 "b" 2}`)
}

func TestParseIncompleteHashMaps(t *testing.T) {
	testpw(t, " { ", "{}")
	testpw(t, ` {"a" `, `{"a" }`) // The value for key "a" is the empty string.
	testpw(t, ` {"a" 1 `, `{"a" 1}`)
}

func TestParseQuote(t *testing.T) {
	testpw(t, " '42 ", "(quote 42)")
	testpw(t, ` '"x" `, `(quote "x")`)
	testpw(t, " '() ", "(quote ())")
	testpw(t, " '(+ 1 1) ", "(quote (+ 1 1))")
	testpw(t, " '[1 2 3] ", "(quote [1 2 3])")
	testpw(t, ` '{ "a" 1 } `, `(quote {"a" 1})`)
}

func TestParseQuasiquote(t *testing.T) {
	testpw(t, " `42 ", "(quasiquote 42)")
	testpw(t, " `\"x\" ", `(quasiquote "x")`)
	testpw(t, " `() ", "(quasiquote ())")
	testpw(t, " `(+ 1 1) ", "(quasiquote (+ 1 1))")
	testpw(t, " `[1 2 3] ", "(quasiquote [1 2 3])")
	testpw(t, " `{ \"a\" 1 } ", `(quasiquote {"a" 1})`)
}

func TestParseQuasiquoteUnquote(t *testing.T) {
	testpw(t, " `~42 ", "(quasiquote (unquote 42))")
	testpw(t, " `(+ ~a ~b) ", "(quasiquote (+ (unquote a) (unquote b)))")
}

func TestParseQuasiquoteSpliceUnquote(t *testing.T) {
	testpw(t, " `~@(42) ", "(quasiquote (splice-unquote (42)))")
	testpw(t, " `(+ ~@(a b) c) ", "(quasiquote (+ (splice-unquote (a b)) c))")
}

func testpw(t *testing.T, i string, e string) {
	r := compiler.NewReader()
	r.Load(i)
	p := compiler.NewParser()
	n := p.Parse(r)
	w := compiler.NewPrinter()
	a := w.Print(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
