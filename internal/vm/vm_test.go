package vm_test

import (
	"math"
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

func TestFConst1(t *testing.T) {
	testToS(t, 3.1415,
		vm.Instr(vm.OpConstF, math.Float64bits(3.1415)),
	)
}

func TestPop(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpPop),
	)
}

func TestAdd0(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpAdd),
	)
}

func TestAdd1(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
	)
}

func TestAddII(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 19),
		vm.Instr(vm.OpConst, 23),
		vm.Instr(vm.OpAdd),
	)
}

func TestAddIF(t *testing.T) {
	testToS(t, 4.1,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConstF, math.Float64bits(3.1)),
		vm.Instr(vm.OpAdd),
	)
}

func TestAddFI(t *testing.T) {
	testToS(t, 4.1,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConstF, math.Float64bits(3.1)),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
	)
}

func TestAddFF(t *testing.T) {
	testToS(t, 6.2,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConstF, math.Float64bits(3.1)),
		vm.Instr(vm.OpConstF, math.Float64bits(3.1)),
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

func TestMul0(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpMul),
	)
}

func TestMul1(t *testing.T) {
	testToS(t, int64(2),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
	)
}

func TestMulII(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 21),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
	)
}

func TestDiv(t *testing.T) {
	testToS(t, 3.5,
		vm.Instr(vm.OpConst, 7),
		vm.Instr(vm.OpConstF, math.Float64bits(2.0)),
		vm.Instr(vm.OpDiv),
	)
}

func TestDivFF(t *testing.T) {
	testToS(t, 2.0,
		vm.Instr(vm.OpConstF, math.Float64bits(6.2)),
		vm.Instr(vm.OpConstF, math.Float64bits(3.1)),
		vm.Instr(vm.OpDiv),
	)
}

func TestMod(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpConst, 7),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMod),
	)
}

// --- BOOL ---

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

// --- AND ---

func TestRunAnd00(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
}

func TestRunAnd01(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
}

func TestRunAnd10(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
}

func TestRunAnd11(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
}

func TestRunAnd010(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpFalse),
	)
}

func TestRunAnd111(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpFalse),
	)
}

// --- OR ---

func TestRunOr00(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
}

func TestRunOr01(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
}

func TestRunOr10(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
}

func TestRunOr11(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
}

func TestRunOr010(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpTrue),
	)
}

func TestRunOr000(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpTrue),
	)
}

// --- COMPARISON ---

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

func TestRunEQ0(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpTrue),
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

func TestRunEQ4(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ5(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ6(t *testing.T) {
	testToS(t, false,
		vm.Str("xxx"),
		vm.Str("yyy"),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ7(t *testing.T) {
	testToS(t, true,
		vm.Str("xxx"),
		vm.Str("xxx"),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ8(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ9(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ10(t *testing.T) {
	testToS(t, true,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunEQ11(t *testing.T) {
	testToS(t, false,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpList),
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

// --- SET ---

func TestRunSet1(t *testing.T) {
	testToS(t, int64(4),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
	)
}

func TestRunSet2(t *testing.T) {
	testToS(t, int64(5),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpJump, 58),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 27),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
}

// --- LET ---

func TestRunLet1(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpEnd),
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
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpDropArgs, 2),
		vm.Instr(vm.OpHalt),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
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

// func TestRunLet5(t *testing.T) {
// 	m := testRun(t,
// 		vm.Instr(vm.OpJump, 96),
// 		vm.Instr(vm.OpPushArgs, 2),
// 		vm.Instr(vm.OpPop),
// 		vm.Instr(vm.OpGetArg, 1),
// 		vm.Instr(vm.OpConst, 0),
// 		vm.Instr(vm.OpEQ),
// 		vm.Instr(vm.OpJumpIfNot, 65),
// 		vm.Instr(vm.OpConst, 1),
// 		vm.Instr(vm.OpJump, 95),
// 		vm.Instr(vm.OpEnd),
// 		vm.Instr(vm.OpGetArg, 1),
// 		vm.Instr(vm.OpConst, 1),
// 		vm.Instr(vm.OpSub),
// 		vm.Instr(vm.OpGetArg, 0),
// 		vm.Instr(vm.OpCall),
// 		vm.Instr(vm.OpReturn),
// 		vm.Instr(vm.OpGetArg, 0),
// 		vm.Instr(vm.OpRef, 1, 9),
// 		vm.Instr(vm.OpPushArgs, 1),
// 		vm.Instr(vm.OpEnd),
// 		vm.Instr(vm.OpConst, 1),
// 		vm.Instr(vm.OpGetArg, 0),
// 		vm.Instr(vm.OpCall),
// 		vm.Instr(vm.OpDropArgs, 1),
// 	)
// 	testVal(t, int64(3), m.InspectStack(0))
// 	testVal(t, nil, m.InspectStack(1))
// 	testVal(t, nil, m.InspectFrames(0))
// }

// --- DEF ---

func TestRunDef1(t *testing.T) {
	testToS(t, int64(2),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpGetGlobal, 0),
	)
}

// --- IF ---

func TestRunIf11(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 19),
		vm.Instr(vm.OpConst, 1),
	)
}

func TestRunIf21(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 37),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestRunIf22(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 37),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestRunIf3(t *testing.T) {
	testToS(t, int64(21),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
}

func TestRunIf4(t *testing.T) {
	testToS(t, int64(42),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
}

// --- COND ---

func TestRunCond1(t *testing.T) {
	testToS(t, int64(1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 56),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 75),
		vm.Instr(vm.OpConst, 3),
	)
}

func TestRunCond2(t *testing.T) {
	testToS(t, int64(2),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 56),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 75),
		vm.Instr(vm.OpConst, 3),
	)
}

func TestRunCond3(t *testing.T) {
	testToS(t, int64(3),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 56),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 75),
		vm.Instr(vm.OpConst, 3),
	)
}

// TODO: Locals with shadowing.

// --- VECTORS ---

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

func TestRunCreateVector3(t *testing.T) {
	e := []vm.Val{int64(0), int64(1), int64(2), int64(3)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpCons),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorAppend(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3), int64(4)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpAppend),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorConcat0(t *testing.T) {
	e := []vm.Val{}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConcat),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorConcat(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpConcat),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorHead(t *testing.T) {
	e := int64(1)
	m := testRun(t,
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpNth),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorTail1(t *testing.T) {
	e := []vm.Val{}
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpDrop),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
}

func TestRunVectorTail2(t *testing.T) {
	e := []vm.Val{int64(2), int64(3)}
	m := testRun(t,
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpDrop),
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

// --- HALT ---

func TestRunHalt(t *testing.T) {
	testToS(t, int64(0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpConst, 1),
	)
}

// --- FN ---

func TestRunAnonymousFun2(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 50),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 19),
		vm.Instr(vm.OpCall),
	)
	testVal(t, int64(2), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunAnonymousFun3(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 60),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 29),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpAdd),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunAnonymousFun4(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpJump, 79),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 61),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 30),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 20),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
	testVal(t, int64(2), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunAnonymousNestedFun(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpJump, 106),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpJump, 78),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 48),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 19),
		vm.Instr(vm.OpCall),
	)
	testVal(t, int64(2), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunLeafFunDef(t *testing.T) {
	m := testRun(t,
		vm.Instr(vm.OpJump, 40),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpAdd),
	)
	testVal(t, int64(3), m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunVariadicFun1(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3), int64(4)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpPushArgs, 0),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 46),
		vm.Instr(vm.OpCall),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunVariadicFun2(t *testing.T) {
	e := []vm.Val{int64(1), int64(2), int64(3), int64(4)}
	m := testRun(t,
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 85),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 46),
		vm.Instr(vm.OpCall),
	)
	testVal(t, e, m.InspectStack(0))
	testVal(t, nil, m.InspectStack(1))
	testVal(t, nil, m.InspectFrames(0))
}

