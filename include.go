package thrifter

type Include struct {
	NodeCommonField
	FilePath string
}

func NewInclude(start *Token, parent Node) *Include {
	return &Include{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Include) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Include) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	tok, err := p.nextString()
	if err != nil {
		return err
	}
	r.EndToken = tok
	r.FilePath = tok.Value

	return
}
