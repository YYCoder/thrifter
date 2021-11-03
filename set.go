package thrifter

type SetType struct {
	NodeCommonField
	Name  string
	Value string
}

func NewSetType(start *Token, parent Node) *SetType {
	return &SetType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *SetType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *SetType) parse(p *Parser) (err error) {
	return
}
