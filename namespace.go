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
	tok := p.nextNonWhitespace()
	if tok.Type != tLEFTPAREN {
		p.buf = tok
		r.EndToken = endIdent
		return
	}

	r.Options = []*Option{}
	var currOption *Option
	for {
		currOption = NewOption(r)
		err = currOption.parse(p)
		if err != nil {
			return err
		}
		r.Options = append(r.Options, currOption)

		tok = p.nextNonWhitespace()
		if tok.Type == tRIGHTPAREN {
			r.EndToken = tok
			break
		}
		// there are more options
		if tok.Type == tCOMMA {
			continue
		}
		// neither a comma nor a ), throw an error
		err = p.unexpected(tok.Value, ", or )")
		return err
	}
	return
}
