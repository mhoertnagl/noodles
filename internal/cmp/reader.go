package cmp

import (
	"bytes"
	"regexp"
)

// Reader tokenizes the input string and provides methods to enumerate the
// tokens sequentially.
type Reader struct {
	re     *regexp.Regexp
	tokens []string
	pos    int
}

// NewReader creates a new Reader instance.
func NewReader() *Reader {
	r := &Reader{}
	pat := buildPattern()
	r.re = regexp.MustCompile(pat)
	return r
}

func (r *Reader) Load(input string) {
	mm := r.re.FindAllStringSubmatch(input, -1)
	r.tokens = []string{}
	for _, m := range mm {
		if m[1] != "" {
			r.tokens = append(r.tokens, m[1])
		}
	}
	r.pos = 0
}

func (r *Reader) Next() string {
	t := r.Peek()
	r.pos++
	return t
}

func (r *Reader) Peek() string {
	if r.pos < len(r.tokens) {
		return r.tokens[r.pos]
	}
	return ""
}

func (r *Reader) Pos() int {
	// Positions visible to the user start at 1.
	return r.pos + 1
}

// https://regex101.com/r/Awgqpk/1
func buildPattern() string {
	var pat bytes.Buffer
	pat.WriteString(`[\s,]*`)                     // whitespace or commas
	pat.WriteString("(")                          // Begin capture group
	pat.WriteString("[\\[\\]{}\\(\\)'~^@]")       // any of [, ], {, }, (, ), ', ~, ^, @
	pat.WriteString("|")                          // or
	pat.WriteString(`"(?:\\"|[^"])*"?`)           // strings with escape characters and an optional " at the end
	pat.WriteString("|")                          // or
	pat.WriteString("[^\\s\\[\\]{}\\('\",;\\)]+") // symbols (including numbers)
	pat.WriteString(")")                          // End capture group
	pat.WriteString("|")                          // or
	pat.WriteString(";[^\n]*(?:$|\n)")            // comments
	return pat.String()
}
