package thrifter

// the definition of what kind of field type
const (
	FIELD_TYPE_IDENT = iota + 1
	FIELD_TYPE_BASE
	FIELD_TYPE_MAP
	FIELD_TYPE_LIST
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

func NewFieldType(parent Node) *FieldType {
	return &FieldType{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *FieldType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *FieldType) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, startTok, endTok := p.nextIdent(true)
	r.StartToken = startTok
	if isBaseTypeToken(fullLit) {
		r.Type = FIELD_TYPE_BASE
		r.BaseType = fullLit
		r.EndToken = endTok
	} else if fullLit == "map" {
		r.Type = FIELD_TYPE_MAP
		r.Map = NewMapType(startTok, r)
		if err = r.Map.parse(p); err != nil {
			return
		}
		r.EndToken = r.Map.EndToken
	} else if fullLit == "set" {
		r.Type = FIELD_TYPE_SET
		r.Set = NewSetType(startTok, r)
		if err = r.Set.parse(p); err != nil {
			return
		}
		r.EndToken = r.Set.EndToken
	} else if fullLit == "list" {
		r.Type = FIELD_TYPE_LIST
		r.List = NewListType(startTok, r)
		if err = r.List.parse(p); err != nil {
			return
		}
		r.EndToken = r.List.EndToken
	} else {
		r.Type = FIELD_TYPE_IDENT
		r.Ident = fullLit
		r.EndToken = endTok
	}
	return
}
