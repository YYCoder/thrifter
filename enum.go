package thrifter

import "strconv"

type Enum struct {
	NodeCommonField
	Ident string
	Elems []*EnumElement
}

func NewEnum(start *Token, parent Node) *Enum {
	return &Enum{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Enum) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Enum) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, _, _ := p.nextIdent(false)
	r.Ident = fullLit
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTCURLY {
		return p.unexpected(string(ru), "{")
	}
	p.next() // consume {
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tRIGHTCURLY {
			r.EndToken = p.next()
			break
		}
		elem := NewEnumElement(r)
		if err = elem.parse(p); err != nil {
			return err
		}
		r.Elems = append(r.Elems, elem)
	}
	return
}

type EnumElement struct {
	NodeCommonField
	ID    int
	Ident string
}

func NewEnumElement(parent Node) *EnumElement {
	return &EnumElement{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *EnumElement) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *EnumElement) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, startTok, endTok := p.nextIdent(false)
	r.StartToken = startTok
	r.Ident = fullLit
	ru := p.peekNonWhitespace()
	// if there is no = after enum field identifier, then directly parse EndToken
	if toToken(string(ru)) != tEQUALS {
		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == tCOMMA || toToken(string(ru)) == tSEMICOLON {
			r.EndToken = p.next()
		} else {
			r.EndToken = endTok
		}
		return
	}
	p.next() // consume =
	p.peekNonWhitespace()
	tok, err, _, isInt := p.nextNumber()
	if err != nil {
		return err
	}
	if !isInt {
		return p.unexpected(tok.Value, "integer")
	}
	id, err := strconv.ParseInt(tok.Value, 10, 64)
	if err != nil {
		return err
	}
	r.ID = int(id)
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) == tCOMMA || toToken(string(ru)) == tSEMICOLON {
		r.EndToken = p.next()
	} else {
		r.EndToken = tok
	}

	return
}
