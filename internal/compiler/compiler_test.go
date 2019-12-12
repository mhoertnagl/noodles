package compiler_test

import (
	"bytes"
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
		t.Errorf("Mismatch at position [%d] Expecting %v but got %v.", x, ee, s)
	}
}
