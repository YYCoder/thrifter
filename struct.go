package thrifter

const (
	STRUCT = iota + 1
	UNION
	EXCEPTION
)

type Struct struct {
	NodeCommonField
	Type     int
	Ident    string
	Elems    []*Field
	Options  []*Option
	ElemsMap map[string]*Field // startToken hash => Field node
}

func NewStruct(start *Token, parent Node) *Struct {
	t := 0
	switch start.Value {
	case "struct":
		t = STRUCT
	case "union":
		t = UNION
	case "exception":
		t = EXCEPTION
	}
	return &Struct{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
		Type:     t,
		ElemsMap: map[string]*Field{},
	}
}

func (r *Struct) NodeType() string {
	switch r.Type {
	case STRUCT:
		return "Struct"
	case UNION:
		return "Union"
	case EXCEPTION:
		return "Exception"
	}
	return ""
}

func (r *Struct) NodeValue() interface{} {
	return *r
}

func (r *Struct) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Struct) patchFieldToMap(node *Field) {
	hash := GenTokenHash(node.StartToken)
	r.ElemsMap[hash] = node
}

func (r *Struct) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	identTok := p.nextIdent(false)
	r.Ident = identTok.Raw
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTCURLY {
		return p.unexpected(string(ru), "{")
	}
	p.next() // consume {
	var rightParenTok *Token
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTCURLY {
			rightParenTok = p.next()
			break
		}
		elem := NewField(r)
		if err = elem.parse(p); err != nil {
			return err
		}
		r.patchFieldToMap(elem)
		r.Elems = append(r.Elems, elem)
	}

	// parse options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		r.EndToken = rightParenTok
		return
	}
	p.next() // consume (
	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}
	return
}
