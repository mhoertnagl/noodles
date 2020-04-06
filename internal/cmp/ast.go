package cmp

import "fmt"

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

func NewError(format string, args ...interface{}) *ErrorNode {
	return &ErrorNode{Msg: fmt.Sprintf(format, args...)}
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

func IsCall(n *ListNode, fn string) bool {
	if x, ok := n.Items[0].(*SymbolNode); ok && x.Name == fn {
		return true
	}
	return false
}

func NewList(items []Node) *ListNode {
	return &ListNode{Items: items}
}

func NewList2(items ...Node) *ListNode {
	return &ListNode{Items: items}
}

func Quote(n Node) *ListNode {
	return NewList2(NewSymbol("quote"), n)
}

func Unquote(n Node) *ListNode {
	return NewList2(NewSymbol("unquote"), n)
}

func Dissolve(n Node) *ListNode {
	return NewList2(NewSymbol("dissolve"), n)
}

func Fn(args []Node, body Node) *ListNode {
	return NewList2(NewSymbol("fn"), args, body)
}

type Map map[string]Node
