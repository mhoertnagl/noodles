package cmp

import (
	"strconv"
	"strings"
)

type Parser struct {
	rd  *Reader
	tok string
	err []*ErrorNode
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r *Reader) Node {
	p.rd = r
	p.err = []*ErrorNode{}
	p.next()
	return p.parse()
}

func (p *Parser) Errors() []*ErrorNode {
	return p.err
}

func (p *Parser) next() {
	p.tok = p.rd.Next()
}

func (p *Parser) consume(exp string) {
	if p.tok == exp {
		p.next()
	} else {
		p.error("Unexpected [%s]. Expecting [%s].\n", p.tok, exp)
	}
}

func (p *Parser) error(format string, args ...interface{}) Node {
	e := NewError(format, args...)
	p.err = append(p.err, e)
	p.next() // Ignore the malign token and move on.
	return e
}

func (p *Parser) parse() Node {
	switch {
	case p.tok == ")":
		return p.error("Unexpected [)].\n")
	case p.tok == "(":
		return p.parseList()
	case p.tok == "]":
		return p.error("Unexpected []].\n")
	case p.tok == "[":
		return p.parseArgs("[", "]")
	case p.tok == "}":
		return p.error("Unexpected [}].\n")
	case p.tok == "{":
		return p.parseHashMap()
	case p.tok == "'":
		p.consume("'")
		return Quote(p.parse())
	case p.tok == "~":
		p.consume("~")
		return Unquote(p.parse())
	case p.tok == "@":
		p.consume("@")
		return Dissolve(p.parse())
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseList() Node {
	return NewList(p.parseArgs("(", ")"))
}

func (p *Parser) parseHashMap() Node {
	n := Map{}
	p.consume("{")
	// TODO: parse first and check if the count is even.
	for p.tok != "}" && p.tok != "" {
		key := p.parse()
		if k, ok := key.(string); ok {
			v := p.parse()
			n[k] = v
		}
	}
	p.consume("}")
	return n
}

// TODO: link to parent node to provide more context for errorNodes.
func (p *Parser) parseArgs(start string, end string) []Node {
	args := []Node{}
	p.consume(start)
	for p.tok != end && p.tok != "" {
		args = append(args, p.parse())
	}
	// TODO: Should we return an error node?
	p.consume(end)
	return args
}

func (p *Parser) parseAtom() Node {
	var n Node
	switch {
	case strings.HasPrefix(p.tok, `"`):
		n = p.parseString()
	case isNumber(p.tok):
		n = p.parseNumber()
	case p.tok == "true":
		n = true
	case p.tok == "false":
		n = false
	case p.tok == "nil":
		n = nil
	default:
		n = p.parseSymbol()
	}
	p.next()
	return n
}

func (p *Parser) parseString() Node {
	if strings.HasSuffix(p.tok, `"`) {
		// TODO: Create a constant for the empty string.
		return normalizeString(p.tok)
	}
	return p.error("Missing [\"].\n")
}

func normalizeString(val string) string {
	val = strings.Trim(val, `"`)
	val = strings.Replace(val, `\n`, "\n", -1)
	return val
}

func (p *Parser) parseNumber() Node {
	if v, err := strconv.ParseInt(p.tok, 10, 64); err == nil {
		// TODO: Create a constant pool for the numbers [-32, 31]?
		return v
	}
	if v, err := strconv.ParseFloat(p.tok, 64); err == nil {
		return v
	}
	return p.error("[%s] is not a number.\n", p.tok)
}

func isNumber(tok string) bool {
	return len(tok) > 0 && '0' <= tok[0] && tok[0] <= '9'
}

// TODO: Keywords :<x> <-> ʞ<x>

func (p *Parser) parseSymbol() Node {
	return NewSymbol(p.tok)
}
