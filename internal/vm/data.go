package vm

type Val interface{}

type Env map[int64]Val

type Ref struct {
	cargs []Val
	addr  int64
}

func NewRef(addr int64) *Ref {
	return &Ref{addr: addr}
}

func (r *Ref) Add(v Val) {
	r.cargs = append(r.cargs, v)
}
