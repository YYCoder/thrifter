package thrifter

import "testing"

func TestTypedef_basic(t *testing.T) {
	parser := newParserOn(`typedef binary b `)
	parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Ident, "b"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.BaseType, "binary"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTypedef_map(t *testing.T) {
	parser := newParserOn(`typedef map<i64, map<i64, bool>> m`)
	parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Ident, "m"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.Map.Value.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTypedef_toString(t *testing.T) {
	src := `typedef map<i64, map<i64, bool>> m`
	parser := newParserOn(src)
	startToken := parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(startToken, nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
