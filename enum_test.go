package thrifter

import "testing"

func TestEnum_basic(t *testing.T) {
	parser := newParserOn(`enum a {
		A = 1
		B = 2;
		C
		D;
	}`)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.Elems), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[0].ID, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[0].Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[1].EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[2].EndToken.Value, "C"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[3].EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "enum"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEnum_elemsMap(t *testing.T) {
	parser := newParserOn(`enum a {
		A = 1
		B = 2;
		C
		D;
	}`)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.ElemsMap), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	for _, ele := range n.Elems {
		hash := GenTokenHash(ele.StartToken)
		if got, want := n.ElemsMap[hash], ele; got != want {
			t.Errorf("got [%v] want [%v]", got, want)
		}
	}
}

func TestEnum_withOptions(t *testing.T) {
	parser := newParserOn(`enum weekdays {
		SUNDAY ( weekend = "yes" ),
		MONDAY,
		TUESDAY,
		WEDNESDAY,
		THURSDAY,
		FRIDAY,
		SATURDAY ( weekend = "yes" )
	  } (foo.bar="baz")`)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.Elems), 7; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "weekdays"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "foo.bar"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "enum"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEnum_toString(t *testing.T) {
	src := `enum a {
		A = 1
		B
	C = 3
	}`
	parser := newParserOn(src)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestEnum_toStringWithComment(t *testing.T) {
	src := `enum a {
		A = 1 // comment
		B = 2 # comment comment
	C = 3 /* comment 1 */
	}`
	parser := newParserOn(src)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestEnum_toStringWithOptions(t *testing.T) {
	src := `enum weekdays {
		SUNDAY ( weekend = "yes" ),
		MONDAY,
		TUESDAY,
		WEDNESDAY,
		THURSDAY,
		FRIDAY,
		SATURDAY ( weekend = "yes" )
	  } (foo.bar="baz")`
	parser := newParserOn(src)
	startTok := parser.next()
	n := NewEnum(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
