package thrifter

type Thrift struct {
	NodeCommonField
	// thrift file name, if it exists
	FileName string
	// since Thrift is the root node, we need a property to access its children
	Nodes []Node
}

func NewThrift(parent Node, FileName string) *Thrift {
	return &Thrift{
		NodeCommonField: NodeCommonField{
			Parent: parent,
		},
		FileName: FileName,
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
		node := NewNamespace(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tENUM:
		node := NewEnum(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tCONST:
		node := NewConst(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tSERVICE:
		node := NewService(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tSTRUCT, tok.Type == tEXCEPTION, tok.Type == tUNION:
		node := NewStruct(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tINCLUDE, tok.Type == tCPP_INCLUDE:
		node := NewInclude(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tTYPEDEF:
		node := NewTypeDef(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == tEOF:
		r.EndToken = tok
		return
	default:
		return p.unexpected(tok.Raw, ".thrift element {namespace|enum|const|service|struct|include|typedef|union|exception}")
	}
	return
}
