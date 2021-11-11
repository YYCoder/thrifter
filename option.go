package thrifter

// Represent a single option, e.g. a = "123"
type Option struct {
	NodeCommonField
	Name  string
	Value string
}

func NewOption(parent Node) *Option {
	return &Option{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
	}
}

func (r *Option) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Option) parse(p *Parser) (err error) {
	// can't use keyword as option name
	name, start, _ := p.nextIdent(false)
	if start == nil || start.Type != tIDENT {
		return p.unexpected(name, "identifier")
	}
	// find equal token
	tok := p.nextNonWhitespace()
	if tok.Type != tEQUALS {
		err = p.unexpected(tok.Value, "=")
		return
	}
	// find next string
	nextRune := p.peekNonWhitespace()
	if nextRune != singleQuoteRune && nextRune != quoteRune {
		err = p.unexpected(tok.Value, "' or \"")
		return
	}
	// if it's string
	tok, err = p.nextString()
	if err != nil {
		return err
	}

	r.Name = name
	r.Value = tok.Raw
	r.Parent = r
	r.StartToken = start
	r.EndToken = tok
	// since Options are always gathered in a slice during parent node parsing, we not need to link each Option with these pointers
	r.Next = nil
	r.Prev = nil

	return
}
