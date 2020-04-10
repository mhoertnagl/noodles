package vm_test

import (
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
// 		vm.Instr(vm.OpEnd),
// 		vm.Instr(vm.OpAdd),
// 	)
// }
//
// func TestAdd1(t *testing.T) {
// 	testToS(t, uint64(2),
// 		vm.Instr(vm.OpEnd),
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

func TestRunNot1(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpNot),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunNot2(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpNot),
		vm.Instr(vm.OpHalt),
	)
}

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

// func TestRunBitAnd(t *testing.T) {
// 	testToS(t, int64(8),
// 		vm.Instr(vm.OpEnd),
// 		vm.Instr(vm.OpConst, 12),
// 		vm.Instr(vm.OpConst, 10),
// 		vm.Instr(vm.OpAnd),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitOr(t *testing.T) {
// 	testToS(t, int64(14),
// 		vm.Instr(vm.OpEnd),
// 		vm.Instr(vm.OpConst, 12),
// 		vm.Instr(vm.OpConst, 10),
// 		vm.Instr(vm.OpOr),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitInv(t *testing.T) {
// 	testToS(t, int64(^2),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpInv),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitInv2(t *testing.T) {
// 	testToS(t, int64(^-3),
// 		vm.Instr(vm.OpConst, 0),
// 		vm.Instr(vm.OpConst, 3),
// 		vm.Instr(vm.OpSub),
// 		vm.Instr(vm.OpInv),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitShiftLeft(t *testing.T) {
// 	testToS(t, int64(32),
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSll),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitShiftRight(t *testing.T) {
// 	testToS(t, int64(2),
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSrl),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestRunBitShiftRightSigned(t *testing.T) {
// 	testToS(t, int64(-2),
// 		vm.Instr(vm.OpConst, 0),
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpSub),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSra),
// 		vm.Instr(vm.OpHalt),
// 	)
// }

// TODO: Environments,locals.

func TestRunLet1(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpDropArgs, 1),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(4), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunLet2(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpDropArgs, 2),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunLet31(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(1), m.InspectFrames(0))
	testVal(t, int64(2), m.InspectFrames(1))
	testVal(t, int64(1), m.InspectStack(0))
	testVal(t, int64(2), m.InspectStack(1))
}

func TestRunLet3(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpDropArgs, 2),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(1), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunLet4(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpDropArgs, 1),
		vm.Instr(vm.OpDropArgs, 1),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
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

func TestRunIf3(t *testing.T) {
	testToS(t, int64(21),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 21),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunIf4(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 21),
		vm.Instr(vm.OpHalt),
	)
}

// TODO: Locals with shadowing.
//       Does not work for the current implementation of let bindings.

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

func TestRunCreateVector2(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorHead(t *testing.T) {
	e := int64(1)
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpFst),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorTail1(t *testing.T) {
	e := []vm.Val{}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpRest),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorTail2(t *testing.T) {
	e := []vm.Val{int64(2), int64(3)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpRest),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorLength1(t *testing.T) {
	e := int64(0)
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpLength),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorLength2(t *testing.T) {
	e := int64(3)
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpLength),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunDissolve(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpDissolve),
	)
	testVal(t, int64(1), m.InspectStack(0))
	testVal(t, int64(2), m.InspectStack(1))
	testVal(t, int64(3), m.InspectStack(2))
	testVal(t, nil, m.InspectStack(3))
}

func TestRunHalt(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestRunAnonymousFun2(t *testing.T) {
	m := testRun(t,
		// ((fn ...) 1)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 21),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
	testVal(t, int64(2), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunAnonymousFun3(t *testing.T) {
	m := testRun(t,
		// ((fn ...) 1)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 31),
		vm.Instr(vm.OpCall),
		// (+ ((fn ...) 1) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunAnonymousFun4(t *testing.T) {
	m := testRun(t,
		// (((fn ...)) 1)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpRef, 53),
		// Call the 0-adic function that returns the 1-adic function.
		vm.Instr(vm.OpCall),
		// Call the 1-adic function.
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		// (fn [] ...)
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpRef, 23),
		vm.Instr(vm.OpReturn),
	)
	testVal(t, int64(2), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunLeafFunDef(t *testing.T) {
	m := testRun(t,
		// (def inc (fn ...))
		vm.Instr(vm.OpRef, 49),
		vm.Instr(vm.OpSetGlobal, 0),
		// (inc 1)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		// (+ (inc ...) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunVariadicFun(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3), int64(4)}
	m := testRun(t,
		// ((fn ...) 1 2 3 4)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 48),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		// (fn [x & xs] (:: x xs))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpReturn),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunFac5(t *testing.T) {
	testToS(t, int64(120),
		vm.Instr(vm.OpRef, 57),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 40),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpReturn),
	)
}

func TestRunStringConst(t *testing.T) {
	testToS(t, "Hello, World!",
		vm.Str("Hello, World!"),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunTest0Bin(t *testing.T) {
	testToS(t, int64(6),
		vm.Instr(vm.OpJump, 0),

		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
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
func testRun(t *testing.T, c ...vm.Ins) *vm.VM {
	t.Helper()
	m := vm.NewVM(1024, 512, 512)
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
		t.Errorf("Expected [%v] but got [%v].", expected, actual)
	}
}
