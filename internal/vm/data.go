package vm

type Val interface{}

type Env map[int64]Val

// TODO: variable type.
type Sym struct {
	name string
	// file string
	// line int
}

type Syms map[int64]Sym
