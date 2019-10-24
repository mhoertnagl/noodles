package vm_test

import (
	"reflect"
	"testing"

	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestConst1(t *testing.T) {
	testToS(t, uint64(42),
		vm.Instr(vm.OpConst, 42),
	)
}

func TestConst2(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpConst, 43),
	)
	testVal(t, m.Inspect(0), uint64(43))
	testVal(t, m.Inspect(1), uint64(42))
}

func TestPop(t *testing.T) {
	testToS(t, uint64(42),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpPop),
	)
}

// func TestAdd0(t *testing.T) {
// 	testToS(t, uint64(1),
// 		vm.Instr(vm.OpStop),
// 		vm.Instr(vm.OpAdd),
// 	)
// }
//
// func TestAdd1(t *testing.T) {
// 	testToS(t, uint64(2),
// 		vm.Instr(vm.OpStop),
// 		vm.Instr(vm.OpConst, 1),
// 		vm.Instr(vm.OpAdd),
// 	)
// }

func TestAdd2(t *testing.T) {
	testToS(t, uint64(42),
		vm.Instr(vm.OpConst, 19),
		vm.Instr(vm.OpConst, 23),
		vm.Instr(vm.OpAdd),
	)
}

func TestSub2(t *testing.T) {
	testToS(t, uint64(42),
		vm.Instr(vm.OpConst, 43),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
	)
}

func TestMul2(t *testing.T) {
	testToS(t, uint64(42),
		vm.Instr(vm.OpConst, 21),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
	)
}

func TestDiv2(t *testing.T) {
	testToS(t, uint64(21),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpDiv),
	)
}

func TestFalse(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
	)
}

func TestTrue(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
	)
}

func TestIfFalse1(t *testing.T) {
	testToS(t, uint64(0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfFalse, 10),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestIfTrue1(t *testing.T) {
	testToS(t, uint64(1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfFalse, 10),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
	)
}

// TODO: If with else.
// TODO: If, elseif, else
// TODO: ||
// TODO: &&

// TODO: Environments,locals.

// testToS executes a sequence of instructions in the vm and tests the top of
// the stack element against an expected value. Will raise an error if the
// types or the values are unequal. The stack is fixed to a maximum size of
// 1024 cells.
func testToS(t *testing.T, expected vm.Val, c ...vm.Ins) {
	t.Helper()
	m := testRun(t, c...)
	testVal(t, expected, m.Inspect(0))
	if m.StackSize() != 1 {
		t.Errorf("Stack size should be [%v] but is [%v].", 1, m.StackSize())
	}
}

// testRun executes a new VM instance with the code provided and returns the
// VM thereafter.
func testRun(t *testing.T, c ...vm.Ins) vm.VM {
	t.Helper()
	m := vm.New(1024)
	m.Run(vm.Concat(c))
	return m
}

// testVal compares the expected and the actual values for equal types and
// values. It will raise an error otherwise.
func testVal(t *testing.T, expected vm.Val, actual vm.Val) {
	t.Helper()
	aType := reflect.TypeOf(actual)
	eType := reflect.TypeOf(expected)
	if aType != eType {
		t.Errorf("Expected [%v] but got [%v].", eType, aType)
	}
	if expected != actual {
		t.Errorf("Expected [%v] but got [%v].", expected, actual)
	}
}
