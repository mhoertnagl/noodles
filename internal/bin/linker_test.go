package bin_test

import (
	"reflect"
	"testing"

	"github.com/mhoertnagl/splis2/internal/bin"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestLinkLib(t *testing.T) {
	lb1 := bin.NewLib()
	lb1.Code = vm.ConcatVar(
		vm.Instr(vm.OpRef, 0),
	)
	lb1.Entries = []uint64{0}
	lb2 := bin.NewLib()
	lb2.Code = vm.ConcatVar(
		vm.Instr(vm.OpRef, 0),
	)
	lb2.Entries = []uint64{0}
	exp := bin.NewLib()
	lnk := bin.NewLinker()
	lnk.Add(lb1)
	lnk.Add(lb2)
	assertLibsEqual(t, lnk.Lib(), exp)
}

func assertLibsEqual(t *testing.T, a *bin.Lib, e *bin.Lib) {
	t.Helper()
	if reflect.DeepEqual(a.Code, e.Code) == false {
		t.Errorf(
			"\nCode mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Code,
			e.Code,
		)
	}
	if reflect.DeepEqual(a.Entries, e.Entries) == false {
		t.Errorf(
			"\nEntries mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Entries,
			e.Entries,
		)
	}
	if reflect.DeepEqual(a.Fns, e.Fns) == false {
		t.Errorf(
			"\nFns mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Fns,
			e.Fns,
		)
	}
	if reflect.DeepEqual(a.Macros, e.Macros) == false {
		t.Errorf(
			"\nMacros mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Macros,
			e.Macros,
		)
	}
}
