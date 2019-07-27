package read

import (
	"bytes"
	//"fmt"
	"regexp"
)

// Reader tokenizes the input string and provides methods to enumerate the
// tokens sequentially.
type Reader interface {
	Load(input string)
	Next() string
	Peek() string
	// Tokens() []string
}

type reader struct {
	re     *regexp.Regexp
	tokens []string
	pos    int
}

// New creates a new Reader instance.
func New() Reader {
	r := &reader{}
	pat := buildPattern()
	r.re = regexp.MustCompile(pat)
	return r
}

func (r *reader) Load(input string) {
	mm := r.re.FindAllStringSubmatch(input, -1)
	// fmt.Println(len(mm))
	r.tokens = make([]string, len(mm))
	for i, m := range mm {
		// fmt.Printf("%d: %s\n", i, m[1])
		r.tokens[i] = m[1]
	}
	// r.tokens = r.re.FindAllString(input, -1)
	// fmt.Println(r.tokens)
	// fmt.Println(len(r.tokens))
	r.pos = 0
}

func (r *reader) Next() string {
	t := r.Peek()
	r.pos++
	return t
}

func (r *reader) Peek() string {
	if r.pos < len(r.tokens) {
		return r.tokens[r.pos]
	}
	return ""
}

// func (r *reader) Tokens() []string {
// 	return r.tokens
// }

func buildPattern() string {
	var pat bytes.Buffer
	pat.WriteString("[\\s,]*")                     // whitespace or commas
	pat.WriteString("(")                           // Begin capture group
	pat.WriteString("~@")                          // ~@
	pat.WriteString("|")                           // or
	pat.WriteString("[\\[\\]{}\\(\\)'`~^@]")       // any of [, ], {, }, (, ), ', `, ~, ^, @
	pat.WriteString("|")                           // or
	pat.WriteString("\"(?:\\.|[^\\\"])*\"?")       // strings with escape characters and an optional " at the end
	pat.WriteString("|")                           // or
	pat.WriteString(";.*")                         // comments???
	pat.WriteString("|")                           // or
	pat.WriteString("[^\\s\\[\\]{}\\('\"`,;\\)]*") // atoms
	pat.WriteString(")")                           // End capture group
	// pat.WriteString("[\\s,]+")                     // whitespace or commas
	// pat.WriteString("|")                           // or
	// pat.WriteString("~@")                          // ~@
	// pat.WriteString("|")                           // or
	// pat.WriteString("[\\[\\]{}\\(\\)'`~^@]")       // any of [, ], {, }, (, ), ', `, ~, ^, @
	// pat.WriteString("|")                           // or
	// pat.WriteString("\"(\\.|[^\\\"])*\"?")         // strings with escape characters and an optional " at the end
	// pat.WriteString("|")                           // or
	// pat.WriteString(";.*")                         // comments???
	// pat.WriteString("|")                           // or
	// pat.WriteString("[^\\s\\[\\]{}\\('\"`,;\\)]*") // atoms
	return pat.String()
}
