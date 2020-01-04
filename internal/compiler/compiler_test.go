package compiler_test

import (
	"bytes"
	"hash/fnv"
	"testing"

	"github.com/mhoertnagl/splis2/internal/compiler"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestCompileBool(t *testing.T) {
	testc(t, "true", vm.Instr(vm.OpTrue))
	testc(t, "false", vm.Instr(vm.OpFalse))
}

func TestCompileInteger(t *testing.T) {
	testc(t, "0", vm.Instr(vm.OpConst, 0))
	testc(t, "1", vm.Instr(vm.OpConst, 1))
}

func TestCompileAdd(t *testing.T) {
	testc(t, "(+)",
		vm.Instr(vm.OpConst, 0),
	)
	testc(t, "(+ 1)",
		vm.Instr(vm.OpConst, 1),
	)
	testc(t, "(+ 1 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1 (+ 2 3))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpAdd),
	)
	testc(t, "(+ (+ 1 2) 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
	)
	testc(t, "(+ 1 2 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpAdd),
	)
}

func TestCompileSub(t *testing.T) {
	testc(t, "(-)",
		vm.Instr(vm.OpConst, 0),
	)
	testc(t, "(- 1)",
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
	)
	testc(t, "(- 2 1)",
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
	)
	testc(t, "(- 3 (- 2 1))",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpSub),
	)
	testc(t, "(- (- 3 2) 1)",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
	)
}

func TestCompileMul(t *testing.T) {
	testc(t, "(*)",
		vm.Instr(vm.OpConst, 1),
	)
	testc(t, "(* 2)",
		vm.Instr(vm.OpConst, 2),
	)
	testc(t, "(* 1 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
	)
	testc(t, "(* 1 (* 2 3))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpMul),
	)
	testc(t, "(* (* 1 2) 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
	)
	testc(t, "(* 1 2 3)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpMul),
	)
}

func TestCompileDiv(t *testing.T) {
	testc(t, "(/)",
		vm.Instr(vm.OpConst, 1),
	)
	testc(t, "(/ 2)",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpDiv),
	)
	testc(t, "(/ 2 1)",
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
	)
	testc(t, "(/ 3 (/ 2 1))",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpDiv),
	)
	testc(t, "(/ (/ 3 2) 1)",
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpDiv),
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
	)
}

func TestCompileDef1(t *testing.T) {
	testc(t, "(def b (+ 1 1))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpSetGlobal, hash("b")),
	)
}

func TestCompileIf1(t *testing.T) {
	testc(t, "(if true 1 0)",
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestCompileIf2(t *testing.T) {
	testc(t, "(if false 1 0)",
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 18),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 9),
		vm.Instr(vm.OpConst, 0),
	)
}

func TestCompileDo(t *testing.T) {
	testc(t, "(do  (def a 1) (def b 2) (+ a b))",
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSetGlobal, hash("a")),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpSetGlobal, hash("b")),
		vm.Instr(vm.OpGetLocal, hash("a")),
		vm.Instr(vm.OpGetLocal, hash("b")),
		vm.Instr(vm.OpAdd),
	)
}

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
