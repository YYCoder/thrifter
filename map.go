package thrifter

type MapType struct {
	NodeCommonField
	// since directly use map structure its hard to index and will lead to loss of order, we use slice to represent map type, use slice index to mapping
	Key     *FieldType
	Value   *FieldType
	CppType string
}

func NewMapType(start *Token, parent Node) *MapType {
	return &MapType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *MapType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *MapType) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Type == tLESS {
		if err = r.parseKeyAndValue(p); err != nil {
			return
		}
	} else if tok.Type == tIDENT && tok.Value == "cpp_type" {
		p.peekNonWhitespace()
		strTok, err := p.nextString()
		if err != nil {
			return err
		}
		r.CppType = strTok.Value
		p.peekNonWhitespace()
		tok := p.next()
		if tok.Type != tLESS {
			return p.unexpected(tok.Value, "<")
		}
		if err = r.parseKeyAndValue(p); err != nil {
			return err
		}
	} else {
		err = p.unexpected(tok.Value, "< or cpp_type")
	}
	return
}

func (r *MapType) parseKeyAndValue(p *Parser) (err error) {
	r.Key = NewFieldType(r)
	if err = r.Key.parse(p); err != nil {
		return
	}
	ru := p.peekNonWhitespace()
	commaTok := toToken(string(ru))
	if commaTok != tCOMMA {
		err = p.unexpected(string(ru), ",")
		return
	}
	// consume comma token
	p.next()
	r.Value = NewFieldType(r)
	if err = r.Value.parse(p); err != nil {
		return
	}
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Type != tGREATER {
		err = p.unexpected(tok.Value, ">")
		return
	}
	r.EndToken = tok
	return
}
