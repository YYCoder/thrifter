package thrifter

import "testing"

func TestStruct_basic(t *testing.T) {
	parser := newParserOn(`struct A {
		1: i32 Test;
		2: bool Haha;
	}`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Type, STRUCT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "struct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestStruct_elemsMap(t *testing.T) {
	parser := newParserOn(`struct Test{
		1: required i32 Foo = 123 (api.test = "./test", api.a = 'asd',);
		2: optional map<i64,string> Bar = { 123: "bar" }
		3: set<i64> Set;
		4: list<test.test.test> List = [aaa, bbb, ccc]; // for compatibility with framework
	}`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.Elems), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	for _, ele := range n.Elems {
		hash := GenTokenHash(ele.StartToken)
		if got, want := n.ElemsMap[hash], ele; got != want {
			t.Errorf("got [%v] want [%v]", got, want)
		}
	}
}

func TestStruct_withOptions(t *testing.T) {
	parser := newParserOn(`struct A {
		1: i32 Test;
		2: bool Haha;
	} (a.b.c = "123")`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Type, STRUCT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "a.b.c"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "struct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestUnion_basic(t *testing.T) {
	parser := newParserOn(`union A {
		1: i32 Test;
		2: bool Haha;
	}`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Type, UNION; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "union"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestException_basic(t *testing.T) {
	parser := newParserOn(`exception A {
		1: i32 Test;
		2: bool Haha;
	}`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Type, EXCEPTION; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "exception"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestStruct_complex(t *testing.T) {
	parser := newParserOn(`struct Test{
		1: required i32 Foo = 123 (api.test = "./test", api.a = 'asd',);
		2: optional map<i64,string> Bar = { 123: "bar" }
		3: set<i64> Set;
		4: list<test.test.test> List = [aaa, bbb, ccc]; // for compatibility with framework
	}`)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "Test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems[3].DefaultValue.List.Elems), 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[3].FieldType.List.Elem.Ident, "test.test.test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestStruct_toString(t *testing.T) {
	src := `struct Test{
		1: required i32 Foo = 123 (api.test = "./test", api.a = 'asd',);
		2: optional map<i64,string> Bar = { 123: "bar" }
		3: set<i64> Set;
		4: list<test.test.test> List = [aaa, bbb, ccc]; // for compatibility with framework
	}`
	parser := newParserOn(src)
	start := parser.next()
	n := NewStruct(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
