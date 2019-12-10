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

func testc(t *testing.T, i string, e ...vm.Ins) {
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
