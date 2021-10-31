package thrifter

import (
	"strings"
	"testing"
)

func newParserOn(def string) *Parser {
	p := NewParser(strings.NewReader(def))
	return p
}

func TestNextIdent_singleIdent(t *testing.T) {
	parser := newParserOn(` ab2 `)
	lit, start, end := parser.nextIdent(false)

	if got, want := start.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, "ab2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_mulitpleWhitespace(t *testing.T) {
	parser := newParserOn(`
	 	ab2 `)
	lit, start, end := parser.nextIdent(false)

	if got, want := start.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, "ab2"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_trailingDot(t *testing.T) {
	parser := newParserOn(` abc.def. `)
	lit, start, end := parser.nextIdent(false)

	if got, want := start.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tDOT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, "abc.def."; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_leadingDot(t *testing.T) {
	parser := newParserOn(` .abc.def `)
	lit, start, end := parser.nextIdent(false)

	if got, want := start.Type, tDOT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, ".abc.def"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_leadingAndTrailingDot(t *testing.T) {
	parser := newParserOn(` .abc.def. `)
	lit, start, end := parser.nextIdent(false)

	if got, want := start.Type, tDOT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tDOT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, ".abc.def."; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_keyword(t *testing.T) {
	parser := newParserOn(`enum.def.struct `)
	lit, start, end := parser.nextIdent(true)

	if got, want := start.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, "enum.def.struct"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextIdent_star(t *testing.T) {
	parser := newParserOn(` * `)
	lit, start, end := parser.nextIdent(true)

	if got, want := start.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := end.Type, tIDENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := lit, "*"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextString_quoteString(t *testing.T) {
	// parser := newParserOn(`"http://thrift.apache.org/ns/ThriftTest"`)
	parser := newParserOn("\"http://thrift.apache.org/ns/ThriftTest\"")
	tok, _ := parser.nextString()

	if got, want := tok.Type, tSTRING; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, "\"http://thrift.apache.org/ns/ThriftTest\""; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, "http://thrift.apache.org/ns/ThriftTest"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextString_singleQuoteString(t *testing.T) {
	parser := newParserOn(`'http://thrift.apache.org/ns/ThriftTest'`)
	tok, _ := parser.nextString()

	if got, want := tok.Type, tSTRING; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, "'http://thrift.apache.org/ns/ThriftTest'"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, "http://thrift.apache.org/ns/ThriftTest"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
func TestNextString_errString(t *testing.T) {
	parser := newParserOn(`'http://thrift.apache.org
	/ns/ThriftTest'`)
	_, err := parser.nextString()

	if !strings.Contains(err.Error(), "EOF or LineBreak") {
		t.Errorf("got [%v] want [%v]", err.Error(), "EOF or LineBreak")
	}
}

func TestNextComment_singleLineBasic(t *testing.T) {
	parser := newParserOn(`/123123 asasd
	`)
	tok, _ := parser.nextComment(SINGLE_LINE_COMMENT)

	if got, want := tok.Type, tCOMMENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, "123123 asasd"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, "//123123 asasd"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextComment_bashBasic(t *testing.T) {
	parser := newParserOn(`123123 asasd
	`)
	tok, _ := parser.nextComment(BASH_LIKE_COMMENT)

	if got, want := tok.Type, tCOMMENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, "123123 asasd"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, "#123123 asasd"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestNextComment_multiLineBasic(t *testing.T) {
	parser := newParserOn(`*123123 asasd
	*/`)
	tok, _ := parser.nextComment(MULTI_LINE_COMMENT)

	if got, want := tok.Type, tCOMMENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, `123123 asasd
	`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, `/*123123 asasd
	*/`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestIsComment_multiLineBasic(t *testing.T) {
	parser := newParserOn(`*123123 asasd
	*/`)
	_, ct := parser.isComment('/')
	tok, _ := parser.nextComment(ct)

	if got, want := tok.Type, tCOMMENT; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Value, `123123 asasd
	`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := tok.Raw, `/*123123 asasd
	*/`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
