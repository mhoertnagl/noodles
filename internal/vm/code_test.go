package vm_test

import (
	"bytes"
	"testing"

	"github.com/mhoertnagl/noodles/internal/vm"
)

func TestInstr(t *testing.T) {
	c := vm.ConcatVar(
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
	)
	e := []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	if bytes.Compare(c, e) != 0 {
		t.Errorf("Expecting %v but got %v.", e, c)
	}
}

func TestStr(t *testing.T) {
	es := "Hello, World!"
	a := vm.Str(es)
	e := []byte{6, 0, 0, 0, 0, 0, 0, 0, 13, 72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}
	if bytes.Compare(a, e) != 0 {
		t.Errorf("Expecting %v but got %v.", e, a)
	}
	as := string(a[9:])
	if as != es {
		t.Errorf("Expecting [%s] but got [%s].", es, as)
	}
}
