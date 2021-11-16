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

func (r *Thrift) NodeType() string {
	return "Thrift"
}

func (r *Thrift) NodeValue() interface{} {
	return *r
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
	case tok.Type == T_COMMENT ||
		tok.Type == T_SPACE ||
		tok.Type == T_LINEBREAK ||
		tok.Type == T_RETURN ||
		tok.Type == T_TAB:
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_NAMESPACE:
		node := NewNamespace(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_SENUM:
	case tok.Type == T_ENUM:
		node := NewEnum(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_CONST:
		node := NewConst(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_SERVICE:
		node := NewService(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_STRUCT, tok.Type == T_EXCEPTION, tok.Type == T_UNION:
		node := NewStruct(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_INCLUDE, tok.Type == T_CPP_INCLUDE:
		node := NewInclude(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_TYPEDEF:
		node := NewTypeDef(tok, r)
		if err = node.parse(p); err != nil {
			return
		}
		r.Nodes = append(r.Nodes, node)
		if err = r.parse(p); err != nil {
			return
		}
	case tok.Type == T_EOF:
		r.EndToken = tok
		return
	default:
		return p.unexpected(tok.Raw, ".thrift element {namespace|enum|const|service|struct|include|typedef|union|exception}")
	}
	return
}
