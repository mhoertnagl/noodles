package data

type Env interface {
	Set(name string, val Node) Node
	Lookup(name string) (Node, bool)
}

type env struct {
	outer Env
	defs  map[string]Node
}

func NewEnv(outer Env) Env {
	return &env{
		outer: outer,
		defs:  make(map[string]Node),
	}
}

func (e *env) Set(name string, val Node) Node {
	e.defs[name] = val
	return val
}

// TODO: Iterative version?
func (e *env) Lookup(name string) (Node, bool) {
	if v, ok := e.defs[name]; ok {
		return v, true
	} else if e.outer != nil {
		return e.outer.Lookup(name)
	}
	return nil, false
}
