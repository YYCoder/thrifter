package thrifter

type ListType struct {
	NodeCommonField
	Name  string
	Value string
}

func NewListType(start *Token, parent Node) *ListType {
	return &ListType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *ListType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ListType) parse(p *Parser) (err error) {
	return
}
