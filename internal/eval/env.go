package eval

import (
	"github.com/mhoertnagl/splis2/internal/read"
)

type SpecialForm func(Env, []read.Node) read.Node

type Env interface {
	Set(name string, val read.Node)
	Lookup(name string) read.Node
  AddSpecialForm(name string, fun SpecialForm)
  FindSpecialForm(name string) (SpecialForm, bool)
}

type env struct {
	outer    Env
  defs     map[string]read.Node
  specials map[string]SpecialForm
}

func NewEnv(outer Env) Env {
	return &env{
    outer: outer, 
    defs: make(map[string]read.Node),
    specials: make(map[string]SpecialForm),
  }
}

func (e *env) Set(name string, val read.Node) {
  if _, ok := e.defs[name]; !ok {
    e.defs[name] = val
  }
}

func (e *env) Lookup(name string) read.Node {
  if v, ok := e.defs[name]; ok {
    return v
  } else if e.outer != nil {
    return e.outer.Lookup(name)
  }
  // TODO: Return error node?
  return nil
}

func (e *env) AddSpecialForm(name string, fun SpecialForm) {
  e.specials[name] = fun
}

func (e *env) FindSpecialForm(name string) (SpecialForm, bool) {
  fun, ok := e.specials[name]
  return fun, ok
}
