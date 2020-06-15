package util_test

import (
	"testing"

	"github.com/mhoertnagl/noodles/internal/util"
)

func TestEmpty(t *testing.T) {
	a := util.IndexOf([]string{}, "&")
	assertEqualInt(t, a, -1)
}

func TestSngleSuccess(t *testing.T) {
	a := util.IndexOf([]string{"&"}, "&")
	assertEqualInt(t, a, 0)
}

func TestSngleFail(t *testing.T) {
	a := util.IndexOf([]string{"x"}, "&")
	assertEqualInt(t, a, -1)
}

func TestMultipleSuccess1(t *testing.T) {
	a := util.IndexOf([]string{"&", "y", "z"}, "&")
	assertEqualInt(t, a, 0)
}

func TestMultipleSuccess2(t *testing.T) {
	a := util.IndexOf([]string{"x", "&", "z"}, "&")
	assertEqualInt(t, a, 1)
}

func TestMultipleSuccess3(t *testing.T) {
	a := util.IndexOf([]string{"x", "y", "&"}, "&")
	assertEqualInt(t, a, 2)
}

func TestMultipleFail(t *testing.T) {
	a := util.IndexOf([]string{"x", "y", "z"}, "&")
	assertEqualInt(t, a, -1)
}

func assertEqualInt(t *testing.T, act int, exp int) {
	t.Helper()
	if act != exp {
		t.Errorf("Expecting [%d] but got [%d]", exp, act)
	}
}
