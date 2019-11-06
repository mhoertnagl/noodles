package eval

import (
	"reflect"

	"github.com/mhoertnagl/splis2/internal/data"
)

func InitCore(e Evaluator) {
	e.AddCoreFun("nil?", eval1n("nil?", isNil))
	e.AddCoreFun("list", list)
	e.AddCoreFun("list?", eval1n("list?", isList))
	e.AddCoreFun("count", count)
	e.AddCoreFun("empty?", eval1n("empty?", isEmpty))
	// TODO: (vector )
	// TODO: (vector? )
	// TODO: (dict? )
	// TODO: (dict )
	e.AddCoreFun("::", cons)
	e.AddCoreFun(":::", concat)
	e.AddCoreFun("head", eval1n("head", head))
	e.AddCoreFun("tail", eval1n("tail", tail))
	// TODO: (join <list/vector> <list/vector>)
	// TODO: (join <string> <string>)
	// TODO: (join <dict> <dict>)
	// TODO: (print ...)
	// e.AddCoreFun("str", printArgs(false))
	e.AddCoreFun("+", evalxf("+", sum))
	e.AddCoreFun("-", eval12f("-", neg, diff))
	e.AddCoreFun("*", evalxf("*", prod))
	e.AddCoreFun("/", eval12f("/", reciproc, div))
	e.AddCoreFun("=", eq)
	e.AddCoreFun("<", eval2f("<", lt))
	e.AddCoreFun(">", eval2f(">", gt))
	e.AddCoreFun("<=", eval2f("<=", le))
	e.AddCoreFun(">=", eval2f(">=", ge))
}

func sum(acc, v float64) float64   { return acc + v }
func neg(n float64) data.Node      { return -n }
func diff(a, b float64) data.Node  { return a - b }
func prod(acc, v float64) float64  { return acc * v }
func reciproc(n float64) data.Node { return 1 / n }
func div(a, b float64) data.Node   { return a / b }

func lt(a, b float64) data.Node { return a < b }
func gt(a, b float64) data.Node { return a > b }
func le(a, b float64) data.Node { return a <= b }
func ge(a, b float64) data.Node { return a >= b }

func isNil(e Evaluator, env data.Env, arg data.Node) data.Node {
	return data.IsNil(arg)
}

func list(e Evaluator, env data.Env, args []data.Node) data.Node {
	return data.NewList(args)
}

