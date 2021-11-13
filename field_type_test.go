package thrifter

import "testing"

func TestFieldType_base(t *testing.T) {
	parser := newParserOn(`i32`)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.BaseType, `i32`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `i32`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken, node.EndToken; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_set(t *testing.T) {
	parser := newParserOn(`   set<abc.def> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_SET; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `set`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_setWithNestedMap(t *testing.T) {
	parser := newParserOn(`   set<map<i64, abc.def>> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_SET; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Map.Key.BaseType, "i64"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Map.Value.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `set`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_setWithCppType(t *testing.T) {
	parser := newParserOn(`   set cpp_type 'haha' <abc.def> `)
	node := NewFieldType(nil)
	if err := node.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := node.Type, FIELD_TYPE_SET; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.CppType, "haha"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `set`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_map(t *testing.T) {
	parser := newParserOn(`   map<string, abc.def> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.BaseType, `string`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `map`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_mapWithNested(t *testing.T) {
	parser := newParserOn(`   map<string, map<i64, bool>> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.BaseType, `string`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Map.Key.BaseType, `i64`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Map.Value.BaseType, `bool`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `map`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_mapWithCppType(t *testing.T) {
	parser := newParserOn(`   map cpp_type 'haha' <string, abc.def> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.CppType, "haha"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Key.BaseType, `string`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Map.Value.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `map`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_list(t *testing.T) {
	parser := newParserOn(`   list<string> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.BaseType, `string`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `list`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_listWithNestedMap(t *testing.T) {
	parser := newParserOn(`   list<map<i64, bool>> `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Map.Key.BaseType, `i64`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Map.Value.BaseType, `bool`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `list`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_listWithCppType(t *testing.T) {
	parser := newParserOn(`   list<string> cpp_type 'haha' `)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.CppType, "haha"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.BaseType, `string`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `list`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `haha`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_toStringMap(t *testing.T) {
	src := `map<   i64, 			bool>`
	parser := newParserOn(src)
	n := NewFieldType(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestFieldType_toStringList(t *testing.T) {
	src := `list<   i64> cpp_type 'haha'`
	parser := newParserOn(src)
	n := NewFieldType(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_toStringSet(t *testing.T) {
	src := `set cpp_type '123' <   i64>`
	parser := newParserOn(src)
	n := NewFieldType(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_toStringWithOptions(t *testing.T) {
	src := `set cpp_type '123' <   i64> (abc.def = 'asdasds')`
	parser := newParserOn(src)
	n := NewFieldType(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestFieldType_toStringWithInsaneNesting(t *testing.T) {
	src := `list<map<set<i32> (python.immutable = ""), map<i32,set<list<map<Insanity,string>(python.immutable = "")> (python.immutable = "")>>>>`
	parser := newParserOn(src)
	n := NewFieldType(nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_setWithOptions(t *testing.T) {
	parser := newParserOn(`   set<abc.def> (test.test = "123")`)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_SET; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Set.Elem.Ident, `abc.def`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(node.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Options[0].Name, "test.test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `set`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `)`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_listWithOptions(t *testing.T) {
	parser := newParserOn(`list<map<set<i32> (python.immutable = ""), map<i32,set<list<map<Insanity,string>(python.immutable = "")> (python.immutable = "")>>>>`)
	node := NewFieldType(nil)
	err := node.parse(parser)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := node.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Map.Key.Type, FIELD_TYPE_SET; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Map.Value.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(node.List.Elem.Map.Key.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.List.Elem.Map.Key.Options[0].Name, "python.immutable"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `list`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `>`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFieldType_basicWithOptions(t *testing.T) {
	parser := newParserOn(` bool (test.test = "123")`)
	node := NewFieldType(nil)
	node.parse(parser)

	if got, want := node.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.BaseType, `bool`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(node.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.Options[0].Name, "test.test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.StartToken.Value, `bool`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := node.EndToken.Value, `)`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
