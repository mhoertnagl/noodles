package cmp

// AddGlobal registers a name with the global definitions.
// NOTE: Every definition has to be registerd in the VM as well.
func (c *Compiler) AddGlobal(name string) uint64 {
	return c.defs.add(name)
}

func (c *Compiler) AddDefaultGlobals() {
	c.AddGlobal("*STD-IN*")
	c.AddGlobal("*STD-OUT*")
	c.AddGlobal("*STD-ERR*")
}
