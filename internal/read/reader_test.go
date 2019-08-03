package read_test

import (
	//"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"testing"
)

func TestIgnoreWhitespace(t *testing.T) {
	test(t, " \t\r\n", "")
}

func TestSpecials(t *testing.T) {
	test(t, "~@", "~@", "")
	test(t, "(", "(", "")
	test(t, ")", ")", "")
	test(t, "[", "[", "")
	test(t, "]", "]", "")
	test(t, "{", "{", "")
	test(t, "}", "}", "")
	test(t, "'", "'", "")
	test(t, "`", "`", "")
	test(t, "~", "~", "")
	test(t, "^", "^", "")
	test(t, "@", "@", "")
}

func TestStrings(t *testing.T) {
	test(t, `""`, `""`, "")
	test(t, `"xxx"`, `"xxx"`, "")
	test(t, `"\n"`, `"\n"`, "")
	test(t, `"([{'^~@}])"`, `"([{'^~@}])"`, "")
	test(t, `"xxx`, `"xxx`, "")
}

func TestSymbols(t *testing.T) {
	test(t, "123.45", "123.45", "")
	test(t, "+", "+", "")
}

func TestComments(t *testing.T) {
	test(t, "; This is a comment", "")
}

func TestList(t *testing.T) {
	test(t, "  (+ 1   2  )   ", "(", "+", "1", "2", ")", "")
}

func test(t *testing.T, i string, es ...string) {
	r := read.NewReader()
	r.Load(i)
	for idx, e := range es {
		tok := r.Next()
		if tok != e {
			t.Errorf("Expecting [%s] at pos [%d] but got [%s]", e, idx+1, tok)
		}
	}
}