func isList(e Evaluator, env data.Env, arg data.Node) data.Node {
	return data.IsList(arg)
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

func isEmpty(e Evaluator, env data.Env, arg data.Node) data.Node {
	switch x := arg.(type) {
	case *data.ListNode:
		return len(x.Items) == 0
	case *data.VectorNode:
		return len(x.Items) == 0
	case *data.HashMapNode:
		return len(x.Items) == 0
	default:
		return e.Error("[%s] cannot be an argument to [empty?].", "")
	}
}

func cons(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 2 {
		return e.Error("[::] expects 2 arguments.")
	}
	switch x := args[1].(type) {
	case *data.ListNode:
		ns := cons2(args[0], x.Items)
		return data.NewList(ns)
	case *data.VectorNode:
		ns := cons2(args[0], x.Items)
		return data.NewVector(ns)
	default:
		return e.Error("Second argument to [::] must be a list or vector.")
	}
}

func cons2(hd data.Node, tl []data.Node) []data.Node {
	ns := make([]data.Node, len(tl)+1)
	ns[0] = hd
	for i, a := range tl {
		ns[i+1] = a
	}
	return ns
}

func concat(e Evaluator, env data.Env, args []data.Node) data.Node {
	ns := []data.Node{}
	for _, arg := range args {
		switch x := arg.(type) {
		case *data.ListNode:
			ns = append(ns, x.Items...)
			// TODO: Allow? Return list?
		// case *data.VectorNode:
		// 	ns = append(ns, x.Items...)
		default:
			return e.Error("Second argument to [cons] must be a list or vector.")
		}
	}
	return data.NewList(ns)
}

func head(e Evaluator, env data.Env, arg data.Node) data.Node {
	switch x := arg.(type) {
	case *data.ListNode:
		if len(x.Items) == 0 {
			return e.Error("Argument to [head] cannot be the empty list.")
		}
		return x.Items[0]
	case *data.VectorNode:
		if len(x.Items) == 0 {
			return e.Error("Argument to [head] cannot be the empty vector.")
		}
		return x.Items[0]
	default:
		return e.Error("Argument to [head] must be a list or vector.")
	}
}

func tail(e Evaluator, env data.Env, arg data.Node) data.Node {
	switch x := arg.(type) {
	case *data.ListNode:
		ln := len(x.Items)
		if ln == 0 {
			return e.Error("Argument to [tail] cannot be the empty list.")
		}
		return data.NewList(x.Items[1:ln])
	case *data.VectorNode:
		ln := len(x.Items)
		if ln == 0 {
			return e.Error("Argument to [tail] cannot be the empty vector.")
		}
		return data.NewVector(x.Items[1:ln])
	default:
		return e.Error("Argument to [tail] must be a list or vector.")
	}
}

// func printArgs(escape bool) CoreFun {
// 	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
// 		var buf bytes.Buffer
// 		for _, arg := range args {
// 			printArg(&buf, arg, escape)
// 		}
// 		return buf.String()
// 	}
// }
//
// func printArg(buf *bytes.Buffer, n data.Node, escape bool) {
// 	switch {
// 	case data.IsError(n):
// 		buf.WriteString("  [ERROR]  ")
// 	case data.IsNil(n):
// 		buf.WriteString("nil")
// 	case data.IsBool(n):
// 		buf.WriteString(strconv.FormatBool(n.(bool)))
// 	case data.IsNumber(n):
// 		buf.WriteString(strconv.FormatFloat(n.(float64), 'f', -1, 64))
// 	case data.IsString(n):
// 		if escape {
// 			buf.WriteString(n.(string))
// 		} else {
// 			buf.WriteString(n.(string))
// 		}
// 	case data.IsSymbol(n):
// 		buf.WriteString(n.(*data.SymbolNode).Name)
// 	case data.IsList(n):
// 		printSeq(buf, n.(*data.ListNode).Items, "(", ")", escape)
// 	case data.IsVector(n):
// 		printSeq(buf, n.(*data.VectorNode).Items, "[", "]", escape)
// 	case data.IsHashMap(n):
// 		printHashMap(buf, n.(*data.HashMapNode).Items, escape)
// 		// case data.IsFuncNode(n):
// 		// 	p.buf.WriteString(n.(*data.FuncNode).Name)
// 	}
// }
//
// func printSeq(buf *bytes.Buffer, items []data.Node, start string, end string, escape bool) {
// 	buf.WriteString(start)
// 	for i, item := range items {
// 		if i > 0 {
// 			buf.WriteString(" ")
// 		}
// 		printArg(buf, item, escape)
// 	}
// 	buf.WriteString(end)
// }
//
// func printHashMap(buf *bytes.Buffer, items data.Map, escape bool) {
// 	buf.WriteString("{")
// 	// TODO: Unfortunate.
// 	init := false
// 	for key, val := range items {
// 		if init {
// 			buf.WriteString(" ")
// 		}
// 		init = true
// 		printArg(buf, key, escape)
// 		buf.WriteString(" ")
// 		printArg(buf, val, escape)
// 	}
// 	buf.WriteString("}")
// }

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
	case *data.SymbolNode:
		y := b.(*data.SymbolNode)
		return x.Name == y.Name
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

func eval1n(name string, f func(Evaluator, data.Env, data.Node) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		switch len(args) {
		case 1:
			return f(e, env, args[0])
		}
		return e.Error("[%s] expects 1 arguments.", name)
	}
}

func evalxf(name string, f func(float64, float64) float64) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		var acc float64
		for i, arg := range args {
			if v, ok := arg.(float64); ok {
				acc = f(acc, v)
			} else {
				// TODO: Add a printer instance to the evaluator to print expressions.
				return e.Error("[%d]. argument [%s] is not a number.", i+1, "")
			}
		}
		return acc
	}
}

func eval12f(name string, f func(float64) data.Node, g func(float64, float64) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		switch len(args) {
		case 1:
			if n, ok := args[0].(float64); ok {
				return f(n)
			}
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("Argument [%s] is not a number.", "")
		case 2:
			n1, ok1 := args[0].(float64)
			n2, ok2 := args[1].(float64)
			if ok1 && ok2 {
				return g(n1, n2)
			}
			if !ok1 {
				return e.Error("First argument [%s] is not a number.", "")
			}
			if !ok2 {
				return e.Error("Second argument [%s] is not a number.", "")
			}
		}
		return e.Error("[%s] requires either 1 or 2 arguments.", name)
	}
}

func eval2f(name string, f func(float64, float64) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		switch len(args) {
		case 2:
			n0, ok0 := args[0].(float64)
			n1, ok1 := args[1].(float64)
			if ok0 && ok1 {
				return f(n0, n1)
			}
			if !ok0 {
				return e.Error("First argument [%s] is not a number.", "")
			}
			if !ok1 {
				return e.Error("Second argument [%s] is not a number.", "")
			}
		}
		return e.Error("[%s] expects 2 arguments.", name)
	}
}
