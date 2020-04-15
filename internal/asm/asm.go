package asm

import "github.com/mhoertnagl/splis2/internal/vm"

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
			vm.Instr(x.Op, a.lbls[x.Name])
			ip += 1 + 8
		case *AsmIns:
			mt, err := vm.LookupMeta(x.Op)
			if err != nil {
				panic(err)
			}
			vm.Instr(x.Op, x.Args...)
			ip += 1 + uint64(mt.Size())
		case *AsmStr:
			mt, err := vm.LookupMeta(vm.OpStr)
			if err != nil {
				panic(err)
			}
			vm.Str(x.Str)
			ip += 1 + uint64(mt.Size()) + uint64(len(x.Str))
		}
	}
}

func (a *Assembler) assemble(code AsmCode) []byte {
	bin := make([]byte, 0)
	for _, line := range code {
		switch x := line.(type) {
		case *AsmLabeled:
			bin = append(bin, vm.Instr(x.Op, a.lbls[x.Name])...)
		case *AsmIns:
			bin = append(bin, vm.Instr(x.Op, x.Args...)...)
		case *AsmStr:
			bin = append(bin, vm.Str(x.Str)...)
		}
	}
	return bin
}

// // Instr creates a new instruction from an opcode and a variable number of
// // arguments.
// func Instr(op vm.Op, args ...uint64) []byte {
// 	if m, err := vm.LookupMeta(op); err == nil {
// 		sz := 1 + m.Size()
// 		ins := make([]byte, sz)
// 		pos := 1
//
// 		ins[0] = op
// 		for i, as := range m.Args {
// 			switch as {
// 			case 1:
// 				ins[pos] = uint8(args[i])
// 			case 2:
// 				binary.BigEndian.PutUint16(ins[pos:pos+2], uint16(args[i]))
// 			case 4:
// 				binary.BigEndian.PutUint32(ins[pos:pos+4], uint32(args[i]))
// 			case 8:
// 				binary.BigEndian.PutUint64(ins[pos:pos+8], args[i])
// 			}
// 			pos += as
// 		}
// 		return ins
// 	}
// 	panic(fmt.Sprintf("could not find meta infor for [%d]", op))
// }
//
// func Str(s string) []byte {
// 	b := []byte(s)
// 	ln := len(b)
// 	sz := 9 + ln
// 	ins := make([]byte, sz)
// 	ins[0] = vm.OpStr
// 	binary.BigEndian.PutUint64(ins[1:9], uint64(ln))
// 	copy(ins[9:sz], b)
// 	return ins
// }
//
// // Concat joins an array of instructions.
// func Concat(is [][]byte) []byte {
// 	return bytes.Join(is, []byte{})
// }
//
// // ConcatVar joins a variable number of instructions.
// func ConcatVar(is ...[]byte) []byte {
// 	return Concat(is)
// }
