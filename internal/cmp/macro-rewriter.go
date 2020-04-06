package cmp

import "fmt"

type macroDefs map[string]*macroDef

type macroDef struct {
	pars []string
	body Node
}

type MacroRewriter struct {
	macros macroDefs
}

func NewMacroRewriter() *MacroRewriter {
	return &MacroRewriter{
		macros: macroDefs{},
	}
}

func (r *MacroRewriter) Rewrite(n Node) Node {
	switch x := n.(type) {
	case []Node:
		return RewriteItems(r, x)
	case *ListNode:
		return r.rewriteList(x)
	default:
		return n
	}
}

func (r *MacroRewriter) rewriteList(n *ListNode) Node {
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
	return NewList(RewriteItems(r, n.Items))
}

func (r *MacroRewriter) addMacro(name Node, pars Node, body Node) {
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
	if pars, ok := parsNode.([]Node); ok {
		names := make([]string, len(pars))
		for i, par := range pars {
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
