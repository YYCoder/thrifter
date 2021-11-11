package thrifter

type Field struct {
	NodeCommonField
	ID           int
	Requiredness bool
	FieldType    *FieldType
	Ident        string
	Value        *ConstValue
}

func NewField(start *Token, parent Node) *Field {
	return &Field{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Field) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Field) parse(p *Parser) (err error) {
	return
}
