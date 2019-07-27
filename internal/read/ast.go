package read

type Node interface {
}

type ErrorNode struct {
}

type ListNode struct {
	Items []Node
}

type VectorNode struct {
	Items []Node
}

type HashMapNode struct {
	Items []Node
}

type StringNode struct {
	Val string
}

type NumberNode struct {
	Val float64
}

type SymbolNode struct {
	Name string
}
