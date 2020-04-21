package vm

import "os"

// AddGlobal assigns a value to an ID in the global definitions.
// NOTE: Every definition has to be registerd in the compiler as well.
func (m *VM) AddGlobal(id uint64, val Val) {
	m.defs[id] = val
}

func (m *VM) AddDefaultGlobals() {

	m.AddGlobal(0, os.Stdin)  // *STD-IN*
	m.AddGlobal(1, os.Stdout) // *STD-OUT*
	m.AddGlobal(2, os.Stderr) // *STD-ERR*
}
