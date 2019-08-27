package eval

import "github.com/mhoertnagl/splis2/internal/data"

func InitCore(e Evaluator) {
	e.AddCoreFun("list", list)
	e.AddCoreFun("list?", isList)
	e.AddCoreFun("+", sum)
	e.AddCoreFun("-", diff)
	e.AddCoreFun("<", eval2f("<", func(n0 float64, n1 float64) data.Node { return n0 < n1 }))
	e.AddCoreFun(">", eval2f(">", func(n0 float64, n1 float64) data.Node { return n0 > n1 }))
	e.AddCoreFun("<=", eval2f("<=", func(n0 float64, n1 float64) data.Node { return n0 <= n1 }))
	e.AddCoreFun(">=", eval2f(">=", func(n0 float64, n1 float64) data.Node { return n0 >= n1 }))
}

func list(e Evaluator, env data.Env, ns []data.Node) data.Node {
	return data.NewList(ns)
}

func isList(e Evaluator, env data.Env, ns []data.Node) data.Node {
	if len(ns) != 1 {
		return e.Error("")
	}
	return data.IsList(ns[0])
}

// evalSum computes the sum of all arguments.
// Non-numeric arguments will be ignored.
func sum(e Evaluator, env data.Env, ns []data.Node) data.Node {
	var sum float64
	for _, n := range ns {
		m := e.EvalEnv(env, n)
		if v, ok := m.(float64); ok {
			sum += v
		} else {
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("[%s] is not a number.", "")
		}
	}
	return sum
}

func diff(e Evaluator, env data.Env, ns []data.Node) data.Node {
	len := len(ns)
	if len == 1 {
		v := e.EvalEnv(env, ns[0])
		if n, ok := v.(float64); ok {
			return -n
		} else {
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("[%s] is not a number.", "")
		}
	}
	if len == 2 {
		v1 := e.EvalEnv(env, ns[0])
		v2 := e.EvalEnv(env, ns[1])
		if n1, ok1 := v1.(float64); ok1 {
			if n2, ok2 := v2.(float64); ok2 {
				return n1 - n2
			}
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("[%s] is not a number.", "")
		}
		// TODO: Add a printer instance to the evaluator to print expressions.
		return e.Error("[%s] is not a number.", "")
	}
	return e.Error("- requires either 1 or 2 arguments.")
}

// func (e *evaluator) eval1n(f func(float64, float64) data.Node) data.SpecialForm {
// 	return func(e.Env() data.Env, name string, ns []data.Node) data.Node {
// 		if len(ns) != 2 {
// 			return e.error("%s requires exactly 2 numeric arguments.", name)
// 		}
// 		v0 := e.EvalEnv(e.Env(), ns[0])
// 		v1 := e.EvalEnv(e.Env(), ns[1])
// 		if n0, ok0 := v0.(float64); ok0 {
// 			if n1, ok1 := v1.(float64); ok1 {
// 				return f(n0, n1)
// 			}
// 		}
// 		// TODO: Angeben welches argument kein float ist.
// 		return e.error("")
// 	}
// }

func eval2f(name string, f func(float64, float64) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, ns []data.Node) data.Node {
		if len(ns) != 2 {
			return e.Error("%s requires exactly 2 numeric arguments.", name)
		}
		v0 := e.EvalEnv(env, ns[0])
		v1 := e.EvalEnv(env, ns[1])
		if n0, ok0 := v0.(float64); ok0 {
			if n1, ok1 := v1.(float64); ok1 {
				return f(n0, n1)
			}
		}
		// TODO: Angeben welches argument kein float ist.
		return e.Error("")
	}
}
