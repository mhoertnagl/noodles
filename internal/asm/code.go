package asm

import (
	"github.com/mhoertnagl/splis2/internal/vm"
)

type AsmCmd interface {
}

type AsmLabel struct {
	Name string
}

type AsmLabeled struct {
	Op   vm.Op
	Name string
}

type AsmIns struct {
	Op   vm.Op
	Args []uint64
}

type AsmStr struct {
	Str string
}

type AsmCode []AsmCmd

func Label(name string) *AsmLabel {
	return &AsmLabel{Name: name}
}

func Labeled(op vm.Op, name string) *AsmLabeled {
	return &AsmLabeled{Op: op, Name: name}
}

func Instr(op vm.Op, args ...uint64) *AsmIns {
	return &AsmIns{Op: op, Args: args}
}

func Str(str string) *AsmStr {
	return &AsmStr{Str: str}
}

// func AsmBool(n bool) *AsmIns {
// 	if n {
// 		return &AsmIns{Op: vm.OpTrue}
// 	}
// 	return &AsmIns{Op: vm.OpTrue}
// }
