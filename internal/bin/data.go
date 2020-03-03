package bin

import "github.com/mhoertnagl/splis2/internal/vm"

type Lib struct {
	macros  []string
	entries []uint64
	fns     vm.Ins
	code    vm.Ins
}

func NewLib() *Lib {
	return &Lib{
		macros:  []string{},
		entries: []uint64{},
		fns:     vm.Ins{},
		code:    vm.Ins{},
	}
}
