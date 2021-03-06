package thrifter

type Namespace struct {
	NodeCommonField
	Name    string
	Value   string
	Options []*Option
}

func NewNamespace(start *Token, parent Node) *Namespace {
	return &Namespace{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Namespace) NodeValue() interface{} {
	return *r
}

func (r *Namespace) NodeType() string {
	return "Namespace"
}

func (r *Namespace) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Namespace) parse(p *Parser) (err error) {
	identTok := p.nextIdent(true)
	r.Name = identTok.Raw
	identTok = p.nextIdent(true)
	r.Value = identTok.Raw
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		r.EndToken = identTok
		return
	}
	p.next() // consume (

	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}
	return
}
