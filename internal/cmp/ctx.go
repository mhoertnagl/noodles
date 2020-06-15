package cmp

type Ctx struct {
	Recurse bool
	// IsVal   bool
	// Last    bool
}

func NewCtx() *Ctx {
	return &Ctx{}
}

func (c *Ctx) NewCtx() *Ctx {
	return &Ctx{
		Recurse: c.Recurse,
		// IsVal:   c.IsVal,
		// Last:    c.Last,
	}
}

func (c *Ctx) NewRecCtx(recurse bool) *Ctx {
	return &Ctx{
		Recurse: recurse,
		// IsVal:   c.IsVal,
		// Last:    c.Last,
	}
}
