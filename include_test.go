package thrifter

import "testing"

func TestInclude_basic(t *testing.T) {
	parser := newParserOn(`include "./../test.thrift"`)
	startTok := parser.next()
	n := NewInclude(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got, want := n.FilePath, "./../test.thrift"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.StartToken.Value, "include"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := n.EndToken.Raw, `"./../test.thrift"`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInclude_toString(t *testing.T) {
	src := `include "./test.thrift"`
	parser := newParserOn(src)
	startTok := parser.next()
	n := NewInclude(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInclude_toStringCpp(t *testing.T) {
	src := `cpp_include "./test.thrift"`
	parser := newParserOn(src)
	startTok := parser.next()
	n := NewInclude(startTok, nil)
	if err := n.parse(parser); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	res := n.String()

	if got, want := res, src; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
