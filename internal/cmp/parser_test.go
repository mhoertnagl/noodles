package cmp_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/cmp"
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
	testpw(t, " -1 ", "-1")
	testpw(t, " 0.3 ", "0.3")
	testpw(t, " -0.3 ", "-0.3")

	testpi(t, " 0 ", 0)
	testpi(t, " 1 ", 1)
	testpi(t, " -1 ", -1)

	testpf(t, " 0.3 ", 0.3)
	testpf(t, " -0.3 ", -0.3)
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

func TestParseQuoteUnquote(t *testing.T) {
	testpw(t, " '~42 ", "(quote (unquote 42))")
	testpw(t, " '(+ ~a ~b) ", "(quote (+ (unquote a) (unquote b)))")
}

func TestParseQuoteSpliceUnquote(t *testing.T) {
	testpw(t, " '~@(42) ", "(quote (unquote (dissolve (42))))")
	testpw(t, " '(+ ~@(a b) c) ", "(quote (+ (unquote (dissolve (a b))) c))")
}

func testpw(t *testing.T, i string, e string) {
	r := cmp.NewReader()
	r.Load(i)
	p := cmp.NewParser()
	n := p.Parse(r)
	a := cmp.PrintAst(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}

func testpi(t *testing.T, i string, e int64) {
	r := cmp.NewReader()
	r.Load(i)
	p := cmp.NewParser()
	n := p.Parse(r)
	if a, ok := n.(int64); !ok || a != e {
		t.Errorf("Expecting [%d] but got [%d]", e, a)
	}
}

func testpf(t *testing.T, i string, e float64) {
	r := cmp.NewReader()
	r.Load(i)
	p := cmp.NewParser()
	n := p.Parse(r)
	if a, ok := n.(float64); !ok || a != e {
		t.Errorf("Expecting [%v] but got [%v]", e, a)
	}
}
