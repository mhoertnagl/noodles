package rwr

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/mhoertnagl/splis2/internal/cmp"
)

type usingsSet map[string]bool

type UseRewriter struct {
	paths  []string
	usings usingsSet
	rdr    *cmp.Reader
	prs    *cmp.Parser
	err    []string
}

func NewUseRewriter(paths []string) *UseRewriter {
	return &UseRewriter{
		paths:  paths,
		usings: usingsSet{},
		rdr:    cmp.NewReader(),
		prs:    cmp.NewParser(),
		err:    make([]string, 0),
	}
}

func (r *UseRewriter) Errors() []string {
	return r.err
}

func (r *UseRewriter) error(format string, args ...interface{}) {
	e := fmt.Sprintf(format, args...)
	r.err = append(r.err, e)
}

func (r *UseRewriter) Rewrite(n cmp.Node) cmp.Node {
	switch x := n.(type) {
	case []cmp.Node:
		return RewriteItems(r, x)
	case *cmp.ListNode:
		return r.rewriteList(x)
	default:
		return n
	}
}

func (r *UseRewriter) rewriteList(n *cmp.ListNode) cmp.Node {
	if len(n.Items) == 0 {
		return n
	}
	// TODO: Length of items should be 2 (use "...")
	if cmp.IsCall(n, "use") {
		if mod, ok := n.Items[1].(string); ok {
			if r.usings[mod] {
				// File has already been included. Skip.
				return nil
			}
			r.usings[mod] = true
			return r.loadUse(mod)
		}
	}
	return cmp.NewList(RewriteItems(r, n.Items))
}

func (r *UseRewriter) loadUse(mod string) cmp.Node {
	s := r.loadModule(r.paths, mod)
	r.rdr.Load(s)
	c := r.prs.Parse(r.rdr)
	return r.Rewrite(c)
}

func (r *UseRewriter) loadModule(dirs []string, mod string) string {
	for _, dir := range dirs {
		modBytes, err := ioutil.ReadFile(path.Join(dir, mod+".splis"))
		if err == nil {
			return string(modBytes)
		}
	}
	r.error("Could not find module [%s] in %v.", mod, dirs)
	return ""
}
