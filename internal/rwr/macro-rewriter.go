package rwr

import (
	"fmt"

	"github.com/mhoertnagl/splis2/internal/cmp"
	"github.com/mhoertnagl/splis2/internal/util"
)

type macroDefs map[string]*macroDef

type macroDef struct {
	man  []string
	opt  string
	body cmp.Node
}

type MacroRewriter struct {
	macros macroDefs
}

func NewMacroRewriter() *MacroRewriter {
	return &MacroRewriter{
		macros: macroDefs{},
	}
}

func (r *MacroRewriter) Rewrite(n cmp.Node) cmp.Node {
	switch x := n.(type) {
	case []cmp.Node:
		return RewriteItems(r, x)
	case *cmp.ListNode:
		return r.rewriteList(x)
	default:
		return n
	}
}

func (r *MacroRewriter) rewriteList(n *cmp.ListNode) cmp.Node {
	if len(n.Items) == 0 {
		return n
	}
	switch x := n.Items[0].(type) {
	case *cmp.SymbolNode:
		switch x.Name {
		case "defmacro":
			r.addMacro(n.Items[1], n.Items[2], n.Items[3])
			return nil
		default:
			if def, ok := r.macros[x.Name]; ok {
				rw := NewArgsRewriter(def.man, def.opt, n.Items[1:])
				return r.Rewrite(rw.Rewrite(def.body))
			}
		}
	}
	return cmp.NewList(RewriteItems(r, n.Items))
}

func (r *MacroRewriter) addMacro(name cmp.Node, pars cmp.Node, body cmp.Node) {
	if sym, ok := name.(*cmp.SymbolNode); ok {
		r.addMacro2(sym.Name, pars, body)
	} else {
		panic(fmt.Sprintf("[defmacro] argument 1 has to be a symbol but is [%T]", name))
	}
}

func (r *MacroRewriter) addMacro2(name string, pars cmp.Node, body cmp.Node) {
	if _, found := r.macros[name]; found {
		panic(fmt.Sprintf("[defmacro] macro [%s] redefined", name))
	}

	man, opt := getParamNames(pars)

	r.macros[name] = &macroDef{
		man:  man,
		opt:  opt,
		body: body,
	}
}

func getParamNames(parsNode cmp.Node) ([]string, string) {
	if params, ok := parsNode.([]cmp.Node); ok {
		return extractParams(params)
	}
	panic(fmt.Sprintf("[defmacro] argument 2 has to be a vector of symbols"))
}

func extractParams(params []cmp.Node) ([]string, string) {
	names := verifyParams(params)
	pos := util.IndexOf(names, "&")
	if pos == -1 {
		return names, ""
	}
	if len(names) == pos+1 {
		panic("[fn] missing optional parameter")
	}
	if len(names) > pos+2 {
		panic("[fn] excess optional parameter")
	}
	return names[:pos], names[pos+1]
}

func verifyParams(params []cmp.Node) []string {
	names := make([]string, 0)
	for pos, param := range params {
		names = append(names, verifyParam(param, pos))
	}
	return names
}

func verifyParam(param cmp.Node, pos int) string {
	switch sym := param.(type) {
	case *cmp.SymbolNode:
		return sym.Name
	default:
		panic(fmt.Sprintf("[fn] parameter [%d] is not a symbol", pos))
	}
}
