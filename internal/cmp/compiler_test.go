package cmp_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/mhoertnagl/noodles/internal/asm"
	"github.com/mhoertnagl/noodles/internal/cmp"
	"github.com/mhoertnagl/noodles/internal/vm"
)

func TestCompileBool(t *testing.T) {
	testc(t, "true",
		asm.Instr(vm.OpTrue),
	)
	testc(t, "false",
		asm.Instr(vm.OpFalse),
	)
}

func TestCompileInteger(t *testing.T) {
	testc(t, "0",
		asm.Instr(vm.OpConst, 0),
	)
	testc(t, "1",
		asm.Instr(vm.OpConst, 1),
	)
}

func TestCompileFloat(t *testing.T) {
	testc(t, "0.1",
		asm.Instr(vm.OpConstF, math.Float64bits(0.1)),
	)
	testc(t, "1.0",
		asm.Instr(vm.OpConstF, math.Float64bits(1)),
	)
}

// --- STR ---

func TestCompileString(t *testing.T) {
	testc(t, `"Hello, World!"`,
		asm.Str("Hello, World!"),
	)
}

func TestCompileStringJoin(t *testing.T) {
	testc(t, `(join "Hello, " "World!")`,
		asm.Instr(vm.OpEnd),
		asm.Str("World!"),
		asm.Str("Hello, "),
		asm.Instr(vm.OpJoin),
	)
}

func TestCompileStringExplode(t *testing.T) {
	testc(t, `(explode "xyz")`,
		asm.Str("xyz"),
		asm.Instr(vm.OpExplode),
	)
}

func TestCompileStringJoinExplode(t *testing.T) {
	testc(t, `(join @(explode "xyz"))`,
		asm.Instr(vm.OpEnd),
		asm.Str("xyz"),
		asm.Instr(vm.OpExplode),
		asm.Instr(vm.OpDissolve),
		asm.Instr(vm.OpJoin),
	)
}

// --- ARITH ---

func TestCompileAdd(t *testing.T) {
	testc(t, "(+)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1 2)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1 (+ 2 3))",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
	)
	testc(t, "(+ (+ 1 2) 3)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1 2 3)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
	)
}

func TestCompileSub(t *testing.T) {
	testc(t, "(-)",
		asm.Instr(vm.OpConst, 0),
	)
	testc(t, "(- 1)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
	)
	testc(t, "(- 2 1)",
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
	)
	testc(t, "(- 3 (- 2 1))",
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpSub),
	)
	testc(t, "(- (- 3 2) 1)",
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
	)
}

func TestCompileMul(t *testing.T) {
	testc(t, "(*)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpMul),
	)
	testc(t, "(* 2)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpMul),
	)
	testc(t, "(* 1 2)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpMul),
	)
	testc(t, "(* 1 (* 2 3))",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpMul),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpMul),
	)
	testc(t, "(* (* 1 2) 3)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpMul),
		asm.Instr(vm.OpMul),
	)
	testc(t, "(* 1 2 3)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpMul),
	)
}

func TestCompileDiv(t *testing.T) {
	testc(t, "(/)",
		asm.Instr(vm.OpConst, 1),
	)
	testc(t, "(/ 2)",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpDiv),
	)
	testc(t, "(/ 2 1)",
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpDiv),
	)
	testc(t, "(/ 3 (/ 2 1))",
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpDiv),
	)
	testc(t, "(/ (/ 3 2) 1)",
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpDiv),
	)
}

func TestCompileMod(t *testing.T) {
	testc(t, "(mod 7 3)",
		asm.Instr(vm.OpConst, 7),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpMod),
	)
}

func TestCompileRandom(t *testing.T) {
	testc(t, "(random 100)",
		asm.Instr(vm.OpConst, 100),
		asm.Instr(vm.OpRand),
	)
}

func TestCompileNot(t *testing.T) {
	testc(t, "(not true)",
		asm.Instr(vm.OpTrue),
		asm.Instr(vm.OpNot),
	)
}

