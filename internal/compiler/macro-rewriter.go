package compiler

import "fmt"

type macroDefs map[string]*macroDef

type macroDef struct {
	pars []string
	body Node
}

type macroRewriter struct {
	macros macroDefs
}

func NewMacroRewriter() *macroRewriter {
	return &macroRewriter{
		macros: macroDefs{},
	}
}

func (r *macroRewriter) addMacro(name Node, pars Node, body Node) {
	if sym, ok := name.(*SymbolNode); ok {
		if _, found := r.macros[sym.Name]; found {
			panic(fmt.Sprintf("[defmacro] macro [%s] redefined", sym.Name))
		}
		r.macros[sym.Name] = &macroDef{
			pars: getParamNames(pars),
			body: body,
		}
	} else {
		panic(fmt.Sprintf("[defmacro] argument 1 has to be a symbol but is [%T]", name))
	}
}

func getParamNames(parsNode Node) []string {
	if pars, ok := parsNode.(*VectorNode); ok {
		names := make([]string, len(pars.Items))
		for i, par := range pars.Items {
			switch sym := par.(type) {
			case *SymbolNode:
				names[i] = sym.Name
			default:
				panic(fmt.Sprintf("[defmacro] parameter [%d] is not a symbol", i))
			}
		}
		return names
	}
	panic(fmt.Sprintf("[defmacro] argument 2 has to be a vector of symbols"))
}

func (r *macroRewriter) Rewrite(n Node) Node {
	switch x := n.(type) {
	case bool:
		return r.rewriteBoolean(x)
	case int64:
		return r.rewriteInteger(x)
	case string:
		return r.rewriteString(x)
	case *SymbolNode:
		return r.rewriteSymbol(x)
	case *VectorNode:
		return r.rewriteVector(x)
	case *ListNode:
		return r.rewriteList(x)
	}
	panic(fmt.Sprintf("Macro-Rewriter: Unsupported node [%v:%T]", n, n))
}

func (r *macroRewriter) rewriteBoolean(n bool) Node {
	return n
}

func (r *macroRewriter) rewriteInteger(n int64) Node {
	return n
}

func (r *macroRewriter) rewriteString(n string) Node {
	return n
}

func (r *macroRewriter) rewriteSymbol(n *SymbolNode) Node {
	return n
}

func (r *macroRewriter) rewriteVector(n *VectorNode) Node {
	return n
}

func (r *macroRewriter) rewriteList(n *ListNode) Node {
	if len(n.Items) == 0 {
		return n
	}
	switch x := n.Items[0].(type) {
	case *SymbolNode:
		switch x.Name {
		case "defmacro":
			r.addMacro(n.Items[1], n.Items[2], n.Items[3])
			return nil
		default:
			if def, ok := r.macros[x.Name]; ok {
				rw := NewArgsRewriter(def.pars, n.Items[1:])
				return r.Rewrite(rw.Rewrite(def.body))
			}
		}
	}
	return NewList(r.rewriteItems(n.Items))
}

func (r *macroRewriter) rewriteItems(ns []Node) []Node {
	ms := []Node{}
	for _, n := range ns {
		// If the list item is a defmacro instruction, rewrite will return nil.
		m := r.Rewrite(n)
		if m != nil {
			ms = append(ms, m)
		}
	}
	return ms
}
