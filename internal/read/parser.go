package read

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser interface {
	Parse(r Reader) Node
}

type parser struct {
	r   Reader
	tok string
}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(r Reader) Node {
	p.r = r
	p.next()
	return p.parse()
}

func (p *parser) next() {
	p.tok = p.r.Next()
}

func (p *parser) peek() string {
	return p.r.Peek()
}

func (p *parser) consume(exp string) {
	if p.tok == exp {
		p.next()
	} else {
		fmt.Printf("%d: Unexpected [%s]. Expecting [%s].\n", p.r.Pos(), p.tok, exp)
	}
}

func (p *parser) error() Node {
	return &ErrorNode{}
}

func (p *parser) parse() Node {
	switch {
	case p.tok == ")":
		fmt.Printf("Unexpected [)].\n")
		return p.error()
	case p.tok == "(":
		return p.parseList()
	case p.tok == "]":
		fmt.Printf("Unexpected []].\n")
		return p.error()
	case p.tok == "[":
		return p.parseVector()
	case p.tok == "}":
		fmt.Printf("Unexpected [}].\n")
		return p.error()
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
	return &HashMapNode{Items: p.parseArgs("{", "}")}
}

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
	pre := p.tok[0]
	switch {
	case pre == '"':
		return p.parseString()
	case isNumber(pre):
		return p.parseNumber()
	case p.tok == "true":
		p.next()
		return TrueObject
	case p.tok == "false":
		p.next()
		return FalseObject
	case p.tok == "nil":
		p.next()
		return NilObject
	default:
		return p.parseSymbol()
	}
}

func (p *parser) parseString() Node {
	if strings.HasSuffix(p.tok, `"`) {
		n := &StringNode{}
		n.Val = strings.Trim(p.tok, `"`)
		n.Val = strings.Replace(n.Val, `\n`, "\n", -1)
		p.next()
		return n
	}
	fmt.Printf("Missing [\"].\n")
	return p.error()
}

func (p *parser) parseNumber() Node {
	if v, err := strconv.ParseFloat(p.tok, 64); err == nil {
		n := &NumberNode{}
		n.Val = v
		p.next()
		return n
	}
	fmt.Printf("[%s] is not a floating point number.\n", p.tok)
	return p.error()
}

func isNumber(c byte) bool {
	return ('0' <= c && c <= '9') || c == '-' || c == '.'
}

// TODO: Keywords :<x> <-> ʞ<x>

func (p *parser) parseSymbol() Node {
	n := &SymbolNode{}
	n.Name = p.tok
	p.next()
	return n
}
