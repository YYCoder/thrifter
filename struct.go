package thrifter

const (
	STRUCT = iota + 1
	UNION
	EXCEPTION
)

type Struct struct {
	NodeCommonField
	Type    int
	Ident   string
	Elems   []*Field
	Options []*Option
}

func NewStruct(start *Token, parent Node) *Struct {
	t := 0
	switch start.Value {
	case "struct":
		t = STRUCT
	case "union":
		t = UNION
	case "exception":
		t = EXCEPTION
	}
	return &Struct{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
		Type: t,
	}
}

func (r *Struct) NodeType() string {
	switch r.Type {
	case STRUCT:
		return "Struct"
	case UNION:
		return "Union"
	case EXCEPTION:
		return "Exception"
	}
	return ""
}

func (r *Struct) NodeValue() interface{} {
	return *r
}

func (r *Struct) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Struct) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, _, _ := p.nextIdent(false)
	r.Ident = fullLit
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTCURLY {
		return p.unexpected(string(ru), "{")
	}
	p.next() // consume {
	var rightParenTok *Token
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tRIGHTCURLY {
			rightParenTok = p.next()
			break
		}
		elem := NewField(r)
		if err = elem.parse(p); err != nil {
			return err
		}
		r.Elems = append(r.Elems, elem)
	}

	// parse options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTPAREN {
		r.EndToken = rightParenTok
		return
	}
	p.next() // consume (
	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}
	return
}
