package cmp_test

import (
	"testing"

	"github.com/mhoertnagl/noodles/internal/cmp"
)

func TestSymIndexOf0(t *testing.T) {
	sym := cmp.NewSymTable()
	sym.AddVar("x", "y")
	testIndexOf(t, sym, "x", 0)
	testIndexOf(t, sym, "y", 1)
}

func TestSymIndexOf1(t *testing.T) {
	sym := cmp.NewSymTable()
	sym.AddVar("a", "b", "c", "lst", "min", "f")
	sub := sym.NewSymTable()
	sub.AddVar("x")
	testIndexOf(t, sub, "x", 0)
	testIndexOf(t, sub, "min", -4)
}

func testIndexOf(t *testing.T, sym *cmp.SymTable, n string, e int) {
	t.Helper()
	a, _ := sym.IndexOf(n)
	if a != e {
		t.Errorf("Symbol [%s]: Expected [%d] but got [%d]\n", n, e, a)
	}
}
