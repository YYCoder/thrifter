package thrifter

import "text/scanner"

type NodeCommonField struct {
	Parent     *Node
	Next       *Node
	Prev       *Node
	StartToken *Token
	EndToken   *Token
}

type Position struct {
	scanner.Position
	OffsetStart int
	OffsetEnd   int
}

type Token struct {
	Type  token
	Raw   string
	Value string
	Next  *Token
	Prev  *Token
	Pos   Position
}

type Node interface {
	// recursively output current node and its children
	String() string
	// recursively parse current node and its children
	parse(p *Parser) error
}

// Nodes have children, e.g. enum/service/struct/union
type Container interface {
	elements() []*Node
	addElement(node *Node)
}
