package asm

import (
	"github.com/mhoertnagl/splis2/internal/vm"
)

type Assembler struct {
	lbls map[string]uint64
}

func NewAssembler() *Assembler {
	return &Assembler{
		lbls: make(map[string]uint64),
	}
}

func (a *Assembler) Assemble(code AsmCode) []byte {
	a.locateLabelPositions(code)
	return a.assemble(code)
}

func (a *Assembler) locateLabelPositions(code AsmCode) {
	ip := uint64(0)
	for _, line := range code {
		switch x := line.(type) {
		case *AsmLabel:
			a.lbls[x.Name] = ip
		case *AsmLabeled:
			ip += a.insInc(x.Op)
		case *AsmRef:
			ip += a.insInc(vm.OpRef)
		case *AsmIns:
			ip += a.insInc(x.Op)
		case *AsmStr:
			ip += a.insInc(vm.OpStr) + uint64(len(x.Str))
		}
	}
}

func (a *Assembler) insInc(op vm.Op) uint64 {
	mt, err := vm.LookupMeta(op)
	if err != nil {
		panic(err)
	}
	return 1 + uint64(mt.Size())
}

func (a *Assembler) assemble(code AsmCode) []byte {
	bin := make([]byte, 0)
	for _, line := range code {
		switch x := line.(type) {
		case *AsmLabeled:
			bin = append(bin, vm.Instr(x.Op, a.lbls[x.Name])...)
		case *AsmRef:
			bin = append(bin, vm.Instr(vm.OpRef, uint64(x.Cargs), a.lbls[x.Name])...)
		case *AsmIns:
			bin = append(bin, vm.Instr(x.Op, x.Args...)...)
		case *AsmStr:
			bin = append(bin, vm.Str(x.Str)...)
		}
	}
	return bin
}
