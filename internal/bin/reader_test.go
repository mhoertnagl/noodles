package bin_test

import (
	"os"
	"testing"

	"github.com/mhoertnagl/splis2/internal/bin"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestReadBin(t *testing.T) {
	exp := vm.ConcatVar(
		vm.Instr(vm.OpJump, 73),

		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),

		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),

		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetGlobal, hash("inc")),
		vm.Instr(vm.OpCall),
	)

	outFile, err := os.Create("../../test/test-0.splis.bin")
	if err != nil {
		t.Error(err)
	}
	bin.WriteStatic(exp, outFile)

	inFile, err := os.Open("../../test/test-0.splis.bin")
	if err != nil {
		t.Error(err)
	}
	act := bin.ReadStatic(inFile)

	assertCodeEqual(t, act, exp)
}
