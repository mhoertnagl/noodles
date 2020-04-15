package cmp

type Ctx struct {
	FnName string
	IsVal  bool
	IsLast bool
}

func NewCtx() *Ctx {
	return &Ctx{}
}

func (c *Ctx) NewCtx() *Ctx {
	return &Ctx{
		FnName: c.FnName,
		IsVal:  c.IsVal,
		IsLast: c.IsLast,
	}
}

func (c *Ctx) NewDefCtx(name string) *Ctx {
	return &Ctx{
		FnName: name,
		IsVal:  c.IsVal,
		IsLast: c.IsLast,
	}
}
