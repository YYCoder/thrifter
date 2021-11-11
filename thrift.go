package thrifter

type Thrift struct {
	NodeCommonField
	// thrift file name, if it exists
	FileName string
	// since Thrift is the root node, we need a property to access its children
	Nodes []Node
}

func NewThrift(start *Token, parent Node) *Thrift {
	return &Thrift{
		NodeCommonField: NodeCommonField{
			Parent:     parent,
			StartToken: start,
		},
	}
}

func (r *Thrift) String() string {
	return toString(r.StartToken, r.EndToken)
}

func (r *Thrift) parse(p *Parser) (err error) {
	tok := p.next()

	if r.StartToken == nil {
		r.StartToken = tok
	}

	switch {
	case tok.Type == tCOMMENT ||
		tok.Type == tSPACE ||
		tok.Type == tLINEBREAK ||
		tok.Type == tRETURN ||
		tok.Type == tTAB:
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tNAMESPACE:
		namespace := NewNamespace(tok, r)
		if err = namespace.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, namespace)
	case tok.Type == tENUM:
	case tok.Type == tCONST:
		cst := NewConst(tok, r)
		if err = cst.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, cst)
	case tok.Type == tSERVICE:
	case tok.Type == tSTRUCT:
	case tok.Type == tINCLUDE:
	case tok.Type == tTYPEDEF:
		td := NewTypeDef(tok, r)
		if err = td.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, td)
	case tok.Type == tUNION:
	case tok.Type == tEXCEPTION:
	case tok.Type == tEOF:
		r.EndToken = tok
		return
	default:
		return p.unexpected(tok.Raw, ".thrift element {namespace|enum|const|service|struct|include|typedef|union|exception}")
	}
	return
}
