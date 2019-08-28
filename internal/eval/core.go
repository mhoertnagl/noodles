package eval

import (
	"reflect"

	"github.com/mhoertnagl/splis2/internal/data"
)

func InitCore(e Evaluator) {
	e.AddCoreFun("list", list)
	e.AddCoreFun("list?", isList)
	e.AddCoreFun("count", count)
	e.AddCoreFun("+", sum)
	e.AddCoreFun("-", diff)
	e.AddCoreFun("=", eq)
	e.AddCoreFun("<", eval2f("<", func(a, b float64) data.Node { return a < b }))
	e.AddCoreFun(">", eval2f(">", func(a, b float64) data.Node { return a > b }))
	e.AddCoreFun("<=", eval2f("<=", func(a, b float64) data.Node { return a <= b }))
	e.AddCoreFun(">=", eval2f(">=", func(a, b float64) data.Node { return a >= b }))
}

func list(e Evaluator, env data.Env, args []data.Node) data.Node {
	return data.NewList(args)
}

func isList(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 1 {
		return e.Error("[list?] expects 1 argument.")
	}
	return data.IsList(args[0])
}

func count(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 1 {
		return e.Error("[count] expects 1 argument.")
	}
	switch x := args[0].(type) {
	case *data.ListNode:
		return float64(len(x.Items))
	case *data.VectorNode:
		return float64(len(x.Items))
	case *data.HashMapNode:
		return float64(len(x.Items))
	default:
		return e.Error("[%s] cannot be an argument to [count].", "")
	}
}

func eq(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 2 {
		return e.Error("[=] expects 2 arguments.")
	}
	return eq2(e, env, args[0], args[1])
}

func eq2(e Evaluator, env data.Env, a, b data.Node) data.Node {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch x := a.(type) {
	case *data.ListNode:
		y := b.(*data.ListNode)
		return eqSeq(e, env, x.Items, y.Items)
	case *data.VectorNode:
		y := b.(*data.VectorNode)
		return eqSeq(e, env, x.Items, y.Items)
	case *data.HashMapNode:
		y := b.(*data.HashMapNode)
		return eqHashMap(e, env, x.Items, y.Items)
	default:
		return a == b
	}
}

func eqSeq(e Evaluator, env data.Env, as, bs []data.Node) data.Node {
	if len(as) != len(bs) {
		return false
	}
	for i := 0; i < len(as); i++ {
		if eq2(e, env, as[i], bs[i]) == false {
			return false
		}
	}
	return true
}

func eqHashMap(e Evaluator, env data.Env, as, bs data.Map) data.Node {
	if len(as) != len(bs) {
		return false
	}
	for k, va := range as {
		vb, ok := bs[k]
		if !ok {
			return false
		}
		if eq2(e, env, va, vb) == false {
			return false
		}
	}
	return true
}

// evalSum computes the sum of all arguments.
// Non-numeric arguments will be ignored.
func sum(e Evaluator, env data.Env, args []data.Node) data.Node {
	var sum float64
	for _, arg := range args {
		if v, ok := arg.(float64); ok {
			sum += v
		} else {
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("[%s] is not a number.", "")
		}
	}
	return sum
}

func diff(e Evaluator, env data.Env, args []data.Node) data.Node {
	switch len(args) {
	case 1:
		if n, ok := args[0].(float64); ok {
			return -n
		}
		// TODO: Add a printer instance to the evaluator to print expressions.
		return e.Error("[%s] is not a number.", "")
	case 2:
		if n1, ok1 := args[0].(float64); ok1 {
			if n2, ok2 := args[1].(float64); ok2 {
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
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		if len(args) != 2 {
			return e.Error("[%s] expects 2 arguments.", name)
		}
		if n0, ok0 := args[0].(float64); ok0 {
			if n1, ok1 := args[1].(float64); ok1 {
				return f(n0, n1)
			}
		}
		// TODO: Angeben welches argument kein float ist.
		return e.Error("")
	}
}
