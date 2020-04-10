package cmp_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mhoertnagl/splis2/internal/cmp"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestCompileBool(t *testing.T) {
	testc(t, "true",
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "false",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileInteger(t *testing.T) {
	testc(t, "0",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "1",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileString(t *testing.T) {
	testc(t, `"Hello, World!"`,
		vm.Str("Hello, World!"),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileAdd(t *testing.T) {
	testc(t, "(+)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(+ 1)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(+ 1 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(+ 1 (+ 2 3))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(+ (+ 1 2) 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(+ 1 2 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileSub(t *testing.T) {
	testc(t, "(-)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(- 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(- 2 1)",
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(- 3 (- 2 1))",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(- (- 3 2) 1)",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileMul(t *testing.T) {
	testc(t, "(*)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(* 2)",
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(* 1 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(* 1 (* 2 3))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(* (* 1 2) 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(* 1 2 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileDiv(t *testing.T) {
	testc(t, "(/)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(/ 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(/ 2 1)",
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(/ 3 (/ 2 1))",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpHalt),
	)
	testc(t, "(/ (/ 3 2) 1)",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileNot(t *testing.T) {
	testc(t, "(not true)",
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpNot),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileLT(t *testing.T) {
	testc(t, "(< 0 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpLT),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileLE(t *testing.T) {
	testc(t, "(<= 0 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpLE),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileGT(t *testing.T) {
	testc(t, "(> 0 1)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpLT),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileGE(t *testing.T) {
	testc(t, "(>= 0 1)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpLE),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileEQ(t *testing.T) {
	testc(t, "(= 0 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileNE(t *testing.T) {
	testc(t, "(!= 0 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpNE),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileLet1(t *testing.T) {
	testc(t, "(let (a (+ 1 1)) (+ a a))",
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
}

func TestCompileLet2(t *testing.T) {
	testc(t, "(let (a 1 b 2) (+ a b))",
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
}

func TestCompileLet3(t *testing.T) {
	testc(t, "(let (a 1 b 2) (- b a))",
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
}

func TestCompileLet4(t *testing.T) {
	testc(t, `
    (let (a 2)
      (let (b 6)
        (/ b a)))`,
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
}

func TestCompileDef1(t *testing.T) {
	testc(t, "(def b (+ 1 1))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileIf1(t *testing.T) {
	testc(t, "(if true 1)",
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 9),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileIf2(t *testing.T) {
	testc(t, "(if false 1 0)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileIf3(t *testing.T) {
	testc(t, "(if (= 1 0) 42 21)",
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

func TestCompileIf4(t *testing.T) {
	testc(t, "(if (= 0 0) 42 21)",
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

func TestCompileAnd0(t *testing.T) {
	testc(t, "(and)",
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileAnd1(t *testing.T) {
	testc(t, "(and false)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileAnd2(t *testing.T) {
	testc(t, "(and false true)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 10),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileAnd3(t *testing.T) {
	testc(t, "(and false true false)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 10),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileOr0(t *testing.T) {
	testc(t, "(or)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileOr1(t *testing.T) {
	testc(t, "(or false)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileOr2(t *testing.T) {
	testc(t, "(or false true)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 10),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileOr3(t *testing.T) {
	testc(t, "(or false true false)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 10),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 1),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpHalt),
	)
}

// func TestCompileBitAnd(t *testing.T) {
// 	testc(t, "(bit-and 12 10)",
// 		vm.Instr(vm.OpConst, 12),
// 		vm.Instr(vm.OpConst, 10),
// 		vm.Instr(vm.OpAnd),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestCompileBitOr(t *testing.T) {
// 	testc(t, "(bit-or 12 10)",
// 		vm.Instr(vm.OpConst, 12),
// 		vm.Instr(vm.OpConst, 10),
// 		vm.Instr(vm.OpOr),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestCompileBitShiftLeft(t *testing.T) {
// 	testc(t, "(bit-shift-left 8 2)",
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSll),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestCompileBitShiftRight(t *testing.T) {
// 	testc(t, "(bit-shift-right 8 2)",
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSrl),
// 		vm.Instr(vm.OpHalt),
// 	)
// }
//
// func TestCompileBitShiftRightSigned(t *testing.T) {
// 	testc(t, "(bit-shift-right-signed (- 8) 2)",
// 		vm.Instr(vm.OpConst, 8),
// 		vm.Instr(vm.OpSub),
// 		vm.Instr(vm.OpConst, 2),
// 		vm.Instr(vm.OpSra),
// 		vm.Instr(vm.OpHalt),
// 	)
// }

func TestCompileDo(t *testing.T) {
	testc(t, "(do (def a 1) (def b 2) (+ a b))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetGlobal, 1),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpGetGlobal, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileVector1(t *testing.T) {
	testc(t, "[1 2 3]",
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileVector2(t *testing.T) {
	testc(t, "(:: 1 (:: 2 (:: 3 [])))",
		vm.Instr(vm.OpEmptyVector),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileFirstVector(t *testing.T) {
	testc(t, "(fst [1 2 3])",
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpFst),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileRestVector(t *testing.T) {
	testc(t, "(rest [1 2 3])",
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpRest),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileLenVector(t *testing.T) {
	testc(t, "(len [1 2 3])",
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpLength),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileAnonymousFun1(t *testing.T) {
	testc(t, `(fn [x] (+ x 1))`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun11(t *testing.T) {
	testc(t, `(fn [] (fn [x] (+ x 1)))`,
		vm.Instr(vm.OpRef, 40),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		// (fn [] ...)
		// 0-adic functions don't require a local environment.
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun2(t *testing.T) {
	testc(t, `((fn [x] (+ x 1)) 1)`,
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
}

func TestCompileAnonymousFun3(t *testing.T) {
	testc(t, `(+ ((fn [x] (+ x 1)) 1) 1)`,
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
}

func TestCompileAnonymousFun4(t *testing.T) {
	testc(t, `(((fn [] (fn [x] (+ x 1)))) 1)`,
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
}

func TestCompileLeafFunDef(t *testing.T) {
	testc(t, `
    (do
      (def inc (fn [x] (+ x 1)))
      (+ (inc 1) 1)
    )`,
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
}

// TODO: Deeply nested function calls

func TestCompileVariadicFun(t *testing.T) {
	testc(t, `((fn [x & xs] (:: x xs)) 1 2 3 4)`,
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
}

func TestCompileSimpleQuote(t *testing.T) {
	testc(t, `'(+ 1 1)`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [] (+ 1 1))
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileReplacementQuote(t *testing.T) {
	testc(t, `'(+ ~a ~b)`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [a b] (+ a b))
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileSpliceQuote(t *testing.T) {
	testc(t, `'(+ ~a ~@b)`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [a b] (+ a @b))
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpDissolve),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileSpliceQuote2(t *testing.T) {
	testc(t, `'(+ ~@a ~@b)`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [a b] (+ @a @b))
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpDissolve),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpDissolve),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileQuote3(t *testing.T) {
	code := `
  (do
    (def cube '(* ~n ~n ~n))
    (cube 3)
  )
  `
	testc(t, code,
		vm.Instr(vm.OpRef, 39),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileFac(t *testing.T) {
	testc(t, `
    (do
      (def fac
        (fn [n]
          (do
            (debug 3)
            (if (= n 0)
              1
              (* n (fac (- n 1)))
            )
          )
        )
      )
      (debug 3)
      (fac 5)
      (debug 3)
    )`,
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

func TestCompileTailFac(t *testing.T) {
	testc(t, `
    (do
      (def _fac
        (fn [n acc]
          (do
            (debug 3)
            (if (= n 0)
              acc
              (_fac (- n 1) (* n acc))
            )
          )
				)
			)
			(def fac (fn [n] (_fac n 1) ))
      (debug 3)
      (fac 5)
      (debug 3)
    )`,
		vm.Instr(vm.OpRef, 75),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpRef, 190),
		vm.Instr(vm.OpSetGlobal, 1),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 1),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpHalt),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpDebug, 3),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpJump, 49),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
	)
}

func testc(t *testing.T, i string, es ...vm.Ins) {
	t.Helper()
	r := cmp.NewReader()
	p := cmp.NewParser()
	c := cmp.NewCompiler()
	w := cmp.NewQuoteRewriter()
	r.Load(i)
	n := p.Parse(r)
	n = w.Rewrite(n)
	s := c.Compile(n)
	e := vm.Concat(es)
	compareAssembly(t, s, e)
}

func compareAssembly(t *testing.T, a []byte, e []byte) {
	t.Helper()
	if bytes.Compare(a, e) != 0 {
		t.Errorf("Expecting \n  %v\n but got \n  %v", e, a)
		d := vm.NewDisassembler()
		da := d.Disassemble(a)
		de := d.Disassemble(e)
		la := len(da)
		le := len(de)
		lm := la
		if le > la {
			lm = le
		}
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("     %-20s%s\n", "Actual", "Expecting"))
		buf.WriteString(fmt.Sprintf("     %-20s%s\n", "------", "---------"))
		for i := 0; i < lm; i++ {
			sa := ""
			if i < la {
				sa = da[i]
			}
			se := ""
			if i < le {
				se = de[i]
			}
			buf.WriteString(fmt.Sprintf("%3d: %-20s%s\n", i, sa, se))
		}
		t.Errorf("\n%s", buf.String())
	}
}