func TestCompileLT(t *testing.T) {
	testc(t, "(< 0 1)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpLT),
	)
}

func TestCompileLE(t *testing.T) {
	testc(t, "(<= 0 1)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpLE),
	)
}

func TestCompileGT(t *testing.T) {
	testc(t, "(> 0 1)",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpLT),
	)
}

func TestCompileGE(t *testing.T) {
	testc(t, "(>= 0 1)",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpLE),
	)
}

func TestCompileEQ(t *testing.T) {
	testc(t, "(= 0 1)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEQ),
	)
}

func TestCompileNE(t *testing.T) {
	testc(t, "(!= 0 1)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpNE),
	)
}

// --- SET ---

func TestCompileSet1(t *testing.T) {
	testc(t, `(do
      (set x 2)
      (+ x x)
    )`,
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
	)
}

func TestCompileSet2(t *testing.T) {
	testc(t, `(do
      (set x 2)
      (def f (fn [y] (+ x y)))
      (f 3)
    )`,
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	)
}

// --- LET ---

func TestCompileLet1(t *testing.T) {
	testc(t, "(let (a (+ 1 1)) (+ a a))",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpDropArgs, 1),
	)
}

func TestCompileLet2(t *testing.T) {
	testc(t, "(let (a 1 b 2) (+ a b))",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpDropArgs, 2),
	)
}

func TestCompileLet3(t *testing.T) {
	testc(t, "(let (a 1 b 2) (- b a))",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpDropArgs, 2),
	)
}

func TestCompileLet4(t *testing.T) {
	testc(t, `
    (let (a 2)
      (let (b 6)
        (/ b a)))`,
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpDropArgs, 1),
		asm.Instr(vm.OpDropArgs, 1),
	)
}

// TODO: GetArg 0 Gibs noch nicht!

func TestCompileLet5(t *testing.T) {
	testc(t, `
    (let
      (b (fn [m] (if (= m 0) 1 (b (- m 1))) ))
      (b 1) )`,
		asm.Labeled(vm.OpJump, "L0"),
		// BEGIN FN
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCall),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		// END FN
		asm.Label("L0"),
		// TODO: An dieser Stelle ist noch nichts in FRAME[0].
		//       Ref erwartet hier als argument sich selbst.
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L1"),
		// BEGIN LET
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpDropArgs, 1),
	)
}

func TestCompileDef1(t *testing.T) {
	testc(t, "(def b (+ 1 1))",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpSetGlobal, 0),
	)
}

// --- IF ---

func TestCompileIf1(t *testing.T) {
	testc(t, "(if true 1)",
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Label("L0"),
	)
}

func TestCompileIf2(t *testing.T) {
	testc(t, "(if false 1 0)",
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 0),
		asm.Label("L1"),
	)
}

func TestCompileIf3(t *testing.T) {
	testc(t, "(if (= 1 0) 42 21)",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	)
}

func TestCompileIf4(t *testing.T) {
	testc(t, "(if (= 0 0) 42 21)",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	)
}

// --- COND ---

func TestCompileCond0(t *testing.T) {
	testc(t, "(cond)")
}

func TestCompileCond1(t *testing.T) {
	testc(t, `
    (cond
      true 1)`,
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Label("L0"),
	)
}

func TestCompileCond2(t *testing.T) {
	testc(t, `
    (cond
      false 1
      true  2)`,
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L1"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 2),
		asm.Label("L0"),
	)
}

func TestCompileCond3(t *testing.T) {
	testc(t, `
    (cond
      false 1
      false 2
      true  3)`,
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L1"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 2),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L2"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 3),
		asm.Label("L0"),
	)
}

// --- AND ---

func TestCompileAnd0(t *testing.T) {
	testc(t, "(and)",
		asm.Instr(vm.OpTrue),
	)
}

