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
	Name  string
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

func (r *Const) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Const) parse(p *Parser) (err error) {
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

func (r *ConstValue) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstValue) parse(p *Parser) (err error) {
	startTok := p.next()

	if startTok.Type == tMINUS {
		tok := p.next()
		isFloat, isInt := p.isNumber(tok.Value)
		if !isFloat && !isInt {
			return p.unexpected(tok.Value, "float or int after - symbol")
		}
		if isFloat {
			r.Type = CONST_VALUE_FLOAT
			r.Value = tok.Value
		} else if isInt {
			r.Type = CONST_VALUE_INT
			r.Value = tok.Value
		}
		r.EndToken = tok
	} else if isFloat, isInt := p.isNumber(startTok.Value); isFloat || isInt {
		if isFloat {
			r.Type = CONST_VALUE_FLOAT
			r.Value = startTok.Value
		} else if isInt {
			r.Type = CONST_VALUE_INT
			r.Value = startTok.Value
		}
		r.EndToken = startTok
	}

	r.StartToken = startTok

	return
}

type ConstMap struct {
	NodeCommonField
	Key   ConstValue
	Value ConstValue
}

func NewConstMap(start *Token, parent Node) *ConstMap {
	return &ConstMap{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *ConstMap) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstMap) parse(p *Parser) (err error) {
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

func (r *ConstList) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *ConstList) parse(p *Parser) (err error) {
	return
}
