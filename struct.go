package thrifter

const (
	STRUCT = iota + 1
	UNION
	EXCEPTION
)

type Struct struct {
	NodeCommonField
	Type    int
	Ident   string
	Elems   []*Field
	Options []*Option
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
		Type: t,
	}
}

func (r *Struct) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Struct) parse(p *Parser) (err error) {
	p.peekNonWhitespace()
	fullLit, _, _ := p.nextIdent(false)
	r.Ident = fullLit
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTCURLY {
		return p.unexpected(string(ru), "{")
	}
	p.next() // consume {
	var rightParenTok *Token
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tRIGHTCURLY {
			rightParenTok = p.next()
			break
		}
		elem := NewField(r)
		if err = elem.parse(p); err != nil {
			return err
		}
		r.Elems = append(r.Elems, elem)
	}

	// parse options
	ru = p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTPAREN {
		r.EndToken = rightParenTok
		return
	}
	r.Options, err = r.parseOptions(p)
	if err != nil {
		return err
	}
	return
}

func (r *Struct) parseOptions(p *Parser) (res []*Option, err error) {
	p.next() // consume (

	res = []*Option{}
	var currOption *Option
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == tRIGHTPAREN {
			r.EndToken = p.next()
			break
		}

		currOption = NewOption(r)
		err = currOption.parse(p)
		if err != nil {
			return
		}
		res = append(res, currOption)

		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == tCOMMA {
			p.next() // consume comma
		}
	}
	return
}