func TestCompileAnd1(t *testing.T) {
	testc(t, "(and false)",
		asm.Instr(vm.OpFalse),
	)
}

func TestCompileAnd2(t *testing.T) {
	testc(t, "(and false true)",
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	)
}

func TestCompileAnd3(t *testing.T) {
	testc(t, "(and false true false)",
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	)
}

// --- OR ---

func TestCompileOr0(t *testing.T) {
	testc(t, "(or)",
		asm.Instr(vm.OpFalse),
	)
}

func TestCompileOr1(t *testing.T) {
	testc(t, "(or false)",
		asm.Instr(vm.OpFalse),
	)
}

func TestCompileOr2(t *testing.T) {
	testc(t, "(or false true)",
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	)
}

func TestCompileOr3(t *testing.T) {
	testc(t, "(or false true false)",
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	)
}

// func TestCompileBitAnd(t *testing.T) {
// 	testc(t, "(bit-and 12 10)",
// 		asm.Instr(vm.OpConst, 12),
// 		asm.Instr(vm.OpConst, 10),
// 		asm.Instr(vm.OpAnd),
//
// 	)
// }
//
// func TestCompileBitOr(t *testing.T) {
// 	testc(t, "(bit-or 12 10)",
// 		asm.Instr(vm.OpConst, 12),
// 		asm.Instr(vm.OpConst, 10),
// 		asm.Instr(vm.OpOr),
//
// 	)
// }
//
// func TestCompileBitShiftLeft(t *testing.T) {
// 	testc(t, "(bit-shift-left 8 2)",
// 		asm.Instr(vm.OpConst, 8),
// 		asm.Instr(vm.OpConst, 2),
// 		asm.Instr(vm.OpSll),
//
// 	)
// }
//
// func TestCompileBitShiftRight(t *testing.T) {
// 	testc(t, "(bit-shift-right 8 2)",
// 		asm.Instr(vm.OpConst, 8),
// 		asm.Instr(vm.OpConst, 2),
// 		asm.Instr(vm.OpSrl),
//
// 	)
// }
//
// func TestCompileBitShiftRightSigned(t *testing.T) {
// 	testc(t, "(bit-shift-right-signed (- 8) 2)",
// 		asm.Instr(vm.OpConst, 8),
// 		asm.Instr(vm.OpSub),
// 		asm.Instr(vm.OpConst, 2),
// 		asm.Instr(vm.OpSra),
//
// 	)
// }

func TestCompileDo(t *testing.T) {
	testc(t, "(do (def a 1) (def b 2) (+ a b))",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpSetGlobal, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetGlobal, 1),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpAdd),
	)
}

// --- VECTORS ---

func TestCompileVector1(t *testing.T) {
	testc(t, "[1 2 3]",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
	)
}

func TestCompileVector2(t *testing.T) {
	testc(t, "(.+ 1 (.+ 2 (.+ 3 [])))",
		asm.Instr(vm.OpEmptyVector),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpCons),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpCons),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpCons),
	)
}

func TestCompileVectorAppend(t *testing.T) {
	testc(t, "(+. [1 2 3] 4)",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpAppend),
	)
}

func TestCompileVectorConcat(t *testing.T) {
	testc(t, "(++ [1 2] [3 4] [5 6])",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpConcat),
	)
}

func TestCompileFirstVector(t *testing.T) {
	testc(t, "(nth 0 [1 2 3])",
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpNth),
	)
}

func TestCompileDropVector(t *testing.T) {
	testc(t, "(drop 1 [1 2 3])",
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpDrop),
	)
}

func TestCompileLenVector(t *testing.T) {
	testc(t, "(len [1 2 3])",
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpLength),
	)
}

// --- FN ---

func TestCompileAnonymousFun0(t *testing.T) {
	testc(t, `(fn [] 1)`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
	)
}

func TestCompileAnonymousFun1(t *testing.T) {
	testc(t, `(fn [x] (+ x 1))`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
	)
}

