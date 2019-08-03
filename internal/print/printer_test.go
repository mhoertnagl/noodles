package print_test

import (
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"testing"
)

func TestError(t *testing.T) {
	test(t, read.NewError(""), "  [ERROR]  ")
}

func TestNil(t *testing.T) {
	test(t, read.NilObject, "nil")
}

func TestBool(t *testing.T) {
	test(t, read.TrueObject, "true")
	test(t, read.FalseObject, "false")
}

func TestNumbers(t *testing.T) {
	test(t, read.NewNumber(1), "1")
	test(t, read.NewNumber(.2), "0.2")
	test(t, read.NewNumber(-.2), "-0.2")
}

func TestStrings(t *testing.T) {
	test(t, read.NewString(""), `""`)
	test(t, read.NewString("x"), `"x"`)
}

func TestSymbols(t *testing.T) {
	test(t, read.NewSymbol(""), ``)
	test(t, read.NewSymbol("x"), `x`)
	test(t, read.NewSymbol("+"), `+`)
	test(t, read.NewSymbol("<<"), `<<`)
}

func TestLists(t *testing.T) {
	test(t, read.NewList2(), `()`)
	test(t, read.NewList2(read.NewNumber(1)), `(1)`)
	test(t, read.NewList2(read.NewNumber(1), read.NewString("")), `(1 "")`)
}

func TestVector(t *testing.T) {
	test(t, read.NewVector2(), `[]`)
	test(t, read.NewVector2(read.NewNumber(1)), `[1]`)
	test(t, read.NewVector2(read.NewNumber(1), read.NewString("")), `[1 ""]`)
}

func TestHashMap(t *testing.T) {
	m := read.NewHashMap2()
	test(t, m, `{}`)
	m.Items[read.NewSymbol("x")] = read.NewNumber(1)
	test(t, m, `{x 1}`)
	// m.Items[read.NewSymbol("y")] = read.NewString("42")
	// TODO: Order is not guaranteed.
	// test(t, m, `{x 1 y "42"}`)
}

func test(t *testing.T, n read.Node, e string) {
	p := print.NewPrinter()
	a := p.Print(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
