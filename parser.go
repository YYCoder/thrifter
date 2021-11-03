package thrifter

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"text/scanner"
)

func NewParser(rd io.Reader) *Parser {
	s := new(scanner.Scanner)
	s.Init(rd)
	s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ' // do not filter tab/lineBreak/return/space, since we want to be non-destructive
	s.Mode = scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanFloats |
		// scanner.ScanChars |
		// scanner.ScanStrings |
		// scanner.ScanComments |
		scanner.ScanRawStrings
	res := &Parser{scanner: s}
	// Scan error callback
	s.Error = func(s *scanner.Scanner, msg string) {
		fmt.Printf("Scan error: %v\n", msg)
		os.Exit(1)
	}
	return res
}

type Parser struct {
	debug     bool
	scanner   *scanner.Scanner
	currToken *Token
	buf       *Token
}

// parse a thrift file
func (p *Parser) Parse() (res *Thrift, err error) {
	res = NewThrift(nil, nil)
	err = res.parse(p)
	return
}

// build token linked-list, and return consumed token
func (p *Parser) next() (res *Token) {
	// if buffer containers a token, consume buffer first
	if p.buf != nil {
		res = p.buf
		p.buf = nil
		return
	}
	s := p.scanner
	t := s.Scan()
	if t == scanner.EOF {
		res = &Token{
			Type: tEOF,
			Prev: p.currToken,
			Pos:  p.scanner.Position,
		}
	} else if isComment, ct := p.isComment(t); isComment {
		var err error
		res, err = p.nextComment(ct)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err.Error())
			os.Exit(1)
		}
	} else {
		str := s.TokenText()
		token := toToken(str)

		res = &Token{
			Type:  token,
			Prev:  p.currToken,
			Raw:   str,
			Value: str,
			Pos:   p.scanner.Position,
		}
	}
	if p.currToken != nil {
		p.currToken.Next = res
	}
	p.currToken = res
	return
}

// Find out if it is a comment token.
func (p *Parser) isComment(r rune) (res bool, commentType int) {
	if r == '/' {
		nextRune := p.peek()
		if nextRune == '/' {
			return true, SINGLE_LINE_COMMENT
		}
		if nextRune == '*' {
			return true, MULTI_LINE_COMMENT
		}
	} else if r == '#' {
		return true, BASH_LIKE_COMMENT
	}
	return false, 0
}

// Scan and return comment token.
func (p *Parser) nextComment(commentType int) (res *Token, err error) {
	var r rune
	var fullLit string
	switch commentType {
	case BASH_LIKE_COMMENT:
		fullLit = "#"
	case SINGLE_LINE_COMMENT:
		fullLit = "/"
	case MULTI_LINE_COMMENT:
		fullLit = "/*"
		p.scanner.Next() // consume * token
	}

	for {
		r = p.peek()
		if r == scanner.EOF {
			if commentType == MULTI_LINE_COMMENT {
				err = fmt.Errorf("unterminated block comment, at %+v", p.scanner.Position)
				return
			}
			break
		}
		p.scanner.Next() // consume next token
		if (r == '\n' || r == '\r') && (commentType == BASH_LIKE_COMMENT || commentType == SINGLE_LINE_COMMENT) {
			break
		}
		if commentType == MULTI_LINE_COMMENT && r == '*' {
			nextRune := p.peek()
			if nextRune == '/' {
				p.scanner.Next() // consume next token
				fullLit += "*/"
				break
			}
		}
		fullLit += string(r)
	}

	res = &Token{
		Value: getCommentValue(fullLit, commentType),
		Raw:   fullLit,
		Type:  tCOMMENT,
		Prev:  p.currToken,
		Pos:   p.scanner.Position,
	}
	return
}

