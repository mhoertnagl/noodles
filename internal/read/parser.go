package read

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mhoertnagl/splis2/internal/data"
)

type Parser interface {
	Parse(r Reader) data.Node
	Errors() []*data.ErrorNode
}

// TODO: Create a type Module this will contain the ast, errors and other infos.

type parser struct {
	rd  Reader
	tok string
	err []*data.ErrorNode
}

func NewParser() Parser {
	return &parser{}
}

func (p *parser) Parse(r Reader) data.Node {
	p.rd = r
	p.err = []*data.ErrorNode{}
	p.next()
	return p.parse()
}

func (p *parser) Errors() []*data.ErrorNode {
	return p.err
}

func (p *parser) next() {
	p.tok = p.rd.Next()
}

func (p *parser) consume(exp string) {
	if p.tok == exp {
		p.next()
	} else {
		p.error("Unexpected [%s]. Expecting [%s].\n", p.tok, exp)
	}
}

func (p *parser) error(format string, args ...interface{}) data.Node {
	e := data.NewError(fmt.Sprintf(format, args...))
	p.err = append(p.err, e)
	p.next() // Ignore the malign token and move on.
	return e
}

func (p *parser) parse() data.Node {
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

func (p *parser) parseList() data.Node {
	return data.NewList(p.parseArgs("(", ")"))
}

func (p *parser) parseVector() data.Node {
	return data.NewVector(p.parseArgs("[", "]"))
}

func (p *parser) parseHashMap() data.Node {
	n := data.NewHashMap2()
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
func (p *parser) parseArgs(start string, end string) []data.Node {
	args := []data.Node{}
	p.consume(start)
	for p.tok != end && p.tok != "" {
		args = append(args, p.parse())
	}
	// TODO: Should we return an error node?
	p.consume(end)
	return args
}

func (p *parser) parseAtom() data.Node {
	var n data.Node
	switch {
	case strings.HasPrefix(p.tok, `"`):
		n = p.parseString()
	case isNumber(p.tok):
		n = p.parseNumber()
	case p.tok == "true":
		n = true // TrueObject
	case p.tok == "false":
		n = false // FalseObject
	case p.tok == "nil":
		n = nil // NilObject
	default:
		n = p.parseSymbol()
	}
	p.next()
	return n
}

func (p *parser) parseString() data.Node {
	if strings.HasSuffix(p.tok, `"`) {
		// TODO: Create a constant for the empty string.
		return normalizeString(p.tok) // NewString(normalizeString(p.tok))
	}
	return p.error("Missing [\"].\n")
}

func normalizeString(val string) string {
	val = strings.Trim(val, `"`)
	val = strings.Replace(val, `\n`, "\n", -1)
	return val
}

func (p *parser) parseNumber() data.Node {
	if v, err := strconv.ParseFloat(p.tok, 64); err == nil {
		// TODO: Create a constant pool for the numbers [-32, 31]?
		return v // NewNumber(v)
	}
	return p.error("[%s] is not a floating point number.\n", p.tok)
}

func isNumber(tok string) bool {
	return len(tok) > 0 && '0' <= tok[0] && tok[0] <= '9'
}

// TODO: Keywords :<x> <-> ʞ<x>

func (p *parser) parseSymbol() data.Node {
	return data.NewSymbol(p.tok)
}
