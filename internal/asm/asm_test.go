package asm_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mhoertnagl/splis2/internal/asm"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestAssembleBool1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleBool2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleInteger(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 42),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 42),
	)
	testa(t, i, e)
}

// --- IF ---

func TestAssembleIf1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Label("L0"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 19),
		vm.Instr(vm.OpConst, 1),
	)
	testa(t, i, e)
}

func TestAssembleIf2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 0),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 37),
		vm.Instr(vm.OpConst, 0),
	)
	testa(t, i, e)
}

func TestAssembleIf3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
	testa(t, i, e)
}

func TestAssembleIf4(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
	testa(t, i, e)
}

// --- AND ---

func TestAssembleAnd0(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleAnd1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleAnd2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleAnd3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

// --- OR ---

func TestAssembleOr0(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleOr1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleOr2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleOr3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

// --- FN ---

func TestAssembleAnonymousFun1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		// (fn [x] (+ x 1))
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Labeled(vm.OpRef, "L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 39),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),

		vm.Instr(vm.OpRef, 9),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun11(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		// (fn [] ...)
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),

		// (fn [x] (+ x 1))
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),

		asm.Label("L2"),
		asm.Labeled(vm.OpRef, "L3"),
		asm.Instr(vm.OpReturn),

		asm.Label("L0"),
		asm.Labeled(vm.OpRef, "L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 59),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 49),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		// (fn [] ...)
		// 0-adic functions don't require a local environment.
		vm.Instr(vm.OpRef, 19),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 9),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun2(t *testing.T) {
	i := []asm.AsmCmd{
		// ((fn ...) 1)
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),

		// (fn [x] (+ x 1))
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),

		asm.Label("L0"),
		asm.Labeled(vm.OpRef, "L1"),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		// ((fn ...) 1)
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 49),
		// (fn [x] (+ x 1))
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 19),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

// --- STRING ---

func TestAssembleString(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Str("Hello, World!"),
	}
	e := vm.ConcatVar(
		vm.Str("Hello, World!"),
	)
	testa(t, i, e)
}

func testa(t *testing.T, i asm.AsmCode, e []byte) {
	t.Helper()
	a := asm.NewAssembler()
	ib := a.Assemble(i)
	compareAssembly(t, ib, e)
}

func compareAssembly(t *testing.T, a []byte, e []byte) {
	t.Helper()

	err := false

	d := asm.NewDisassembler()
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