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
	etrLen := len(lnk.lib.Entries)
	fnsLen := len(lnk.lib.Fns)
	// Join the macro arrays.
	lnk.lib.Macros = append(lnk.lib.Macros, unit.Macros...)
	// Add the new local functions to the library.
	lnk.lib.Fns = append(lnk.lib.Fns, unit.Fns...)
	// Append and update the function entry points of the unit. Add the
	// length of the library's function block before the merge to each
	// entry point in the added unit.
	for unitEntry := range unit.Entries {
		lnk.lib.Entries = append(lnk.lib.Entries, uint64(unitEntry+fnsLen))
	}
	// Update the entry ids of the local function calls.
	updateRefIndexes(unit.Code, uint64(etrLen))
	// Add the unit's updated code segment to the library.
	lnk.lib.Code = append(lnk.lib.Code, unit.Code...)
}

func (lnk *Linker) Lib() *Lib {
	return lnk.lib
}

func (lnk *Linker) Code() []byte {
	updateRefAddresses(lnk.lib)
	// TODO: Hier fehlen die Funktionen.
	// TODO: Append the functions to the end. Add OpHalt
	//       Or add to the front and insert jump to code.
	return lnk.lib.Code
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
	for ip := 0; ip < len(lib.Code); {
		op := lib.Code[ip]
		mt, err := vm.LookupMeta(op)
		if err != nil {
			panic(err)
		}
		switch op {
		case vm.OpRef:
			// The only argument of Ref is the index of the referenced function.
			// The library's entries array contains the real memory address for
			// this function. Replace the Ref argument with this entry address.
			id := binary.BigEndian.Uint64(lib.Code[ip+1 : ip+9])
			vm.Correct(lib.Code, ip+1, lib.Entries[id])
		case vm.OpStr:
			// String commands are of variable length. The first argument is the
			// length of the string. We need to skip over the string as well and
			// thus add the string length to the position pointer i.
			strLen := binary.BigEndian.Uint64(lib.Code[ip+1 : ip+9])
			ip += int(strLen)
		}
		ip += mt.Size() + 1
	}
}
