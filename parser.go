package thrifter

import (
	"fmt"
	"io"
	"runtime"
	"text/scanner"
)

func NewParser(rd io.Reader) *Parser {
	s := new(scanner.Scanner)
	s.Init(rd)
	s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ' // do not filter tab/lineBreak/return/space, since we want to be non-destructive
	s.Mode = scanner.ScanRawStrings |
		scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanFloats |
		scanner.ScanChars |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments
	res := &Parser{scanner: s}
	return res
}

type Parser struct {
	debug     bool
	scanner   *scanner.Scanner
	currToken *Token
}

// parse a thrift file
func (p *Parser) Parse() (res *Thrift, err error) {
	res = new(Thrift)
	err = res.parse(p)
	return
}

func (p *Parser) next() (res *Token, err error) {
	s := p.scanner
	_ = s.Scan()
	str := s.TokenText()
	token := toToken(str)

	var curPos Position
	if p.currToken != nil {
		curPos = p.currToken.Pos
	} else {
		curPos = Position{
			p.scanner.Position,
			p.scanner.Position.Offset,
			1,
		}
	}
	res = &Token{
		Type:  token,
		Prev:  p.currToken,
		Raw:   str,
		Value: str,
		Pos: Position{
			p.scanner.Position,
			p.scanner.Position.Offset,
			curPos.OffsetStart + 1,
		},
	}
	if p.currToken != nil {
		p.currToken.Next = res
	}
	p.currToken = res
	return
}

// TODO: 待优化输出格式以及确认是否使用 runtime.Caller
func (p *Parser) unexpected(found, expected string, obj interface{}) error {
	debug := ""
	if p.debug {
		_, file, line, _ := runtime.Caller(1)
		debug = fmt.Sprintf(" at %s:%d (with %#v)", file, line, obj)
	}
	return fmt.Errorf("%v: found %q but expected [%s]%s", p.scanner.Position, found, expected, debug)
}
