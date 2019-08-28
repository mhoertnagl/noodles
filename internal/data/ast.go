package data

// TODO: Move to separate namespace?

type Node interface {
}

type ErrorNode struct {
	Msg string
}

func IsError(n Node) bool {
	_, ok := n.(*ErrorNode)
	return ok
}

// TODO: Make it a varargs version with fmt.
func NewError(msg string) *ErrorNode {
	return &ErrorNode{Msg: msg}
}

func IsNil(n Node) bool {
	return n == nil
}

func IsBool(n Node) bool {
	_, ok := n.(bool)
	return ok
}

func IsNumber(n Node) bool {
	_, ok := n.(float64)
	return ok
}

func IsString(n Node) bool {
	_, ok := n.(string)
	return ok
}

type SymbolNode struct {
	Name string
}

func IsSymbol(n Node) bool {
	_, ok := n.(*SymbolNode)
	return ok
}

func NewSymbol(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

type ListNode struct {
	Items []Node
}

func IsList(n Node) bool {
	_, ok := n.(*ListNode)
	return ok
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

func IsVector(n Node) bool {
	_, ok := n.(*VectorNode)
	return ok
}

func NewVector(items []Node) *VectorNode {
	return &VectorNode{Items: items}
}

func NewVector2(items ...Node) *VectorNode {
	return &VectorNode{Items: items}
}

type Map map[string]Node

type HashMapNode struct {
	Items Map
}

func IsHashMap(n Node) bool {
	_, ok := n.(*HashMapNode)
	return ok
}

func NewHashMap(items Map) *HashMapNode {
	return &HashMapNode{Items: items}
}

func NewEmptyHashMap() *HashMapNode {
	return NewHashMap(make(Map))
}

type FuncNode struct {
	Env  Env
	Pars []string
	Fun  Node
}

func NewFuncNode(env Env, pars []string, fun Node) Node {
	return &FuncNode{Env: env, Pars: pars, Fun: fun}
}

func IsFuncNode(n Node) bool {
	_, ok := n.(*FuncNode)
	return ok
}
