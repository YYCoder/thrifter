package thrifter

import (
	"testing"
	"text/scanner"
)

func TestIsNumber_int(t *testing.T) {
	isFloat, isInt := IsNumber(`123`)

	if got, want := isInt, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := isFloat, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIsNumber_float0(t *testing.T) {
	isFloat, isInt := IsNumber(`123.123`)

	if got, want := isInt, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := isFloat, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIsNumber_float1(t *testing.T) {
	isFloat, isInt := IsNumber(`0.123`)

	if got, want := isInt, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := isFloat, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIsNumber_invalidFloat(t *testing.T) {
	isFloat, isInt := IsNumber(`.123`)

	if got, want := isInt, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := isFloat, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestGenTokenHash(t *testing.T) {
	hash1 := GenTokenHash(&Token{
		Type:  T_COMMENT,
		Raw:   "// asdasd",
		Value: " asdasd",
		Pos: scanner.Position{
			Line:   1,
			Column: 1,
		},
	})
	hash2 := GenTokenHash(&Token{
		Type:  T_COMMENT,
		Raw:   "// asdasd",
		Value: " asdasd",
		Pos: scanner.Position{
			Line:   1,
			Column: 1,
		},
	})
	hash3 := GenTokenHash(&Token{
		Type:  T_COMMENT,
		Raw:   "// asdasd",
		Value: " asdasd",
		Pos: scanner.Position{
			Line:   3,
			Column: 1,
		},
	})

	if got, want := hash1, hash2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := hash1, hash3; got == want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
