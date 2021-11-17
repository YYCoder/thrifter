package thrifter

import "testing"

func TestService_basic(t *testing.T) {
	parser := newParserOn(`service A {
		double       testDouble(1: double thing) // test double
		oneway void testOneway(1:i32 secondsToSleep)
		void testException(1: string arg) throws(1: Xception err1),
		map<UserId, map<Numberz,Insanity>> testInsanity(1: Insanity argument);
	}`)
	start := parser.next()
	n := NewService(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 4; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Elems[0].FunctionType.BaseType, "double"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "service"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestService_elemsMap(t *testing.T) {
	parser := newParserOn(`service A {
		double       testDouble(1: double thing) // test double
		oneway void testOneway(1:i32 secondsToSleep)
		void testException(1: string arg) throws(1: Xception err1),
		map<UserId, map<Numberz,Insanity>> testInsanity(1: Insanity argument);
	}`)
	start := parser.next()
	n := NewService(start, nil)
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

func TestService_withOptions(t *testing.T) {
	parser := newParserOn(`service foo_service {
		void foo() ( foo = "bar" )
	  } (a.b="c")`)
	start := parser.next()
	n := NewService(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "foo_service"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "a.b"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "service"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestService_withExtends(t *testing.T) {
	parser := newParserOn(`service A extends B{
	}`)
	start := parser.next()
	n := NewService(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "A"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Elems), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Extends, "B"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "service"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, "}"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFunction_withOnewayAndVoid(t *testing.T) {
	parser := newParserOn(`oneway void testOneway(1:i32 secondsToSleep)`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "testOneway"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Void, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Oneway, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Args), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "oneway"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestFunction_withOptions(t *testing.T) {
	parser := newParserOn(`void test(1:i32 secondsToSleep) (api.get='/empty/msg',api.serializer='json')`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "test"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Void, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Oneway, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Args), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "api.get"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "void"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestFunction_multipleArgs(t *testing.T) {
	parser := newParserOn(`Xtruct testMulti(1: i8 arg0, 2: i32 arg1, 3: i64 arg2, 4: map<i16, string> arg3, 5: Numberz arg4, 6: UserId arg5)`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "testMulti"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Void, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Oneway, false; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FunctionType.Ident, "Xtruct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Args[3].FieldType.Type, FIELD_TYPE_MAP; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Args), 6; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "Xtruct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFunction_argsMap(t *testing.T) {
	parser := newParserOn(`Xtruct testMulti(1: i8 arg0, 2: i32 arg1, 3: i64 arg2, 4: map<i16, string> arg3, 5: Numberz arg4, 6: UserId arg5)`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.Args), 6; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	for _, ele := range n.Args {
		hash := GenTokenHash(ele.StartToken)
		if got, want := n.ArgsMap[hash], ele; got != want {
			t.Errorf("got [%v] want [%v]", got, want)
		}
	}
}

func TestFunction_throwsMap(t *testing.T) {
	parser := newParserOn(`void testException(1: string arg) throws(1: Xception err1),`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := len(n.Throws), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	for _, ele := range n.Throws {
		hash := GenTokenHash(ele.StartToken)
		if got, want := n.ThrowsMap[hash], ele; got != want {
			t.Errorf("got [%v] want [%v]", got, want)
		}
	}
}

func TestFunction_basic(t *testing.T) {
	parser := newParserOn(`double       testDouble(1: double thing) // test double`)
	n := NewFunction(nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.Ident, "testDouble"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.FunctionType.Type, FIELD_TYPE_BASE; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Args), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "double"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Value, ")"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestService_toString(t *testing.T) {
	src := `service A /* 123123 */
	{
		double       testDouble(1: double thing) // test double
		oneway void testOneway(1:i32 secondsToSleep) # 123123
		void testException(1: string arg) throws(1: Xception err1),
		map<UserId, map<Numberz,Insanity>> testInsanity(1: Insanity argument); /* hahaha */
		/* hahaha */
		set<UserId> withOptions(1: Insanity argument) (api.get='/empty/msg',api.serializer='json');
		# asdasd
		Xtruct testMulti(1: i8 arg0, 2: i32 arg1, 3: i64 arg2, 4: map<i16, string> arg3, 5: Numberz arg4, 6: UserId arg5)
	}`
	parser := newParserOn(src)
	start := parser.next()
	n := NewService(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestService_toStringWithOptions(t *testing.T) {
	src := `service foo_service {
		void foo() ( foo = "bar" )
	  } (a.b="c")`
	parser := newParserOn(src)
	start := parser.next()
	n := NewService(start, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
