package thrifter

import "testing"

func TestTypedef_basic(t *testing.T) {
	parser := newParserOn(`typedef binary b `)
	parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(nil, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
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
func TestTypedef_withOptions(t *testing.T) {
	parser := newParserOn(`typedef string ( unicode.encoding = "UTF-16" ) non_latin_string (foo="bar")`)
	parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(nil, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got, want := n.Ident, "non_latin_string"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Type.BaseType, "string"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Type.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "foo"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTypedef_map(t *testing.T) {
	parser := newParserOn(`typedef map<i64, map<i64, bool>> m`)
	parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(nil, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
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
	src := `typedef list< double ( cpp.fixed_point = "16" )> tiny_float_list`
	parser := newParserOn(src)
	startToken := parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(startToken, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTypedef_toStringWithOptions(t *testing.T) {
	src := `typedef string ( unicode.encoding = "UTF-16" ) non_latin_string (foo="bar")`
	parser := newParserOn(src)
	startToken := parser.next() // consume keyword token [typedef] first
	n := NewTypeDef(startToken, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
