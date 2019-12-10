package compiler_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/compiler"
	"github.com/mhoertnagl/splis2/internal/print"
)

// func TestPrintError(t *testing.T) {
// 	testPrint(t, compiler.NewError(""), "  [ERROR]  ")
// }
//
// func TestPrintNil(t *testing.T) {
// 	testPrint(t, nil, "nil")
// }
//
// func TestPrintBool(t *testing.T) {
// 	testPrint(t, true, "true")
// 	testPrint(t, false, "false")
// }
//
// func TestPrintNumbers(t *testing.T) {
// 	testPrint(t, 1., "1")
// 	testPrint(t, .2, "0.2")
// 	testPrint(t, -.2, "-0.2")
// }
//
// func TestPrintStrings(t *testing.T) {
// 	testPrint(t, "", `""`)
// 	testPrint(t, "x", `"x"`)
// }
//
// func TestPrintSymbols(t *testing.T) {
// 	testPrint(t, compiler.NewSymbol(""), ``)
// 	testPrint(t, compiler.NewSymbol("x"), `x`)
// 	testPrint(t, compiler.NewSymbol("+"), `+`)
// 	testPrint(t, compiler.NewSymbol("<<"), `<<`)
// }
//
// func TestPrintLists(t *testing.T) {
// 	testPrint(t, compiler.NewList2(), `()`)
// 	testPrint(t, compiler.NewList2(1.), `(1)`)
// 	testPrint(t, compiler.NewList2(1., ""), `(1 "")`)
// }
//
// func TestPrintVector(t *testing.T) {
// 	testPrint(t, compiler.NewVector2(), `[]`)
// 	testPrint(t, compiler.NewVector2(1.), `[1]`)
// 	testPrint(t, compiler.NewVector2(1., ""), `[1 ""]`)
// }
//
// func TestPrintHashMap(t *testing.T) {
// 	m := compiler.NewEmptyHashMap()
// 	testPrint(t, m, `{}`)
// 	m.Items["x"] = 1.
// 	testPrint(t, m, `{"x" 1}`)
// 	// m.Items[NewSymbol("y")] = NewString("42")
// 	// TODO: Order is not guaranteed.
// 	// test(t, m, `{x 1 y "42"}`)
// }

func testPrint(t *testing.T, n compiler.Node, e string) {
	t.Helper()
	p := print.NewPrinter()
	a := p.Print(n)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
