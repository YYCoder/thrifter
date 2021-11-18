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

func (r *Option) NodeType() string {
	return "Option"
}

func (r *Option) NodeValue() interface{} {
	return *r
}

func (r *Option) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Option) parse(p *Parser) (err error) {
	// can't use keyword as option name
	identTok := p.nextIdent(false)
	if identTok == nil || identTok.Type != T_IDENT {
		return p.unexpected(identTok.Raw, "identifier")
	}
	// if there is no = token
	tok := p.nextNonWhitespace()
	if tok.Type != T_EQUALS {
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

	r.Name = identTok.Raw
	r.Value = tok.Raw
	r.Parent = r
	r.StartToken = identTok
	r.EndToken = tok
	// since Options are always gathered in a slice during parent node parsing, we not need to link each Option with these pointers
	r.Next = nil
	r.Prev = nil

	return
}

func parseOptions(p *Parser, parent Node) (res []*Option, rightParenTok *Token, err error) {
	res = []*Option{}
	var currOption *Option
	for {
		ru := p.peekNonWhitespace()
		if toToken(string(ru)) == T_RIGHTPAREN {
			rightParenTok = p.next()
			break
		}

		currOption = NewOption(parent)
		err = currOption.parse(p)
		if err != nil {
			return
		}
		res = append(res, currOption)

		ru = p.peekNonWhitespace()
		if toToken(string(ru)) == T_COMMA {
			p.next() // consume comma
		}
	}
	return
}
