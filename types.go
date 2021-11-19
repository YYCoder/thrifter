package thrifter

import "text/scanner"

type NodeCommonField struct {
	Parent     Node
	Next       Node
	Prev       Node
	StartToken *Token
	EndToken   *Token
}

type Token struct {
	Type  token
	Raw   string // tokens raw value, e.g. comments contain prefix, like // or /* or #; strings contain ' or "
	Value string // tokens transformed value
	Next  *Token
	Prev  *Token
	Pos   scanner.Position
}

type Node interface {
	// recursively output current node and its children
	String() string
	// recursively parse current node and its children
	parse(p *Parser) error
	// get node value
	NodeValue() interface{}
	// get node type, value specified from each node
	NodeType() string
}
