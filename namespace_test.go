package thrifter

import "testing"

func TestNamespace_basic(t *testing.T) {
	parser := newParserOn(`namespace a b `)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "b"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNamespace_withStar(t *testing.T) {
	parser := newParserOn(`namespace * b `)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "*"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "b"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNamespace_withDotSeparator(t *testing.T) {
	parser := newParserOn(`namespace * .b.a.c `)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "*"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, ".b.a.c"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNamespace_withSingleOption(t *testing.T) {
	parser := newParserOn(`namespace * b.a.c (a = 'a/b/c')`)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "*"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "b.a.c"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Value, "'a/b/c'"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNamespace_withMultipleOptions(t *testing.T) {
	parser := newParserOn(`namespace * b.a.c (a = 'a/b/c', b="4565a")`)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	if err != nil {
		t.Error(err)
	}
	if got, want := n.Name, "*"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Value, "b.a.c"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(n.Options), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Name, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[0].Value, "'a/b/c'"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[1].Name, "b"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.Options[1].Value, "\"4565a\""; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNamespace_withExtraComma(t *testing.T) {
	parser := newParserOn(`namespace * b.a.c (a = 'a/b/c', b="4565a", )`)
	parser.next() // consume keyword token [namespace] first
	n := NewNamespace(nil, nil)
	err := n.parse(parser)
	// expect an error
	if err == nil {
		t.Error(err)
	}
}

func TestNamespace_toString(t *testing.T) {
	src := `namespace * b.a.c (a = 'a/b/c', b="4565a")`
	parser := newParserOn(src)
	startToken := parser.next() // consume keyword token [namespace] first
	n := NewNamespace(startToken, nil)
	_ = n.parse(parser)
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
