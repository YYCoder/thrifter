package thrifter

type TypeDef struct {
	NodeCommonField
	Type  *FieldType // except for identifier
	Ident string
}

func NewTypeDef(start *Token, parent Node) *TypeDef {
	return &TypeDef{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *TypeDef) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *TypeDef) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	r.Type = NewFieldType(r)
	if err = r.Type.parse(p); err != nil {
		return
	}
	if r.Type.Type == FIELD_TYPE_IDENT {
		return p.unexpected(r.Type.Ident, "base type or map or list or set")
	}

	fullLit, _, endTok := p.nextIdent(true)
	r.Ident = fullLit
	r.EndToken = endTok
	return
}
