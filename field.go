package thrifter

import "strconv"

// Field represent a field within struct/union/exception
type Field struct {
	NodeCommonField
	ID           int
	Requiredness string
	FieldType    *FieldType
	Ident        string
	DefaultValue *ConstValue
	Options      []*Option
}

func NewField(parent Node) *Field {
	return &Field{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *Field) NodeType() string {
	return "Field"
}

func (r *Field) NodeValue() interface{} {
	return *r
}

func (r *Field) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Field) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	// parse ID
	idToken, err, _, isInt := p.nextNumber()
	if err != nil {
		return err
	}
	if !isInt {
		return p.unexpected(idToken.Value, "integer")
	}
	id, err := strconv.ParseInt(idToken.Value, 10, 64)
	if err != nil {
		return err
	}
	r.ID = int(id)
	r.StartToken = idToken
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tCOLON {
		return p.unexpected(string(ru), ":")
	}
	p.next() // consume :

	// parse requiredness
	p.peekNonWhitespace()
	tok := p.next()
	if tok.Value == "required" || tok.Value == "optional" {
		r.Requiredness = tok.Value
	} else {
		p.buf = tok
	}

	// parse field type
	r.FieldType = NewFieldType(r)
	if err = r.FieldType.parse(p); err != nil {
		return err
	}

	// parse identifier
	p.peekNonWhitespace()
	ident, _, endTok := p.nextIdent(false)
	r.Ident = ident

	// parse DefaultValue/Options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) == tEQUALS {
		p.next() // consume =
		cst, err := r.parseDefaultValue(p)
		if err != nil {
			return err
		}
		r.DefaultValue = cst
		// see if there are options
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tLEFTPAREN {
			p.next() // consume (
			var rightParenTok *Token
			r.Options, rightParenTok, err = r.parseOptions(p)
			if err != nil {
				return err
			}
			if err = r.parseEnd(rightParenTok, p); err != nil {
				return err
			}
		} else {
			if err = r.parseEnd(cst.EndToken, p); err != nil {
				return err
			}
		}
	} else if toToken(string(ru)) == tLEFTPAREN {
		p.next() // consume (
		var rightParenTok *Token
		r.Options, rightParenTok, err = r.parseOptions(p)
		if err != nil {
			return err
		}
		if err = r.parseEnd(rightParenTok, p); err != nil {
			return err
		}
	} else {
		if err = r.parseEnd(endTok, p); err != nil {
			return err
		}
	}

	return
}

func (r *Field) parseEnd(defaultEnd *Token, p *Parser) (err error) {
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) == tCOMMA || toToken(string(ru)) == tSEMICOLON {
		r.EndToken = p.next()
	} else {
		r.EndToken = defaultEnd
	}
	return
}

func (r *Field) parseDefaultValue(p *Parser) (res *ConstValue, err error) {
	p.peekNonWhitespace()
	res = NewConstValue(r)
	if err = res.parse(p); err != nil {
		return nil, err
	}
	return
}

func (r *Field) parseOptions(p *Parser) (options []*Option, rightParenTok *Token, err error) {
	var currOption *Option
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tRIGHTPAREN {
			rightParenTok = p.next()
			break
		}

		currOption = NewOption(r)
		err = currOption.parse(p)
		if err != nil {
			return
		}
		options = append(options, currOption)

		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == tCOMMA {
			p.next() // consume comma
		}
	}
	return
}
