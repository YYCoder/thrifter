package thrifter

const (
	CONST_VALUE_INT = iota + 1
	CONST_VALUE_FLOAT
	CONST_VALUE_IDENT
	CONST_VALUE_LITERAL
	CONST_VALUE_MAP
	CONST_VALUE_LIST
)

type Const struct {
	NodeCommonField
	Ident string
	Type  *FieldType
	Value *ConstValue
}

func NewConst(start *Token, parent Node) *Const {
	return &Const{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Const) NodeType() string {
	return "Const"
}

func (r *Const) NodeValue() interface{} {
	return *r
}

func (r *Const) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Const) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	r.Type = NewFieldType(r)
	if err = r.Type.parse(p); err != nil {
		return
	}
	p.peekNonWhitespace()
	identTok := p.nextIdent(false)
	r.Ident = identTok.Raw
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_EQUALS {
		return p.unexpected(string(ru), "=")
	}
	p.next() // consume T_EQUALS
	r.Value = NewConstValue(r)
	if err = r.Value.parse(p); err != nil {
		return
	}
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) == T_COMMA || toToken(string(ru)) == T_SEMICOLON {
		r.EndToken = p.next() // consume comma or semicolon
	} else {
		r.EndToken = r.Value.EndToken
	}

	return
}

type ConstValue struct {
	NodeCommonField
	Type  int
	Value string
	Map   *ConstMap
	List  *ConstList
}

func NewConstValue(parent Node) *ConstValue {
	return &ConstValue{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *ConstValue) NodeType() string {
	return "ConstValue"
}

func (r *ConstValue) NodeValue() interface{} {
	return *r
}

func (r *ConstValue) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstValue) parse(p *Parser) (err error) {
	ru := p.peekNonWhitespace()
	tok := toToken(string(ru))

	// if it's minus symbol or a digit
	if tok == T_MINUS || IsDigit(ru) {
		numberTok, err, isFloat, isInt := p.nextNumber()
		if err != nil {
			return err
		}
		if isFloat {
			r.Type = CONST_VALUE_FLOAT
		} else if isInt {
			r.Type = CONST_VALUE_INT
		}
		r.StartToken = numberTok
		r.EndToken = numberTok
		r.Value = numberTok.Value
	} else if ru == singleQuoteRune || ru == quoteRune {
		strTok, err := p.nextString()
		if err != nil {
			return err
		}
		r.Type = CONST_VALUE_LITERAL
		r.StartToken = strTok
		r.EndToken = strTok
		r.Value = strTok.Raw
	} else if tok == T_LEFTSQUARE {
		leftSquareTok := p.next()
		r.Type = CONST_VALUE_LIST
		r.StartToken = leftSquareTok
		r.List = NewConstList(leftSquareTok, r)
		err = r.List.parse(p)
		if err != nil {
			return err
		}
		r.EndToken = r.List.EndToken
	} else if tok == T_LEFTCURLY {
		leftCurlyTok := p.next()
		r.Type = CONST_VALUE_MAP
		r.StartToken = leftCurlyTok
		r.Map = NewConstMap(leftCurlyTok, r)
		err = r.Map.parse(p)
		if err != nil {
			return err
		}
		r.EndToken = r.Map.EndToken
	} else {
		identTok := p.nextIdent(false)
		if identTok.Type != T_IDENT {
			return p.unexpected("identifier", identTok.Raw)
		}
		r.Type = CONST_VALUE_IDENT
		r.StartToken = identTok
		r.EndToken = identTok
		r.Value = identTok.Raw
	}

	return
}

type ConstMap struct {
	NodeCommonField
	MapKeyList   []ConstValue
	MapValueList []ConstValue
	// since directly use map structure its hard to index, we use slice to represent const map, use slice index to mapping
	// Map map[ConstValue]ConstValue
}

func NewConstMap(start *Token, parent Node) *ConstMap {
	return &ConstMap{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
		// Map: map[ConstValue]ConstValue{},
	}
}

func (r *ConstMap) NodeValue() interface{} {
	return *r
}

func (r *ConstMap) NodeType() string {
	return "ConstMap"
}

func (r *ConstMap) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstMap) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	for {
		ru := p.peekNonWhitespace() // consume white spaces
		ruTok := toToken(string(ru))
		if ruTok == T_RIGHTCURLY {
			r.EndToken = p.next() // consume right curly
			break
		}

		keyNode := NewConstValue(r)
		err = keyNode.parse(p)
		if err != nil {
			return err
		}
		ru = p.peekNonWhitespace() // consume white spaces
		if toToken(string(ru)) != T_COLON {
			return p.unexpected(":", string(ru))
		}
		p.next()              // consume T_COLON
		p.peekNonWhitespace() // consume white spaces
		valNode := NewConstValue(r)
		err = valNode.parse(p)
		if err != nil {
			return err
		}

		r.MapKeyList = append(r.MapKeyList, *keyNode)
		r.MapValueList = append(r.MapValueList, *valNode)

		ru = p.peekNonWhitespace() // consume white spaces
		ruTok = toToken(string(ru))
		if ruTok == T_COMMA || ruTok == T_SEMICOLON {
			p.next() // consume separator
		}

	}
	return
}

type ConstList struct {
	NodeCommonField
	Elems []*ConstValue
}

func NewConstList(start *Token, parent Node) *ConstList {
	return &ConstList{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *ConstList) NodeType() string {
	return "ConstList"
}

func (r *ConstList) NodeValue() interface{} {
	return *r
}

func (r *ConstList) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstList) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	for {
		ru := p.peekNonWhitespace()
		nextTok := toToken(string(ru))
		// const list end
		if nextTok == T_RIGHTSQUARE {
			r.EndToken = p.next() // consume right square
			break
		}
		valNode := NewConstValue(r)
		err = valNode.parse(p)
		if err != nil {
			return err
		}
		r.Elems = append(r.Elems, valNode)
		ru = p.peekNonWhitespace()
		nextTok = toToken(string(ru))
		if nextTok == T_COMMA || nextTok == T_SEMICOLON {
			p.next() // consume list separator
		}
	}
	return
}
