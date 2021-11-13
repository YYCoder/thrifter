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
	ident, _, _ := p.nextIdent(false)
	r.Ident = ident
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tEQUALS {
		return p.unexpected(string(ru), "=")
	}
	p.next() // consume tEQUALS
	r.Value = NewConstValue(r)
	if err = r.Value.parse(p); err != nil {
		return
	}
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) == tCOMMA || toToken(string(ru)) == tSEMICOLON {
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
	if tok == tMINUS || isDigit(ru) {
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
	} else if tok == tLEFTSQUARE {
		leftSquareTok := p.next()
		r.Type = CONST_VALUE_LIST
		r.StartToken = leftSquareTok
		r.List = NewConstList(leftSquareTok, r)
		err = r.List.parse(p)
		if err != nil {
			return err
		}
		r.EndToken = r.List.EndToken
	} else if tok == tLEFTCURLY {
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
		fullLit, startTok, endTok := p.nextIdent(false)
		if startTok.Type != tIDENT || endTok.Type != tIDENT {
			return p.unexpected("identifier", fullLit)
		}
		r.Type = CONST_VALUE_IDENT
		r.StartToken = startTok
		r.EndToken = endTok
		r.Value = fullLit
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
		if ruTok == tRIGHTCURLY {
			r.EndToken = p.next() // consume right curly
			break
		}

		keyNode := NewConstValue(r)
		err = keyNode.parse(p)
		if err != nil {
			return err
		}
		ru = p.peekNonWhitespace() // consume white spaces
		if toToken(string(ru)) != tCOLON {
			return p.unexpected(":", string(ru))
		}
		p.next()              // consume tCOLON
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
		if ruTok == tCOMMA || ruTok == tSEMICOLON {
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
		if nextTok == tRIGHTSQUARE {
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
		if nextTok == tCOMMA || nextTok == tSEMICOLON {
			p.next() // consume list separator
		}
	}
	return
}
