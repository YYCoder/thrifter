package thrifter

import "strconv"

type Enum struct {
	NodeCommonField
	Ident   string
	Elems   []*EnumElement
	Options []*Option
}

func NewEnum(start *Token, parent Node) *Enum {
	return &Enum{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Enum) NodeType() string {
	return "Enum"
}

func (r *Enum) NodeValue() interface{} {
	return *r
}

func (r *Enum) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Enum) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, _, _ := p.nextIdent(false)
	r.Ident = fullLit
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTCURLY {
		return p.unexpected(string(ru), "{")
	}
	p.next() // consume {
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTCURLY {
			r.EndToken = p.next()
			break
		}
		elem := NewEnumElement(r)
		if err = elem.parse(p); err != nil {
			return err
		}
		r.Elems = append(r.Elems, elem)
	}

	// parse options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		return
	}
	p.next() // consume (
	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}
	return
}

type EnumElement struct {
	NodeCommonField
	ID      int
	Ident   string
	Options []*Option
}

func NewEnumElement(parent Node) *EnumElement {
	return &EnumElement{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *EnumElement) NodeType() string {
	return "EnumElement"
}

func (r *EnumElement) NodeValue() interface{} {
	return *r
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
	if toToken(string(ru)) != T_EQUALS {
		// parse options
		ru = p.peekNonWhitespace()
		if toToken(string(ru)) != T_LEFTPAREN {
			r.EndToken = endTok
			// parse separator
			if err = r.parseSeparator(p); err != nil {
				return err
			}
			return
		}
		p.next() // consume (
		r.Options, r.EndToken, err = parseOptions(p, r)
		if err != nil {
			return err
		}

		// parse separator
		if err = r.parseSeparator(p); err != nil {
			return err
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

	// parse options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) != T_LEFTPAREN {
		r.EndToken = tok
		// parse separator
		if err = r.parseSeparator(p); err != nil {
			return err
		}
		return
	}
	p.next() // consume (
	r.Options, r.EndToken, err = parseOptions(p, r)
	if err != nil {
		return err
	}

	// parse separator
	if err = r.parseSeparator(p); err != nil {
		return err
	}

	return
}

func (r *EnumElement) parseSeparator(p *Parser) (err error) {
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) == T_COMMA || toToken(string(ru)) == T_SEMICOLON {
		r.EndToken = p.next()
	}
	return
}
