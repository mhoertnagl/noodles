package read_test

import (
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"testing"
)

func TestParseNil(t *testing.T) {
	testp(t, " nil ", read.NilObject)
}

func TestParseBool(t *testing.T) {
	testp(t, " true ", read.TrueObject)
	testp(t, " false ", read.FalseObject)
}

func TestParseNumbers(t *testing.T) {
	testpw(t, " 1 ", "1")
	testpw(t, " 0.3 ", "0.3")
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

func testp(t *testing.T, i string, e read.Node) {
	r := read.NewReader()
	r.Load(i)
	p := read.NewParser()
	n := p.Parse(r)
	if n != e {
		w := print.NewPrinter()
		t.Errorf("Expecting [%s] but got [%s]", w.Print(e), w.Print(n))
	}
}

func testpw(t *testing.T, i string, e string) {
	r := read.NewReader()
	r.Load(i)
	p := read.NewParser()
	n := p.Parse(r)
	w := print.NewPrinter()
	a := w.Print(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
