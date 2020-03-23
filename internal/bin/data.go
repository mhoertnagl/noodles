package bin

import "github.com/mhoertnagl/splis2/internal/vm"

type Lib struct {
	Usings  []string
	Macros  []string
	Entries []uint64
	Fns     vm.Ins
	Code    vm.Ins
}

func NewLib() *Lib {
	return &Lib{
		Usings:  []string{},
		Macros:  []string{},
		Entries: []uint64{},
		Fns:     vm.Ins{},
		Code:    vm.Ins{},
	}
}
