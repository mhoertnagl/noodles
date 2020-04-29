package cmp

import "fmt"

type Node interface {
}

type ErrorNode struct {
	Msg string
}

func NewError(format string, args ...interface{}) *ErrorNode {
	return &ErrorNode{Msg: fmt.Sprintf(format, args...)}
}

func IsError(n Node) bool {
	_, ok := n.(*ErrorNode)
	return ok
}

func IsNil(n Node) bool {
	return n == nil
}

func IsBool(n Node) bool {
	_, ok := n.(bool)
	return ok
}

func IsInteger(n Node) bool {
	_, ok := n.(int64)
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

func NewSymbol(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

func IsSymbol(n Node) bool {
	_, ok := n.(*SymbolNode)
	return ok
}

func (s *SymbolNode) String() string {
	return s.Name
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

func (l *ListNode) Empty() bool {
	return len(l.Items) == 0
}

func (l *ListNode) Len() int {
	return len(l.Items)
}

func (l *ListNode) First() Node {
	return l.Items[0]
}

func (l *ListNode) Rest() []Node {
	return l.Items[1:]
}

func IsList(n Node) bool {
	_, ok := n.(*ListNode)
	return ok
}

func IsCallN(n Node, fn string) bool {
	if x, ok := n.(*ListNode); ok {
		return IsCall(x, fn)
	}
	return false
}

func IsCall(n *ListNode, fn string) bool {
	if len(n.Items) == 0 {
		return false
	}
	if x, ok := n.Items[0].(*SymbolNode); ok {
		return x.Name == fn
	}
	return false
}

func Call(name string, n Node) *ListNode {
	return NewList2(NewSymbol(name), n)
}

func CallVar(name string, n ...Node) *ListNode {
	m := []Node{NewSymbol(name)}
	m = append(m, n...)
	return NewList(m)
}

func Quote(n Node) *ListNode {
	return Call("quote", n)
}

func Unquote(n Node) *ListNode {
	return Call("unquote", n)
}

func Dissolve(n Node) *ListNode {
	return Call("dissolve", n)
}

func Fn(args []Node, body Node) *ListNode {
	return NewList2(NewSymbol("fn"), args, body)
}

type Map map[string]Node
