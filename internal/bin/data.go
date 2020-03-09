package bin

import "github.com/mhoertnagl/splis2/internal/vm"

type Lib struct {
	Macros  []string
	Entries []uint64
	Fns     vm.Ins
	Code    vm.Ins
}

func NewLib() *Lib {
	return &Lib{
		Macros:  []string{},
		Entries: []uint64{},
		Fns:     vm.Ins{},
		Code:    vm.Ins{},
	}
}