func TestCompileAnonymousFun11(t *testing.T) {
	testc(t, `(fn [] (fn [x] (+ x 1)))`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Ref(0, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
	)
}

func TestCompileAnonymousFun2(t *testing.T) {
	testc(t, `((fn [x] (+ x 1)) 1)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileAnonymousFun3(t *testing.T) {
	testc(t, `(+ ((fn [x] (+ x 1)) 1) 1)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpAdd),
	)
}

func TestCompileAnonymousFun4(t *testing.T) {
	testc(t, `(((fn [] (fn [x] (+ x 1)))) 1)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Ref(0, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileAnonymousNestedFun(t *testing.T) {
	testc(t, `((fn [n] ((fn [m] (/ n m)) 3)) 6)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Labeled(vm.OpJump, "L2"),
		// (fn [m] (/ n m))
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpReturn),
		// end
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileLeafFunDef(t *testing.T) {
	testc(t, `
    (do
      (def inc (fn [x] (+ x 1)))
      (+ (inc 1) 1)
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpAdd),
	)
}

func TestCompileVariadicFun1(t *testing.T) {
	testc(t, `((fn [& xs] xs) 1 2 3 4)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 0),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileVariadicFun2(t *testing.T) {
	testc(t, `((fn [x & xs] (.+ x xs)) 1 2 3 4)`,
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCons),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileFac(t *testing.T) {
	testc(t, `
    (do
      (def fac
        (fn [n]
          (do
            (if (= n 0)
              1
              (* n (fac (- n 1)))
            )
          )
        )
      )
      (fac 5)
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpMul),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileTailFac(t *testing.T) {
	testc(t, `
    (do
      (def _fac
        (fn [n acc]
          (do
            (if (= n 0)
              acc
              (rec (_fac (- n 1) (* n acc)))
            )
          )
				)
			)
			(def fac (fn [n] (_fac n 1) ))
      (fac 5)
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpGetArg, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpMul),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpRecCall),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Labeled(vm.OpJump, "L4"),
		asm.Label("L5"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L4"),
		asm.Ref(0, "L5"),
		asm.Instr(vm.OpSetGlobal, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpGetGlobal, 1),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileClosure1(t *testing.T) {
	testc(t, `
    (do
      (def ggg (fn [n] (fn [] n)) )
      ((ggg 6))
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpCall),
	)
}

func TestCompileClosure2(t *testing.T) {
	testc(t, `
    (do
      (def m2
        (fn [n]
          (do (set min 1)
              ((fn [x] (+ x min)) n) )))
      (m2 6)
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 1),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	)
}

//
// func TestCompileClosure3(t *testing.T) {
// 	fmt.Println("\n\n=== Closure 3")
// 	testc(t, `
//     (do
//       (def setclosure1 (fn [a b c]
//         (do (set lst [a b c])
//             (set min (nth 0 lst))
//             (set f (fn [x] (= x min) ))
//             (f (nth 0 lst)))))
//       (setclosure1 1 2 3)
//     )`,
// 		asm.Labeled(vm.OpJump, "L0"),
// 		asm.Label("L1"),
// 		asm.Instr(vm.OpPushArgs, 3),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpEnd),
// 		asm.Instr(vm.OpGetArg, 2),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpList),
// 		asm.Instr(vm.OpPushArgs, 1),
// 		asm.Instr(vm.OpConst, 0),
// 		asm.Instr(vm.OpGetArg, 3),
// 		asm.Instr(vm.OpNth),
// 		asm.Instr(vm.OpPushArgs, 1),
// 		asm.Labeled(vm.OpJump, "L2"),
// 		asm.Label("L3"),
// 		asm.Instr(vm.OpPushArgs, 2),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpEQ),
// 		asm.Instr(vm.OpReturn),
// 		asm.Label("L2"),
// 		asm.Instr(vm.OpGetArg, 4),
// 		asm.Ref(1, "L3"),
// 		asm.Instr(vm.OpPushArgs, 1),
//
// 		asm.Instr(vm.OpConst, 1),
// 		asm.Instr(vm.OpPushArgs, 1),
// 		asm.Instr(vm.OpEnd),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Labeled(vm.OpJump, "L2"),
// 		asm.Label("L3"),
// 		asm.Instr(vm.OpPushArgs, 2),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpEnd),
// 		asm.Instr(vm.OpGetArg, 0),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Instr(vm.OpAdd),
// 		asm.Instr(vm.OpReturn),
// 		asm.Label("L2"),
// 		asm.Instr(vm.OpGetArg, 1),
// 		asm.Ref(1, "L3"),
// 		asm.Instr(vm.OpCall),
// 		asm.Instr(vm.OpReturn),
// 		asm.Label("L0"),
// 		asm.Ref(0, "L1"),
// 		asm.Instr(vm.OpSetGlobal, 0),
// 		asm.Instr(vm.OpEnd),
// 		asm.Instr(vm.OpConst, 6),
// 		asm.Instr(vm.OpGetGlobal, 0),
// 		asm.Instr(vm.OpCall),
// 	)
// }

func TestCompileClosure(t *testing.T) {
	fmt.Println("--- Closure")
	testc(t, `
    (do
      (def divN (fn [n]
        (fn [x] (/ x n)) ))
      ((divN 3) 9)
    )`,
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		// def divN
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		// ((divN 3) 9 )
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 9),
		// (divN 3)
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		// --- (divN 3)
		asm.Instr(vm.OpCall),
		// --- ((divN 3) 9 )
	)
}

// func TestCompilePrimitiveAsArgument(t *testing.T) {
// 	testc(t, `((fn [op a b] (op a b)) +)`,
// 		asm.Labeled(vm.OpJump, "L0"),
// 		asm.Label("L1"),
// 		asm.Instr(vm.OpPop),
// 		asm.Instr(vm.OpConst, 1),
// 		asm.Instr(vm.OpReturn),
// 		asm.Label("L0"),
// 		asm.Labeled(vm.OpRef, "L1"),
// 	)
// }

// --- WRITE ---

func TestCompileWrite1(t *testing.T) {
	testcd(t, `(write *STD-OUT* "Hello, World!\n")`,
		asm.Instr(vm.OpEnd),
		asm.Str("Hello, World!\n"),
		asm.Instr(vm.OpGetGlobal, 1),
		asm.Instr(vm.OpWrite),
	)
}

// --- HALT ---

func TestCompileHalt(t *testing.T) {
	testc(t, "(halt)",
		asm.Instr(vm.OpHalt),
	)
}

func testc(t *testing.T, i string, e ...asm.AsmCmd) {
	t.Helper()
	r := cmp.NewReader()
	p := cmp.NewParser()
	c := cmp.NewCompiler()

	r.Load(i)
	n := p.Parse(r)
	s := c.Compile(n)

	compareAssembly(t, s, e)
}

func testcd(t *testing.T, i string, e ...asm.AsmCmd) {
	t.Helper()
	r := cmp.NewReader()
	p := cmp.NewParser()
	c := cmp.NewCompiler()

	c.AddDefaultGlobals()

	r.Load(i)
	n := p.Parse(r)
	s := c.Compile(n)

	compareAssembly(t, s, e)
}

func compareAssembly(t *testing.T, a []asm.AsmCmd, e []asm.AsmCmd) {
	t.Helper()

	err := false

	d := asm.NewAsmPrinter()
	da := d.Print(a)
	de := d.Print(e)
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

		buf.WriteString(fmt.Sprintf("%3d: %-20s%-20s", i, sa, se))

		if sa != se {
			err = true
			buf.WriteString("<--")
		}

		buf.WriteString("\n")
	}

	if err {
		t.Errorf("\n%s\n", buf.String())
	}
}
