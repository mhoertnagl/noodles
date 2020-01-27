package compiler_test

import (
	"bytes"
	"hash/fnv"
	"testing"

	"github.com/mhoertnagl/splis2/internal/compiler"
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
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetLocal, hash("a")),
		vm.Instr(vm.OpGetLocal, hash("a")),
		vm.Instr(vm.OpGetLocal, hash("a")),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileLet2(t *testing.T) {
	testc(t, "(let (a 1 b 2) (+ a b))",
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetLocal, hash("a")),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetLocal, hash("b")),
		vm.Instr(vm.OpGetLocal, hash("a")),
		vm.Instr(vm.OpGetLocal, hash("b")),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileDef1(t *testing.T) {
	testc(t, "(def b (+ 1 1))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetGlobal, hash("b")),
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

func TestCompileDo(t *testing.T) {
	testc(t, "(do (def a 1) (def b 2) (+ a b))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetGlobal, hash("a")),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetGlobal, hash("b")),
		vm.Instr(vm.OpGetLocal, hash("a")),
		vm.Instr(vm.OpGetLocal, hash("b")),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
	)
}

func TestCompileVector1(t *testing.T) {
	testc(t, "[1 2 3]",
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

func TestCompileAnonymousFun1(t *testing.T) {
	testc(t, `(fn [x] (+ x 1))`,
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun11(t *testing.T) {
	testc(t, `(fn [] (fn [x] (+ x 1)))`,
		vm.Instr(vm.OpRef, 41),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
		// (fn [] ...)
		// 0-adic functions don't require a local environment.
		vm.Instr(vm.OpRef, 10),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun2(t *testing.T) {
	testc(t, `((fn [x] (+ x 1)) 1)`,
		// ((fn ...) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 20),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun3(t *testing.T) {
	testc(t, `(+ ((fn [x] (+ x 1)) 1) 1)`,
		// ((fn ...) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 30),
		vm.Instr(vm.OpCall),
		// (+ ((fn ...) 1) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
}

func TestCompileAnonymousFun4(t *testing.T) {
	testc(t, `(((fn [] (fn [x] (+ x 1)))) 1)`,
		// (((fn ...)) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 52),
		// Call the 0-adic function that returns the 1-adic function.
		vm.Instr(vm.OpCall),
		// Call the 1-adic function.
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
		// (fn [] ...)
		vm.Instr(vm.OpRef, 21),
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
		vm.Instr(vm.OpRef, 48),
		vm.Instr(vm.OpSetGlobal, hash("inc")),
		// (inc 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetGlobal, hash("inc")),
		vm.Instr(vm.OpCall),
		// (+ (inc ...) 1)
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpHalt),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
}

// TODO: Deeply nested function call
// TODO: A function that returns a function

func testc(t *testing.T, i string, e ...vm.Ins) {
	t.Helper()
	r := compiler.NewReader()
	p := compiler.NewParser()
	c := compiler.NewCompiler()
	r.Load(i)
	n := p.Parse(r)
	s := c.Compile(n)
	ee := vm.Concat(e)
	x := bytes.Compare(s, ee)
	if x != 0 {
		t.Errorf("Mismatch [%d] Expecting \n  %v\n but got \n  %v.", x, ee, s)
	}
}

func hash(sym string) uint64 {
	hg := fnv.New64()
	hg.Reset()
	hg.Write([]byte(sym))
	return hg.Sum64()
}