func TestRunFac(t *testing.T) {
	testToS(t, int64(120),
		vm.Instr(vm.OpJump, 107),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 65),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 106),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
}

func TestRunTailFac(t *testing.T) {
	testToS(t, int64(120),
		vm.Instr(vm.OpJump, 116),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 65),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpJump, 115),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpRecCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpJump, 191),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 151),
		vm.Instr(vm.OpSetGlobal, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 1),
		vm.Instr(vm.OpCall),
	)
}

func TestRunClosure1(t *testing.T) {
	testToS(t, int64(6),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 48),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 28),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
}

func TestRunClosure2(t *testing.T) {
	testToS(t, int64(7),
		vm.Instr(vm.OpJump, 115),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpJump, 87),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpRef, 1, 56),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
}

func TestRunClosure(t *testing.T) {
	testToS(t, int64(3),
		vm.Instr(vm.OpJump, 85),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 58),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 28),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 9),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
}

// --- STR ---

func TestRunStringConst(t *testing.T) {
	testToS(t, "Hello, World!",
		vm.Str("Hello, World!"),
		vm.Instr(vm.OpHalt),
	)
}

func TestRunStringJoin(t *testing.T) {
	testToS(t, "Hello, World!",
		vm.Instr(vm.OpEnd),
		vm.Str("World!"),
		vm.Str("Hello, "),
		vm.Instr(vm.OpJoin),
	)
}

func TestRunStringExplode(t *testing.T) {
	testToS(t, []vm.Val{"x", "y", "z"},
		vm.Str("xyz"),
		vm.Instr(vm.OpExplode),
	)
}

func TestRunStringExplodeAndJoin(t *testing.T) {
	testToS(t, "xyz",
		vm.Instr(vm.OpEnd),
		vm.Str("xyz"),
		vm.Instr(vm.OpExplode),
		vm.Instr(vm.OpDissolve),
		vm.Instr(vm.OpJoin),
	)
}

// --- WRITE ---

func TestRunWrite1(t *testing.T) {
	m := vm.NewVM(1024, 512, 512)
	m.AddDefaultGlobals()

	m.Run(vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Str("Hello, World!\n"),
		vm.Instr(vm.OpGetGlobal, 1), // *STD-OUT*
		vm.Instr(vm.OpWrite),
	))

	// EXPECTED: Hello, World!
}

func TestRunWrite2(t *testing.T) {
	m := vm.NewVM(1024, 512, 512)
	m.AddDefaultGlobals()

	m.Run(vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Str("\n"),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpGetGlobal, 1), // *STD-OUT*
		vm.Instr(vm.OpWrite),
	))

	// EXPECTED: [1 2 3]
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
