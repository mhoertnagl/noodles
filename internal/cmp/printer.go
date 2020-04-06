package cmp

import (
	"bytes"
	"strconv"
)

func PrintAst(node Node) string {
	var buf bytes.Buffer
	printNode(&buf, node)
	return buf.String()
}

func printNode(buf *bytes.Buffer, node Node) {
	switch x := node.(type) {
	case *ErrorNode:
		buf.WriteString("  [ERROR]  ")
	case nil:
		buf.WriteString("nil")
	case bool:
		buf.WriteString(strconv.FormatBool(x))
	case int64:
		buf.WriteString(strconv.FormatInt(x, 10))
	// case float64:
	// 	buf.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	case string:
		printString(buf, x)
	case *SymbolNode:
		buf.WriteString(x.Name)
	case *ListNode:
		printSeq(buf, x.Items, "(", ")")
	case []Node:
		printSeq(buf, x, "[", "]")
	case Map:
		printHashMap(buf, x)
	}
}

func printString(buf *bytes.Buffer, s string) {
	buf.WriteString(`"`)
	buf.WriteString(s)
	buf.WriteString(`"`)
}

func printSeq(buf *bytes.Buffer, items []Node, start string, end string) {
	buf.WriteString(start)
	for i, item := range items {
		if i > 0 {
			buf.WriteString(" ")
		}
		printNode(buf, item)
	}
	buf.WriteString(end)
}

func printHashMap(buf *bytes.Buffer, m Map) {
	buf.WriteString("{")
	// TODO: Unfortunate.
	init := false
	for key, val := range m {
		if init {
			buf.WriteString(" ")
		}
		init = true
		printNode(buf, key)
		buf.WriteString(" ")
		printNode(buf, val)
	}
	buf.WriteString("}")
}
