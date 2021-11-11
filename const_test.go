package thrifter

import "testing"

func TestConst_basic(t *testing.T) {
	parser := newParserOn(`i32 test = 123;`)
	constNode := NewConst(nil, nil)
	constNode.parse(parser)

	if got, want := constNode.Type.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.BaseType, `i32`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Ident, "test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Type, CONST_VALUE_INT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Value, "123"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.EndToken.Value, ";"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConst_map(t *testing.T) {
	parser := newParserOn(`map<i64, string> test = {
		123: "hahah",
	}`)
	constNode := NewConst(nil, nil)
	err := constNode.parse(parser)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := constNode.Type.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Key.BaseType, "i64"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Value.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Value.BaseType, "string"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Ident, "test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Type, CONST_VALUE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constNode.Value.Map.MapKeyList), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Map.MapKeyList[0].Value, "123"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Map.MapValueList[0].Value, `"hahah"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.EndToken.Value, `}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConst_mapWithIdent(t *testing.T) {
	parser := newParserOn(`map<i64, abc.def> test = {
		123: "hahah"
	}`)
	constNode := NewConst(nil, nil)
	err := constNode.parse(parser)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := constNode.Type.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Key.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Key.BaseType, "i64"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Value.Type, FIELD_TYPE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.Map.Value.BaseType, "abc.def"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Ident, "test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Type, CONST_VALUE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constNode.Value.Map.MapKeyList), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Map.MapKeyList[0].Value, "123"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Map.MapValueList[0].Value, `"hahah"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.EndToken.Value, `}`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConst_list(t *testing.T) {
	parser := newParserOn(`list<i64> test = [
		123, 234, 345, 456
	]`)
	constNode := NewConst(nil, nil)
	err := constNode.parse(parser)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := constNode.Type.Type, FIELD_TYPE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.List.Elem.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Type.List.Elem.BaseType, "i64"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Ident, "test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.Type, CONST_VALUE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constNode.Value.List.Elems), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.List.Elems[0].Type, CONST_VALUE_INT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.Value.List.Elems[0].Value, `123`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constNode.EndToken.Value, `]`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConstValue_positiveNumber(t *testing.T) {
	parser := newParserOn(`0.123`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_FLOAT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Value, `0.123`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.StartToken, constVal.EndToken; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConstValue_negativeNumber(t *testing.T) {
	parser := newParserOn(`-0.123`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_FLOAT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Value, `-0.123`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.StartToken, constVal.EndToken; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConstValue_literal(t *testing.T) {
	parser := newParserOn(`"123"`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_LITERAL; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Value, `"123"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.StartToken, constVal.EndToken; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConstValue_ident(t *testing.T) {
	parser := newParserOn(`a1bc`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Value, `a1bc`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.StartToken, constVal.EndToken; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestConstValue_identWithDot(t *testing.T) {
	parser := newParserOn(`a.b.c`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Value, `a.b.c`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	// TODO: concat ident token
	// if got, want := constVal.StartToken, constVal.EndToken; got != want {
	// 	t.Errorf("got [%v] want [%v]", got, want)
	// }
}

func TestConstValue_listBasic(t *testing.T) {
	parser := newParserOn(`[a,b,c]`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.List.Elems), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[0].Value, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[0].Type, CONST_VALUE_IDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.EndToken.Value, `]`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
}

func TestConstValue_listLiteral(t *testing.T) {
	parser := newParserOn(`[
		"a","b","c"
	]`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.List.Elems), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[0].Value, `"a"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[0].Type, CONST_VALUE_LITERAL; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.EndToken.Value, `]`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
}

func TestConstValue_listFloat(t *testing.T) {
	parser := newParserOn(`[0.5, -12.123, 3.23]`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_LIST; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.List.Elems), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[1].Value, `-12.123`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.List.Elems[0].Type, CONST_VALUE_FLOAT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.EndToken.Value, `]`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
}

func TestConstValue_mapBasic(t *testing.T) {
	parser := newParserOn(`{"a": b.c.d}`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.Map.MapKeyList), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.Map.MapValueList), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[0].Value, `"a"`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[0].Type, CONST_VALUE_LITERAL; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[0].Value, "b.c.d"; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[0].Type, CONST_VALUE_IDENT; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
}

func TestConstValue_mapMultiple(t *testing.T) {
	parser := newParserOn(`{
		"a": b.c.d, c.d : 123.123; -123: "a"
	}`)
	constVal := NewConstValue(nil)
	constVal.parse(parser)

	if got, want := constVal.Type, CONST_VALUE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.Map.MapKeyList), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(constVal.Map.MapValueList), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[1].Value, `c.d`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[1].Type, CONST_VALUE_IDENT; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[1].Type, CONST_VALUE_FLOAT; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[1].Value, "123.123"; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[2].Type, CONST_VALUE_INT; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapKeyList[2].Value, "-123"; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[2].Type, CONST_VALUE_LITERAL; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.Map.MapValueList[2].Value, `"a"`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
	if got, want := constVal.EndToken.Value, `}`; got != want {
		t.Errorf("got [%+v] want [%+v]", got, want)
	}
}

func TestConst_toString(t *testing.T) {
	src := `const 		i64  	 a =   123;`
	parser := newParserOn(src)
	startToken := parser.next() // consume keyword token [const] first
	n := NewConst(startToken, nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
