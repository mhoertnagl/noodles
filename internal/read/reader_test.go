package read_test

import (
	"testing"

	"github.com/mhoertnagl/splis2/internal/read"
)

func TestIgnoreWhitespace(t *testing.T) {
	testr(t, " \t\r\n", "")
}

func TestSpecials(t *testing.T) {
	testr(t, "~@", "~@", "")
	testr(t, "(", "(", "")
	testr(t, ")", ")", "")
	testr(t, "[", "[", "")
	testr(t, "]", "]", "")
	testr(t, "{", "{", "")
	testr(t, "}", "}", "")
	testr(t, "'", "'", "")
	testr(t, "`", "`", "")
	testr(t, "~", "~", "")
	testr(t, "^", "^", "")
	testr(t, "@", "@", "")
}

func TestStrings(t *testing.T) {
	testr(t, `""`, `""`, "")
	testr(t, `"xxx"`, `"xxx"`, "")
	testr(t, `"\n"`, `"\n"`, "")
	testr(t, `"([{'^~@}])"`, `"([{'^~@}])"`, "")
	testr(t, `"xxx`, `"xxx`, "")
	testr(t, `"\n"`, `"\n"`, "")
	testr(t, `"\\"`, `"\\"`, "")
	testr(t, `"\""`, `"\""`, "")
}

func TestSymbols(t *testing.T) {
	testr(t, "123.45", "123.45", "")
	testr(t, "+", "+", "")
}

func TestComments(t *testing.T) {
	testr(t, "; This is a comment.", "")
}

func TestComments2(t *testing.T) {
	testr(t, "  ;; This is a comment.", "")
}

func TestComments3(t *testing.T) {
	src := `
    ;; Returns the negation of x.
    1
  `
	testr(t, src, "1", "")
}

func TestList(t *testing.T) {
	testr(t, "  (+ 1   2  )   ", "(", "+", "1", "2", ")", "")
}

func TestProg(t *testing.T) {
	src := `
  ;; Returns the negation of x.
  (def! not (fn* [x]
    (if x false true)))
  `
	testr(t, src,
		"(", "def!", "not", "(", "fn*", "[", "x", "]",
		"(", "if", "x", "false", "true", ")", ")", ")", "")
}

func TestPlusOne(t *testing.T) {
	testr(t, "  (+1    0 )   ", "(", "+1", "0", ")", "")
}

func testr(t *testing.T, i string, es ...string) {
	r := read.NewReader()
	r.Load(i)
	for idx, e := range es {
		tok := r.Next()
		if tok != e {
			t.Errorf("Expecting [%s] at pos [%d] but got [%s]", e, idx+1, tok)
		}
	}
}
