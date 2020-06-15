package cmp

import (
	"github.com/mhoertnagl/noodles/internal/vm"
)

type fnDef struct {
	addr uint64
	code vm.Ins
}

type fnDefs []*fnDef

type specFun func([]Node, *SymTable, *Ctx)

type specDefs map[string]specFun

func (d specDefs) add(name string, fn specFun) {
	d[name] = fn
}

type primDef struct {
	name  string
	op    vm.Op
	nargs int
	rev   bool
}

type primDefs map[string]primDef

func (d primDefs) add(name string, op vm.Op, nargs int, rev bool) {
	d[name] = primDef{name: name, op: op, nargs: nargs, rev: rev}
}

type varPrimDef struct {
	name    string
	op      vm.Op
	argsMin int
}

type varPrimDefs map[string]varPrimDef

func (d varPrimDefs) add(name string, op vm.Op, argsMin int) {
	d[name] = varPrimDef{name: name, op: op, argsMin: argsMin}
}

type defMap struct {
	index uint64
	ids   map[string]uint64
	names map[uint64]string
}

func newDefMap() *defMap {
	return &defMap{
		index: 0,
		ids:   make(map[string]uint64),
		names: make(map[uint64]string),
	}
}

func (d *defMap) get(name string) (uint64, bool) {
	id, ok := d.ids[name]
	return id, ok
}

func (d *defMap) add(name string) uint64 {
	next := d.nextID()
	d.ids[name] = next
	d.names[next] = name
	return next
}

// getOrAdd assigns an ID to the name if it does not exist already. In each
// case it returns the (already) assigned ID.
func (d *defMap) getOrAdd(name string) uint64 {
	if id, ok := d.get(name); ok {
		return id
	}
	return d.add(name)
}

func (d *defMap) nextID() uint64 {
	next := d.index
	d.index++
	return next
}
