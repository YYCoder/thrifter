package thrifter

// the definition of what kind of field type
const (
	FIELD_TYPE_IDENT = iota + 1
	FIELD_TYPE_BASE
	FIELD_TYPE_MAP
	FIELD_TYPE_List
	FIELD_TYPE_SET
)

type FieldType struct {
	NodeCommonField
	Type     int
	Ident    string
	BaseType string
	Map      *MapType
	List     *ListType
	Set      *SetType
}

func NewFieldType(start *Token, parent Node) *FieldType {
	return &FieldType{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *FieldType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *FieldType) parse(p *Parser) (err error) {
	return
}
