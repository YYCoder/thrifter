package thrifter

type MapType struct {
	NodeCommonField
	KeyType string // base type tokens
	Value   string
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
	return
}
