package eval

import (
	"github.com/mhoertnagl/splis2/internal/read"
)

type Env interface {
	Set(name string, val read.Node)
	Lookup(name string) read.Node
}

type env struct {
	outer Env
  defs map[string]read.Node
}

func NewEnv(outer Env) Env {
	return &env{outer: outer, defs: make(map[string]read.Node)}
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
