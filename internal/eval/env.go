package eval

import (
	"bytes"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
)

type SpecialForm func(Env, []read.Node) read.Node

type Env interface {
	Set(name string, val read.Node) read.Node
	Lookup(name string) read.Node
	AddSpecialForm(name string, fun SpecialForm)
	FindSpecialForm(name string) (SpecialForm, bool)
	String() string
}

type env struct {
	outer    Env
	defs     map[string]read.Node
	specials map[string]SpecialForm
}

func NewEnv(outer Env) Env {
	return &env{
		outer:    outer,
		defs:     make(map[string]read.Node),
		specials: make(map[string]SpecialForm),
	}
}

func (e *env) Set(name string, val read.Node) read.Node {
	// if _, ok := e.defs[name]; !ok {
	// 	e.defs[name] = val
	// }
	e.defs[name] = val
	return val
}

// TODO: Iterative version?
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

// TODO: Iterative version?
func (e *env) FindSpecialForm(name string) (SpecialForm, bool) {
	if fun, ok := e.specials[name]; ok {
		return fun, ok
	} else if e.outer != nil {
		return e.outer.FindSpecialForm(name)
	}
	return nil, false
}

func (e *env) String() string {
	var buf bytes.Buffer
	w := print.NewPrinter()
	buf.WriteString("- DEFS ---------------------------------\n")
	for k, v := range e.defs {
		buf.WriteString("  ")
		buf.WriteString(k)
		buf.WriteString(" = ")
		buf.WriteString(w.Print(v))
		buf.WriteString("\n")
	}
	return buf.String()
}
