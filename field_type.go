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
	Options  []*Option
}

func NewFieldType(parent Node) *FieldType {
	return &FieldType{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *FieldType) NodeType() string {
	return "FieldType"
}

func (r *FieldType) NodeValue() interface{} {
	return *r
}

func (r *FieldType) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *FieldType) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	identTok := p.nextIdent(true)
	r.StartToken = identTok
	if isBaseTypeToken(identTok.Raw) {
		r.Type = FIELD_TYPE_BASE
		r.BaseType = identTok.Raw
		r.EndToken = identTok
	} else if identTok.Raw == "map" {
		r.Type = FIELD_TYPE_MAP
		r.Map = NewMapType(identTok, r)
		if err = r.Map.parse(p); err != nil {
			return
		}
		r.EndToken = r.Map.EndToken
	} else if identTok.Raw == "set" {
		r.Type = FIELD_TYPE_SET
		r.Set = NewSetType(identTok, r)
		if err = r.Set.parse(p); err != nil {
			return
		}
		r.EndToken = r.Set.EndToken
	} else if identTok.Raw == "list" {
		r.Type = FIELD_TYPE_LIST
		r.List = NewListType(identTok, r)
		if err = r.List.parse(p); err != nil {
			return
		}
		r.EndToken = r.List.EndToken
	} else {
		r.Type = FIELD_TYPE_IDENT
		r.Ident = identTok.Raw
		r.EndToken = identTok
	}

	// parse options
	// list type may save token to buffer, since it need to scan next cpp_type token
	if p.buf != nil {
		if p.buf.Type != T_LEFTPAREN {
			return
		}
		p.buf = nil
		r.Options, r.EndToken, err = parseOptions(p, r)
		if err != nil {
			return err
		}
	} else {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) != T_LEFTPAREN {
			return
		}
		p.next() // consume (
		r.Options, r.EndToken, err = parseOptions(p, r)
		if err != nil {
			return err
		}
	}

	return
}
