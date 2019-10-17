package vm_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestInstr(t *testing.T) {
	c := vm.ConcatVar(
		vm.Instr(vm.OpConst, 1),
	)

}
