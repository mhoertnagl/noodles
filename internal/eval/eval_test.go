package eval_test

import (
	"github.com/mhoertnagl/splis2/internal/eval"
	"github.com/mhoertnagl/splis2/internal/print"
	"github.com/mhoertnagl/splis2/internal/read"
	"testing"
)

func TestEvalNumbers(t *testing.T) {
	test(t, "1", "1")
	test(t, "1.1", "1.1")
}

func TestEvalSum(t *testing.T) {
	test(t, "(+)", "0")
	test(t, "(+ 1)", "1")
	test(t, "(+ 1 1)", "2")
	test(t, "(+ 1 1 1)", "3")
}

func test(t *testing.T, i string, e string) {
	r := read.NewReader()
	r.Load(i)
	p := read.NewParser()
	n := p.Parse(r)
	w := print.NewPrinter()
	v := eval.NewEvaluator(eval.NewEnv(nil))
	m := v.Eval(n)
	a := w.Print(m)
	if a != e {
		t.Errorf("Expecting [%s] but got [%s]", e, a)
	}
}
