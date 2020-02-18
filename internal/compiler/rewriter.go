package compiler

type Rewriter interface {
	Rewrite(n Node) Node
}
