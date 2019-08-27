package print_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/data"
	"github.com/mhoertnagl/splis2/internal/print"
)

func TestError(t *testing.T) {
	test(t, data.NewError(""), "  [ERROR]  ")
}

func TestNil(t *testing.T) {
	test(t, nil, "nil")
}

func TestBool(t *testing.T) {
	test(t, true, "true")
	test(t, false, "false")
}

func TestNumbers(t *testing.T) {
	test(t, 1., "1")
	test(t, .2, "0.2")
	test(t, -.2, "-0.2")
}

func TestStrings(t *testing.T) {
	test(t, "", `""`)
	test(t, "x", `"x"`)
}

func TestSymbols(t *testing.T) {
	test(t, data.NewSymbol(""), ``)
	test(t, data.NewSymbol("x"), `x`)
	test(t, data.NewSymbol("+"), `+`)
	test(t, data.NewSymbol("<<"), `<<`)
}

func TestLists(t *testing.T) {
	test(t, data.NewList2(), `()`)
	test(t, data.NewList2(1.), `(1)`)
	test(t, data.NewList2(1., ""), `(1 "")`)
}

func TestVector(t *testing.T) {
	test(t, data.NewVector2(), `[]`)
	test(t, data.NewVector2(1.), `[1]`)
	test(t, data.NewVector2(1., ""), `[1 ""]`)
}

func TestHashMap(t *testing.T) {
	m := data.NewHashMap2()
	test(t, m, `{}`)
	m.Items[data.NewSymbol("x")] = 1.
	test(t, m, `{x 1}`)
	// m.Items[data.NewSymbol("y")] = data.NewString("42")
	// TODO: Order is not guaranteed.
	// test(t, m, `{x 1 y "42"}`)
}

func test(t *testing.T, n data.Node, e string) {
	p := print.NewPrinter()
	a := p.Print(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
