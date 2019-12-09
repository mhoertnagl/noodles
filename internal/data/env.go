package data

import "os"

type Binding = map[string]Node

type Env interface {
	Set(name string, val Node) Node
	Lookup(name string) (Node, bool)
	Bindings() []Binding
}

type env struct {
	outer Env
	defs  Binding
}

func NewRootEnv() Env {
	e := NewEnv(nil)
	e.Set("*STDIN*", os.Stdin)
	e.Set("*STDOUT*", os.Stdout)
	e.Set("*STDERR*", os.Stderr)
	return e
}

func NewEnv(outer Env) Env {
	return &env{
		outer: outer,
		defs:  make(Binding),
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

func (e *env) Bindings() []Binding {
	bs := make([]Binding, 1)
	bs = append(bs, e.defs)
	if e.outer != nil {
		bs = append(bs, e.outer.Bindings()...)
	}
	return bs
}