// TODO: concat dot-separated ident into one token
// Find next identifier, it will consume white spaces during scanning.
// 1. Allow leading && trailing dot.
// 2. If keywordAllowed == true, it will allow keyword inside an identifier, e.g. enum.aaa.struct. In this case, the token for keyword will be replace to tIDENT, since
// the meaning for is no more a keyword.
// 3. For dot-separated identifier, it will automatically connected to a single string.
func (p *Parser) nextIdent(keywordAllowed bool) (res string, startToken *Token, endToken *Token) {
	t := p.nextNonWhitespace()
	tok, lit := t.Type, t.Value
	if tIDENT != tok && tDOT != tok {
		// can be keyword, change its token.Type
		if isKeyword(tok) && keywordAllowed {
			t.Type = tIDENT
		} else {
			return
		}
		// proceed with keyword as first literal
	}
	startToken, endToken = t, t
	fullLit := lit
	// if we have a leading dot, we need to skip dot handling in first iteration
	skipDot := false
	if t.Type == tDOT {
		skipDot = true
	}

	for {
		if skipDot {
			skipDot = false
			fullLit = ""
		} else {
			r := p.peek()
			if '.' != r {
				break
			}
			endToken = p.next() // consume dot
		}
		// scan next token, see if it's a identifier or keyword, if not, save it to p.buf until next p.next() calling
		tok := p.next()
		if isKeyword(tok.Type) && keywordAllowed {
			tok.Type = tIDENT
		}
		if tIDENT != tok.Type {
			fullLit = fmt.Sprintf("%s.", fullLit)
			p.buf = tok
			break
		}
		fullLit = fmt.Sprintf("%s.%s", fullLit, tok.Value)
		endToken = tok
	}
	return fullLit, startToken, endToken
}

func (p *Parser) peek() rune {
	return p.scanner.Peek()
}

// Scan next Unicode character, consumes white spaces and first non-whitespaces character.
func (p *Parser) nextNonWhitespace() (res *Token) {
	r := p.peek()
	if isWhitespace(toToken(string(r))) {
		p.next() // consume whitespaces
		return p.nextNonWhitespace()
	}
	return p.next()
}

// Scan next Unicode character, only consume white spaces, will not consume first non-whitespaces character.
func (p *Parser) peekNonWhitespace() (r rune) {
	r = p.peek()
	if isWhitespace(toToken(string(r))) {
		p.next() // consume whitespaces
		return p.peekNonWhitespace()
	}
	return r
}

// Note: assume we found next token is ' or ", concat unicode character into a single string token.
// we can't use p.next() to scan token, because if string contains // or /* characters it will be parsed as comment.
func (p *Parser) nextString() (res *Token, err error) {
	r := p.scanner.Next()
	tok := toToken(string(r))
	if tok != tSINGLEQUOTE && tok != tQUOTE {
		err = p.unexpected(string(r), "single quote || quote")
		return
	}
	quoteType := tok
	var fullLit string
	if quoteType == tSINGLEQUOTE {
		fullLit = singleQuoteString
	} else {
		fullLit = doubleQuoteString
	}

	for {
		r := p.scanner.Next()
		// invalid string
		if r == scanner.EOF || r == '\n' || r == '\r' {
			err = p.unexpected("EOF or LineBreak", "single quote || quote")
			return
		}
		fullLit += string(r)
		// find the ending quote
		if toToken(string(r)) == quoteType {
			break
		}
	}

	val, _ := unQuote(fullLit)
	res = &Token{
		Type:  tSTRING,
		Raw:   fullLit,
		Value: val,
		Prev:  p.currToken,
		Pos:   p.scanner.Position,
	}
	if p.currToken != nil {
		p.currToken.Next = res
	}
	p.currToken = res
	return
}

// determine whether it is an integer or a float number
func (p *Parser) isNumber(str string) (isFloat bool, isInt bool) {
	isFloat, _ = regexp.MatchString("^\\d+\\.\\d+$", str)
	isInt, _ = regexp.MatchString("^\\d+$", str)
	return
}

// assume we found next token is a number
func (p *Parser) nextNumber() (res *Token, err error) {
	p.next()
	return
}

func (p *Parser) unexpected(found, expected string) error {
	return fmt.Errorf("%v: found %q but expected [%s]", p.scanner.Position, found, expected)
}
