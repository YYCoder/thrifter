package thrifter

type Namespace struct {
	NodeCommonField
	Name    string
	Value   string
	Options []*Option
}

func NewNamespace(start *Token, parent Node) *Namespace {
	return &Namespace{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Namespace) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Namespace) parse(p *Parser) (err error) {
	r.Name, _, _ = p.nextIdent(true)
	var endIdent *Token
	r.Value, _, endIdent = p.nextIdent(true)
	ru := p.peekNonWhitespace()
	if toToken(string(ru)) != tLEFTPAREN {
		r.EndToken = endIdent
		return
	}
	p.next() // consume (

	r.Options = []*Option{}
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
		r.Options = append(r.Options, currOption)

		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == tCOMMA {
			p.next() // consume comma
		}
	}
	return
}
