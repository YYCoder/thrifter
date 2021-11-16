package thrifter

type TypeDef struct {
	NodeCommonField
	Type    *FieldType // except for identifier
	Ident   string
	Options []*Option
}

func NewTypeDef(start *Token, parent Node) *TypeDef {
	return &TypeDef{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *TypeDef) NodeType() string {
	return "TypeDef"
}

func (r *TypeDef) NodeValue() interface{} {
	return *r
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

	// parse options
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		return
	}
	p.next() // consume (
	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}
	return
}
