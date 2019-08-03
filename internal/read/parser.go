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

// func (p *parser) peek() string {
// 	return p.rd.Peek()
// }

func (p *parser) consume(exp string) {
	if p.tok == exp {
		p.next()
	} else {
		fmt.Printf("Unexpected [%s]. Expecting [%s].\n", p.tok, exp)
	}
}

func (p *parser) error(format string, args ...interface{}) Node {
	e := NewError(fmt.Sprintf(format, args...))
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
	return NewList(p.parseArgs("(", ")"))
}

func (p *parser) parseVector() Node {
	return NewVector(p.parseArgs("[", "]"))
}

func (p *parser) parseHashMap() Node {
	n := NewHashMap2()
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
	// TODO: Should we return an error node?
	p.consume(end)
	return args
}

func (p *parser) parseAtom() Node {
	var n Node
	switch {
	case strings.HasPrefix(p.tok, `"`):
		n = p.parseString()
	case isNumber(p.tok):
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
	p.next()
	return n
}

func (p *parser) parseString() Node {
	if strings.HasSuffix(p.tok, `"`) {
		// TODO: Create a constant for the empty string.
		return NewString(normalizeString(p.tok))
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
		// TODO: Create a constant pool for the numbers [-32, 31]?
		return NewNumber(v)
	}
	return p.error("[%s] is not a floating point number.\n", p.tok)
}

func isNumber(tok string) bool {
	return len(tok) > 0 && '0' <= tok[0] && tok[0] <= '9'
}

// TODO: Keywords :<x> <-> ʞ<x>

func (p *parser) parseSymbol() Node {
	return NewSymbol(p.tok)
}
