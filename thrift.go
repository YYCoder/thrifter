package thrifter

type Thrift struct {
	NodeCommonField
	// thrift file name, if it exists
	FileName string
	Nodes    []*Node
}

func (t *Thrift) parse(p *Parser) (err error) {
	tok, err := p.next()
	if err != nil {
		// TODO:
	}

	if t.StartToken == nil {
		t.StartToken = tok
	}

	switch {
	case tok.Type == tCOMMENT ||
		tok.Type == tSPACE ||
		tok.Type == tLINEBREAK ||
		tok.Type == tRETURN ||
		tok.Type == tTAB:
		err = t.parse(p)
	case tok.Type == tNAMESPACE:
	case tok.Type == tENUM:
	case tok.Type == tCONST:
	case tok.Type == tSERVICE:
	case tok.Type == tSTRUCT:
	case tok.Type == tINCLUDE:
	case tok.Type == tTYPEDEF:
	case tok.Type == tUNION:
	case tok.Type == tEXCEPTION:
	case tok.Type == tEOF:
		t.EndToken = tok
		goto done
	default:
		// err = t.parse(p)
		return p.unexpected(tok.Raw, ".thrift element {namespace|enum|const|service|struct|include|typedef|union|exception}", p)
	}
done:
	return
}
