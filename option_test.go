package thrifter

import "testing"

func TestOption_singleQuoteValue(t *testing.T) {
	parser := newParserOn(`a = '123'`)
	n := NewOption(nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "'123'"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOption_quoteValue(t *testing.T) {
	parser := newParserOn(`a = "123"`)
	n := NewOption(nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "\"123\""; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOption_withWhiteSpaces(t *testing.T) {
	parser := newParserOn(`a        =	
		"123"`)
	n := NewOption(nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "\"123\""; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOption_toString(t *testing.T) {
	src := `a        =	
	"123"`
	parser := newParserOn(src)
	n := NewOption(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
