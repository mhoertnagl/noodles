package read

// TODO: Move to separate namespace?

type Node interface {
}

type ErrorNode struct {
	Msg string
}

// TODO: Make it a varargs version with fmt.
func NewError(msg string) *ErrorNode {
	return &ErrorNode{Msg: msg}
}

type NilNode struct{}

var NilObject *NilNode = &NilNode{}

type TrueNode struct{}
type FalseNode struct{}

var TrueObject *TrueNode = &TrueNode{}
var FalseObject *FalseNode = &FalseNode{}

func NewBool(val bool) Node {
	if val {
    return TrueObject
  }
  return FalseObject
}

type StringNode struct {
	Val string
}

func NewString(val string) *StringNode {
	return &StringNode{Val: val}
}

type NumberNode struct {
	Val float64
}

func NewNumber(val float64) *NumberNode {
	return &NumberNode{Val: val}
}

type SymbolNode struct {
	Name string
}

func NewSymbol(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

type ListNode struct {
	Items []Node
}

func NewList(items []Node) *ListNode {
	return &ListNode{Items: items}
}

func NewList2(items ...Node) *ListNode {
	return &ListNode{Items: items}
}

type VectorNode struct {
	Items []Node
}

func NewVector(items []Node) *VectorNode {
	return &VectorNode{Items: items}
}

func NewVector2(items ...Node) *VectorNode {
	return &VectorNode{Items: items}
}

type HashMapNode struct {
	Items map[Node]Node
}

func NewHashMap(items map[Node]Node) *HashMapNode {
	return &HashMapNode{Items: items}
}

// TODO: Rename to NewEmptyHashMap()
func NewHashMap2() *HashMapNode {
	return NewHashMap(make(map[Node]Node))
}
