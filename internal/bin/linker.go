package bin

import (
	"encoding/binary"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type Linker struct {
	lib *Lib
}

func NewLinker() *Linker {
	return &Linker{
		lib: NewLib(),
	}
}

func (lnk *Linker) Add(unit *Lib) {
	etrLen := len(lnk.lib.entries)
	fnsLen := len(lnk.lib.fns)
	// Join the macro arrays.
	lnk.lib.macros = append(lnk.lib.macros, unit.macros...)
	// Add the new local functions to the library.
	lnk.lib.fns = append(lnk.lib.fns, unit.fns...)
	// Append and update the function entry points of the unit. Add the
	// length of the library's function block before the merge to each
	// entry point in the added unit.
	for unitEntry := range unit.entries {
		lnk.lib.entries = append(lnk.lib.entries, uint64(unitEntry+fnsLen))
	}
	// Update the entry ids of the local function calls.
	updateRefIndexes(unit.code, uint64(etrLen))
	// Add the unit's updated code segment to the library.
	lnk.lib.code = append(lnk.lib.code, unit.code...)
}

func (lnk *Linker) Lib() *Lib {
	return lnk.lib
}

func (lnk *Linker) Code() []byte {
	updateRefAddresses(lnk.lib)
	return lnk.lib.code
}

func updateRefIndexes(code []byte, offset uint64) {
	for ip := 0; ip < len(code); {
		op := code[ip]
		mt, err := vm.LookupMeta(op)
		if err != nil {
			panic(err)
		}
		switch op {
		case vm.OpRef:
			// Shift the index of the Ref cell by offset bytes.
			id := binary.BigEndian.Uint64(code[ip+1 : ip+9])
			vm.Correct(code, ip+1, id+offset)
		case vm.OpStr:
			// String commands are of variable length. The first argument is the
			// length of the string. We need to skip over the string as well and
			// thus add the string length to the position pointer i.
			strLen := binary.BigEndian.Uint64(code[ip+1 : ip+9])
			ip += int(strLen)
		}
		ip += mt.Size() + 1
	}
}

func updateRefAddresses(lib *Lib) {
	for ip := 0; ip < len(lib.code); {
		op := lib.code[ip]
		mt, err := vm.LookupMeta(op)
		if err != nil {
			panic(err)
		}
		switch op {
		case vm.OpRef:
			// The only argument of Ref is the index of the referenced function.
			// The library's entries array contains the real memory address for
			// this function. Replace the Ref argument with this entry address.
			id := binary.BigEndian.Uint64(lib.code[ip+1 : ip+9])
			vm.Correct(lib.code, ip+1, lib.entries[id])
		case vm.OpStr:
			// String commands are of variable length. The first argument is the
			// length of the string. We need to skip over the string as well and
			// thus add the string length to the position pointer i.
			strLen := binary.BigEndian.Uint64(lib.code[ip+1 : ip+9])
			ip += int(strLen)
		}
		ip += mt.Size() + 1
	}
}
