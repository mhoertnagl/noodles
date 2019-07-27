package eval

type Env interface {
	Set()
	Lookup()
}

type env struct {
	outer Env
}

func NewEnv() Env {
	return &env{}
}
