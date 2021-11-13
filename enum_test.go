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
