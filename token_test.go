package thrifter

import "testing"

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
