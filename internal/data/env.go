package data

// type SpecialForm func(Env, string, []Node) Node

type Env interface {
	Set(name string, val Node) Node
	Lookup(name string) (Node, bool)
	// AddSpecialForm(name string, fun SpecialForm)
	// FindSpecialForm(name string) (SpecialForm, bool)
}

type env struct {
	outer Env
	defs  map[string]Node
	// specials map[string]SpecialForm
}

func NewEnv(outer Env) Env {
	return &env{
		outer: outer,
		defs:  make(map[string]Node),
		// specials: make(map[string]SpecialForm),
	}
}

func (e *env) Set(name string, val Node) Node {
	// if _, ok := e.defs[name]; !ok {
	// 	e.defs[name] = val
	// }
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

// func (e *env) AddSpecialForm(name string, fun SpecialForm) {
// 	e.specials[name] = fun
// }
//
// // TODO: Iterative version?
// func (e *env) FindSpecialForm(name string) (SpecialForm, bool) {
// 	if fun, ok := e.specials[name]; ok {
// 		return fun, ok
// 	} else if e.outer != nil {
// 		return e.outer.FindSpecialForm(name)
// 	}
// 	return nil, false
// }
