package thrifter

import "testing"

func TestField_basic(t *testing.T) {
	parser := newParserOn(`1: required map< string , string > Test = { "abc": "def" } (api.test = "./test", api.a = 'asd',);`)
	n := NewField(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "Test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.ID, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Requiredness, "required"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FieldType.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.DefaultValue.Type, CONST_VALUE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "api.test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Value, `"./test"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestField_withoutDefaultValue(t *testing.T) {
	parser := newParserOn(`1: required map< string , string > Test      (api.test = "./test", api.a = 'asd',);`)
	n := NewField(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "Test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.ID, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Requiredness, "required"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FieldType.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "api.test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Value, `"./test"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestField_withoutDefaultValueAndOptions(t *testing.T) {
	parser := newParserOn(`1: required map< string , string > Test;`)
	n := NewField(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "Test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.ID, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Requiredness, "required"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FieldType.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestField_withoutDefaultValueAndOptionsAndRequiredness(t *testing.T) {
	parser := newParserOn(`1: map< string , string > Test;`)
	n := NewField(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "Test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.ID, 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FieldType.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "1"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestField_toString(t *testing.T) {
	src := `1: required map< string , string > Test = { "abc": "def" } (api.test = "./test", api.a = 'asd',);`
	parser := newParserOn(src)
	n := NewField(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
