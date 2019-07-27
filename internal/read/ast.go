package read

// TODO: Move to separate namespace?

type Node interface {
}

type ErrorNode struct {
	Msg string
}

type ListNode struct {
	Items []Node
}

type VectorNode struct {
	Items []Node
}

type HashMapNode struct {
	Items map[Node]Node
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

type TrueNode struct{}
type FalseNode struct{}
type NilNode struct{}

var TrueObject *TrueNode = &TrueNode{}
var FalseObject *FalseNode = &FalseNode{}
var NilObject *NilNode = &NilNode{}
