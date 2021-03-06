package thrifter

type SetType struct {
	NodeCommonField
	Elem    *FieldType
	CppType string
}

func NewSetType(start *Token, parent Node) *SetType {
	return &SetType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *SetType) NodeType() string {
	return "SetType"
}

func (r *SetType) NodeValue() interface{} {
	return *r
}

func (r *SetType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *SetType) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Type == T_LESS {
		if err = r.parseElem(p); err != nil {
			return
		}
	} else if tok.Type == T_IDENT && tok.Value == "cpp_type" {
		p.peekNonWhitespace()
		strTok, err := p.nextString()
		if err != nil {
			return err
		}
		r.CppType = strTok.Value
		p.peekNonWhitespace()
		tok := p.next()
		if tok.Type != T_LESS {
			return p.unexpected(tok.Value, "<")
		}
		if err = r.parseElem(p); err != nil {
			return err
		}
	} else {
		err = p.unexpected(tok.Value, "< or cpp_type")
	}
	return
}

func (r *SetType) parseElem(p *Parser) (err error) {
	r.Elem = NewFieldType(r)
	if err = r.Elem.parse(p); err != nil {
		return
	}
	tok := p.next()
	if tok.Type != T_GREATER {
		err = p.unexpected(tok.Value, ">")
		return
	}
	r.EndToken = tok
	return
}
