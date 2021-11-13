package thrifter

type ListType struct {
	NodeCommonField
	Elem    *FieldType
	CppType string
}

func NewListType(start *Token, parent Node) *ListType {
	return &ListType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *ListType) NodeValue() interface{} {
	return *r
}

func (r *ListType) NodeType() string {
	return "ListType"
}

func (r *ListType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ListType) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Type != tLESS {
		return p.unexpected(tok.Value, "<")
	}
	if err = r.parseElem(p); err != nil {
		return
	}
	return
}

func (r *ListType) parseElem(p *Parser) (err error) {
	r.Elem = NewFieldType(r)
	if err = r.Elem.parse(p); err != nil {
		return
	}
	p.peekNonWhitespace()
	greaterTok := p.next()
	if greaterTok.Type != tGREATER {
		err = p.unexpected(greaterTok.Value, ">")
		return
	}
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Type == tIDENT && tok.Value == "cpp_type" {
		p.peekNonWhitespace()
		strTok, err := p.nextString()
		if err != nil {
			return err
		}
		r.CppType = strTok.Value
		r.EndToken = strTok
	} else {
		// save it to buffer
		p.buf = tok
		r.EndToken = greaterTok
	}
	return
}
