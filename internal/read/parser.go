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

func New() Parser {
	return &parser{}
}

func (p *parser) Parse(r Reader) Node {
	p.r = r
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
		fmt.Printf("%d: Unexpected [%s]. Expecting [%s]", p.r.Pos(), p.tok, exp)
	}
}

func (p *parser) error() Node {
	return &ErrorNode{}
}

func (p *parser) parse() Node {
	p.next()
	switch {
	case p.tok == ")":
		fmt.Printf("Unexpected [)].\n")
		return p.error()
	case p.tok == "(":
		return p.parseList()
	default:
		return p.parseAtom()
	}
}

func (p *parser) parseList() Node {
	n := &ListNode{}
	p.consume("(")
	for p.tok != ")" && p.tok != "" {
		p.parse()
	}
	p.consume(")")
	return n
}

func (p *parser) parseAtom() Node {
	pre := p.tok[0]
	switch {
	// case strings.HasPrefix(p.tok, `"`):
	case pre == '"':
		return p.parseString()
	case isNumber(pre):
		return p.parseNumber()
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
	if v, err := strconv.ParseFloat(p.tok, 64); err != nil {
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
