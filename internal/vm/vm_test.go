package vm_test

import (
	"hash/fnv"
	"reflect"
	"testing"

	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestConst1(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 42),
	)
}

func TestConst2(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpConst, 43),
	)
	testVal(t, int64(43), m.InspectStack(0))
	testVal(t, int64(42), m.InspectStack(1))
}

func TestPop(t *testing.T) {
	testToS(t, int64(42),
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
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 19),
		vm.Instr(vm.OpConst, 23),
		vm.Instr(vm.OpAdd),
	)
}

func TestSub2(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 43),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
	)
}

func TestMul2(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 21),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
	)
}

func TestDiv2(t *testing.T) {
	testToS(t, int64(21),
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
	testToS(t, int64(0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 10),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestIfTrue1(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 10),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestIfElseFalse1(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestIfElseTrue1(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

// TODO: If, elseif, else

func TestRunAnd1(t *testing.T) {
	testToS(t, false,
		// a
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 10),
		// b
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunAnd2(t *testing.T) {
	testToS(t, false,
		// a
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 10),
		// b
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunAnd3(t *testing.T) {
	testToS(t, false,
		// a
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 10),
		// b
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunAnd4(t *testing.T) {
	testToS(t, true,
		// a
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 10),
		// b
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunOr1(t *testing.T) {
	testToS(t, false,
		// a
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 10),
		// b
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunOr2(t *testing.T) {
	testToS(t, true,
		// a
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 10),
		// b
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunOr3(t *testing.T) {
	testToS(t, true,
		// a
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 10),
		// b
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunOr4(t *testing.T) {
	testToS(t, true,
		// a
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 10),
		// b
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLT1(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpLT),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLT2(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpLT),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLT3(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpLT),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLE1(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpLE),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLE2(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpLE),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunLE3(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ1(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ2(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ3(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunNE1(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpNE),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunNE2(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpNE),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunNE3(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpNE),
		vm.Instr(vm.OpHalt),
	)
}

// TODO: Environments,locals.

func TestLocals1(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpPopEnv),
	)
	testVal(t, int64(4), m.InspectStack(0))
}

func TestLocals2(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetLocal, 1),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpSetLocal, 2),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpGetLocal, 1),
		vm.Instr(vm.OpGetLocal, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpPopEnv),
	)
	testVal(t, int64(6), m.InspectStack(0))
}

func TestLocals3(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpPopEnv),
	)
	testVal(t, int64(6), m.InspectStack(0))
}

func TestRunLet1(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
	)
	testVal(t, int64(4), m.InspectStack(0))
}

func TestRunDef1(t *testing.T) {
	testToS(t, int64(2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpGetGlobal, 0),
	)
}

func TestRunIf1(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestRunIf2(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

// TODO: Locals with shadowing.

func TestRunCreateVector1(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3)}
	testToS(t, e,
		vm.Instr(vm.OpEmptyVector),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpCons),
	)
}

func TestRunHalt(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestRunFunCall(t *testing.T) {
	testToS(t, int64(3),
		vm.Instr(vm.OpConst, 1),
		// vm.Instr(vm.OpDebug, 1),
		// -| 1
		vm.Instr(vm.OpConst, 39),
		// vm.Instr(vm.OpDebug, 1),
		// -| 1 @30
		// call fn (1) => 1 + 1
		vm.Instr(vm.OpCall),
		// -| 2
		// -| 2 1
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpAdd),
		// -| 3
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x x))
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpNewEnv),
		// Stack contains arguments in reverse order.
		vm.Instr(vm.OpSetLocal, 0),
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpGetLocal, 0),
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpDebug, vm.DbgStack),
		vm.Instr(vm.OpPopEnv),
		// -| 2*x
		vm.Instr(vm.OpReturn),
	)
}

// testToS executes a sequence of instructions in the vm and tests the top of
// the stack element against an expected value. Will raise an error if the
// types or the values are unequal. The stack is fixed to a maximum size of
// 1024 cells.
func testToS(t *testing.T, expected vm.Val, c ...vm.Ins) {
	t.Helper()
	m := testRun(t, c...)
	testVal(t, expected, m.InspectStack(0))
	if m.StackSize() != 1 {
		t.Errorf("Stack size should be [%v] but is [%v].", 1, m.StackSize())
	}
}

// testRun executes a new VM instance with the code provided and returns the
// VM thereafter.
func testRun(t *testing.T, c ...vm.Ins) vm.VM {
	t.Helper()
	m := vm.New(1024, 512, 256, 128)
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
	if reflect.DeepEqual(expected, actual) == false {
		// if expected != actual {
		t.Errorf("Expected [%v] but got [%v].", expected, actual)
	}
}

func hash(sym string) uint64 {
	hg := fnv.New64()
	hg.Reset()
	hg.Write([]byte(sym))
	return hg.Sum64()
}
