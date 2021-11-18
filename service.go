package thrifter

type Service struct {
	NodeCommonField
	Ident    string
	Elems    []*Function
	Extends  string
	Options  []*Option
	ElemsMap map[string]*Function // startToken hash => Function node
}

func NewService(start *Token, parent Node) *Service {
	return &Service{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
		ElemsMap: map[string]*Function{},
	}
}

func (r *Service) NodeType() string {
	return "Service"
}

func (r *Service) NodeValue() interface{} {
	return *r
}

func (r *Service) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Service) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	identTok := p.nextIdent(false)
	r.Ident = identTok.Raw
	tok := p.nextNonWhitespace()
	if tok.Type == T_LEFTCURLY {
		r.Elems, err = r.parseFunctions(p)
		if err != nil {
			return err
		}
		// parse options
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) != T_LEFTPAREN {
			return
		}
		p.next() // consume (
		r.Options, r.EndToken, err = parseOptions(p, r)
		if err != nil {
			return err
		}
	} else if tok.Value == "extends" {
		identTok := p.nextIdent(false)
		r.Extends = identTok.Raw
		tok := p.nextNonWhitespace()
		if tok.Type == T_LEFTCURLY {
			r.Elems, err = r.parseFunctions(p)
			if err != nil {
				return err
			}
		} else {
			return p.unexpected(tok.Value, "{")
		}
	} else {
		return p.unexpected(tok.Value, "extends or {")
	}
	return
}

func (r *Service) parseFunctions(p *Parser) (funcs []*Function, err error) {
	for {
		if p.buf != nil && p.buf.Type == T_RIGHTCURLY {
			r.EndToken = p.buf
			p.buf = nil
			break
		}
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTCURLY {
			r.EndToken = p.next()
			break
		}
		elem := NewFunction(r)
		if err = elem.parse(p); err != nil {
			return nil, err
		}
		elem.patchToParentMap()
		funcs = append(funcs, elem)
	}
	return
}

const (
	FIELD_PARENT_TYPE_ARGS = iota + 1
	FIELD_PARENT_TYPE_THROWS
)

type Function struct {
	NodeCommonField
	Ident        string
	Throws       []*Field
	Oneway       bool
	FunctionType *FieldType
	Void         bool
	Args         []*Field
	Options      []*Option
	ArgsMap      map[string]*Field // startToken hash => Argument Field
	ThrowsMap    map[string]*Field // startToken hash => Throws Field
}

func NewFunction(parent Node) *Function {
	return &Function{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
		ArgsMap:   map[string]*Field{},
		ThrowsMap: map[string]*Field{},
	}
}

func (r *Function) NodeType() string {
	return "Function"
}

func (r *Function) NodeValue() interface{} {
	return *r
}

func (r *Function) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Function) patchToParentMap() {
	hash := GenTokenHash(r.StartToken)
	parent := r.Parent.(*Service)
	parent.ElemsMap[hash] = r
}

func (r *Function) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	identTok := p.nextIdent(true)
	if identTok.Raw == "oneway" {
		r.StartToken = identTok
		r.Oneway = true
		p.peekNonWhitespace()
		identTok := p.nextIdent(true)
		if identTok.Raw == "void" {
			r.Void = true
		} else {
			r.FunctionType = NewFieldType(r)
			if err = r.FunctionType.parse(p); err != nil {
				return err
			}
		}
	} else if identTok.Raw == "void" {
		r.StartToken = identTok
		r.Void = true
	} else {
		p.buf = identTok
		r.FunctionType = NewFieldType(r)
		if err = r.FunctionType.parse(p); err != nil {
			return err
		}
		r.StartToken = r.FunctionType.StartToken
	}
	p.peekNonWhitespace()
	identTok = p.nextIdent(false)
	r.Ident = identTok.Raw

	// parse argument fields
	var rightParenTok *Token
	r.Args, rightParenTok, err = r.parseFields(p, FIELD_PARENT_TYPE_ARGS)
	if err != nil {
		return err
	}

	// parse options
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) == T_LEFTPAREN {
		r.Options, rightParenTok, err = r.parseOptions(p)
		if err != nil {
			return err
		}
	}

	// parse throws
	tok := p.nextNonWhitespace()
	if tok.Type == T_THROWS {
		r.Throws, rightParenTok, err = r.parseFields(p, FIELD_PARENT_TYPE_THROWS)
		if err != nil {
			return err
		}
		tok := p.nextNonWhitespace()
		if tok.Type == T_COMMA || tok.Type == T_SEMICOLON {
			r.EndToken = tok
		} else {
			p.buf = tok
			r.EndToken = rightParenTok
		}
	} else if tok.Type == T_COMMA || tok.Type == T_SEMICOLON {
		r.EndToken = tok
	} else {
		p.buf = tok
		r.EndToken = rightParenTok
	}
	return
}

func (r *Function) parseOptions(p *Parser) (res []*Option, rightParenTok *Token, err error) {
	p.next() // consume (
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTPAREN {
			rightParenTok = p.next()
			break
		}
		elem := NewOption(r)
		if err = elem.parse(p); err != nil {
			return nil, nil, err
		}
		res = append(res, elem)

		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == T_COMMA {
			p.next() // consume comma
		}
	}
	return
}

func (r *Function) patchFieldToMap(t int, node *Field) {
	hash := GenTokenHash(node.StartToken)
	switch t {
	case FIELD_PARENT_TYPE_ARGS:
		r.ArgsMap[hash] = node
	case FIELD_PARENT_TYPE_THROWS:
		r.ThrowsMap[hash] = node
	}
}

func (r *Function) parseFields(p *Parser, t int) (fields []*Field, rightParenTok *Token, err error) {
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		return nil, nil, p.unexpected(string(ru), "(")
	}
	p.next() // consume (
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTPAREN {
			rightParenTok = p.next()
			break
		}
		elem := NewField(r)
		if err = elem.parse(p); err != nil {
			return nil, nil, err
		}
		r.patchFieldToMap(t, elem)
		fields = append(fields, elem)
	}
	return
}
