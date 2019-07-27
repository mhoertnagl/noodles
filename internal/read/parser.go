package read

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(r Reader) Node
	Errors() []*ErrorNode
}

// TODO: Create a type Module this will contain the ast, errors and other infos.

type parser struct {
	rd  Reader
	tok string
	err []*ErrorNode
}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(r Reader) Node {
	p.rd = r
	p.err = []*ErrorNode{}
	p.next()
	return p.parse()
}

// func (p *parser) success() bool {
// 	return len(p.err) == 0
// }

func (p *parser) Errors() []*ErrorNode {
	return p.err
}

func (p *parser) next() {
	p.tok = p.rd.Next()
}

func (p *parser) peek() string {
	return p.rd.Peek()
}

func (p *parser) consume(exp string) {
	if p.tok == exp {
		p.next()
	} else {
		fmt.Printf("Unexpected [%s]. Expecting [%s].\n", p.tok, exp)
	}
}

func (p *parser) error(format string, args ...interface{}) Node {
	e := &ErrorNode{Msg: fmt.Sprintf(format, args...)}
	p.err = append(p.err, e)
	p.next() // Ignore the malign token and move on.
	return e
}

func (p *parser) parse() Node {
	switch {
	case p.tok == ")":
		return p.error("Unexpected [)].\n")
	case p.tok == "(":
		return p.parseList()
	case p.tok == "]":
		return p.error("Unexpected []].\n")
	case p.tok == "[":
		return p.parseVector()
	case p.tok == "}":
		return p.error("Unexpected [}].\n")
	case p.tok == "{":
		return p.parseHashMap()
	default:
		return p.parseAtom()
	}
}

func (p *parser) parseList() Node {
	return &ListNode{Items: p.parseArgs("(", ")")}
}

func (p *parser) parseVector() Node {
	return &VectorNode{Items: p.parseArgs("[", "]")}
}

func (p *parser) parseHashMap() Node {
	n := &HashMapNode{Items: make(map[Node]Node)}
	p.consume("{")
	// TODO: parse first and check if the count is even.
	for p.tok != "}" && p.tok != "" {
		key := p.parse()
		val := p.parse()
		n.Items[key] = val
	}
	p.consume("}")
	return n
}

// TODO: link to parent node to provide more context for errorNodes.
func (p *parser) parseArgs(start string, end string) []Node {
	args := []Node{}
	p.consume(start)
	for p.tok != end && p.tok != "" {
		args = append(args, p.parse())
	}
	p.consume(end)
	return args
}

func (p *parser) parseAtom() Node {
	var n Node
	// TODO: Can tok be nil?
	pre := p.tok[0]
	switch {
	case pre == '"':
		n = p.parseString()
	case isNumber(pre):
		n = p.parseNumber()
	case p.tok == "true":
		n = TrueObject
	case p.tok == "false":
		n = FalseObject
	case p.tok == "nil":
		n = NilObject
	default:
		n = p.parseSymbol()
	}
	// TODO: Remove this in final version if we can guarantee that default is symbol for all inputs.
	// This should never happen. Defaults to symbol.
	if n == nil {
		n = p.error("Unrecognized token [%s].", p.tok)
	}
	p.next()
	return n
}

func (p *parser) parseString() Node {
	if strings.HasSuffix(p.tok, `"`) {
		return &StringNode{Val: normalizeString(p.tok)}
	}
	return p.error("Missing [\"].\n")
}

func normalizeString(val string) string {
	val = strings.Trim(val, `"`)
	val = strings.Replace(val, `\n`, "\n", -1)
	return val
}

func (p *parser) parseNumber() Node {
	if v, err := strconv.ParseFloat(p.tok, 64); err == nil {
		return &NumberNode{Val: v}
	}
	return p.error("[%s] is not a floating point number.\n", p.tok)
}

func isNumber(c byte) bool {
	return ('0' <= c && c <= '9') || c == '-' || c == '.'
}

// TODO: Keywords :<x> <-> ʞ<x>

func (p *parser) parseSymbol() Node {
	return &SymbolNode{Name: p.tok}
}
