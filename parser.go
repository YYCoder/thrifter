package thrifter

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"text/scanner"
)

func NewParser(rd io.Reader, debug bool) *Parser {
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
	// Scan error callback
	s.Error = func(s *scanner.Scanner, msg string) {
		fmt.Printf("Scan error: %v\n", msg)
		os.Exit(1)
	}
	res := &Parser{scanner: s, debug: debug}
	return res
}

type Parser struct {
	debug     bool
	scanner   *scanner.Scanner
	currToken *Token
	buf       *Token
}

// parse a thrift file
func (p *Parser) Parse(fileName string) (res *Thrift, err error) {
	res = NewThrift(nil, fileName)
	err = res.parse(p)
	return
}

// build token linked-list, and return consumed token
func (p *Parser) next() (res *Token) {
	// if buffer contains a token, consume buffer first
	if p.buf != nil {
		res = p.buf
		p.buf = nil
		return
	}
	s := p.scanner
	t := s.Scan()
	if t == scanner.EOF {
		res = &Token{
			Type: T_EOF,
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
	p.chainToken(res)
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
		if commentType == BASH_LIKE_COMMENT || commentType == SINGLE_LINE_COMMENT {
			if r == '\n' || r == '\r' {
				break
			} else {
				p.scanner.Next() // consume next token
			}
		} else if commentType == MULTI_LINE_COMMENT {
			if r == '*' {
				p.scanner.Next() // consume next token
				nextRune := p.peek()
				if nextRune == '/' {
					p.scanner.Next() // consume next token
					fullLit += "*/"
					break
				}
			} else {
				p.scanner.Next() // consume next token
			}
		}
		fullLit += string(r)
	}

	res = &Token{
		Value: getCommentValue(fullLit, commentType),
		Raw:   fullLit,
		Type:  T_COMMENT,
		Prev:  p.currToken,
		Pos:   p.scanner.Position,
	}
	return
}

// TODO: concat dot-separated ident into one token
// Find next identifier, it will consume white spaces during scanning.
// 1. Allow leading && trailing dot.
// 2. If keywordAllowed == true, it will allow keyword inside an identifier, e.g. enum.aaa.struct. In this case, the token for keyword will be replace to T_IDENT, since the meaning for it is no more a keyword.
// 3. For dot-separated identifier, it will automatically connected to a single string.
func (p *Parser) nextIdent(keywordAllowed bool) (res string, startToken *Token, endToken *Token) {
	var fullLit string
	var skipDot bool
	// if buffer containers a token, consume buffer first
	if p.buf != nil && p.buf.Type == T_IDENT {
		startToken, endToken = p.buf, p.buf
		fullLit = p.buf.Value
		p.buf = nil
	} else {
		t := p.nextNonWhitespace()
		tok, lit := t.Type, t.Value
		if T_IDENT != tok && T_DOT != tok {
			// can be keyword, change its token.Type
			if IsKeyword(tok) && keywordAllowed {
				t.Type = T_IDENT
			} else {
				return
			}
			// proceed with keyword as first literal
		}
		startToken, endToken = t, t
		fullLit = lit
		// if we have a leading dot, we need to skip dot handling in first iteration
		skipDot = false
		if t.Type == T_DOT {
			skipDot = true
		}
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
		if IsKeyword(tok.Type) && keywordAllowed {
			tok.Type = T_IDENT
		}
		if T_IDENT != tok.Type {
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

// Scan next Unicode character, consumes white spaces, comments and first non-whitespaces token.
func (p *Parser) nextNonWhitespace() (res *Token) {
	r := p.peek()

	// see if it's a comment
	if string(r) == "/" || string(r) == "#" {
		p.scanner.Next() // consume comment first unicode character
		isComment, commentType := p.isComment(r)
		if isComment {
			var err error
			tok, err := p.nextComment(commentType)
			if err != nil {
				fmt.Printf("Scan error: %v\n", err.Error())
				os.Exit(1)
			}
			p.chainToken(tok)
			return p.nextNonWhitespace()
		} else {
			tok := &Token{
				Type:  T_IDENT,
				Raw:   string(r),
				Value: string(r),
				Prev:  p.currToken,
				Pos:   p.scanner.Position,
			}
			p.chainToken(tok)
			return tok
		}
	} else if IsWhitespace(toToken(string(r))) {
		r = p.scanner.Next() // consume whitespaces
		tok := &Token{
			Type:  toToken(string(r)),
			Raw:   string(r),
			Value: string(r),
			Prev:  p.currToken,
			Pos:   p.scanner.Position,
		}
		p.chainToken(tok)
		return p.nextNonWhitespace()
	}
	return p.next()
}

// Scan next Unicode character, only consume white spaces, will not consume first non-whitespaces character.
func (p *Parser) peekNonWhitespace() (r rune) {
	r = p.peek()

	// see if it's a comment
	if string(r) == "/" || string(r) == "#" {
		p.scanner.Next() // consume comment first unicode character
		isComment, commentType := p.isComment(r)
		if isComment {
			var err error
			tok, err := p.nextComment(commentType)
			if err != nil {
				fmt.Printf("Scan error: %v\n", err.Error())
				os.Exit(1)
			}
			p.chainToken(tok)
			return p.peekNonWhitespace()
		} else {
			return r
		}
	} else if IsWhitespace(toToken(string(r))) {
		r = p.scanner.Next() // consume whitespaces
		tok := &Token{
			Type:  toToken(string(r)),
			Raw:   string(r),
			Value: string(r),
			Prev:  p.currToken,
			Pos:   p.scanner.Position,
		}
		p.chainToken(tok)
		return p.peekNonWhitespace()
	}
	return r
}

// Note: assume we found next token is ' or ", concat unicode character into a single string token.
// we can't use p.next() to scan token, because if string contains // or /* characters it will be parsed as comment.
func (p *Parser) nextString() (res *Token, err error) {
	r := p.scanner.Next()
	tok := toToken(string(r))
	if tok != T_SINGLEQUOTE && tok != T_QUOTE {
		err = p.unexpected(string(r), "single quote || quote")
		return
	}
	quoteType := tok
	var fullLit string
	if quoteType == T_SINGLEQUOTE {
		fullLit = singleQuoteString
	} else {
		fullLit = quoteString
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
		Type:  T_STRING,
		Raw:   fullLit,
		Value: val,
		Prev:  p.currToken,
		Pos:   p.scanner.Position,
	}
	p.chainToken(res)
	return
}

// assume we found next token is a number, if it is minus number, concat minus symbol and digits into one token
func (p *Parser) nextNumber() (res *Token, err error, isFloat bool, isInt bool) {
	var fullLit string
	r := p.peekNonWhitespace()
	if IsDigit(r) {
		p.scanner.Scan()
		fullLit = p.scanner.TokenText()
		if isFloat, isInt = IsNumber(fullLit); !isFloat && !isInt {
			err = p.unexpected("digit", fullLit)
			return
		}
		res = &Token{
			Type:  T_NUMBER,
			Raw:   fullLit,
			Value: fullLit,
			Prev:  p.currToken,
			Pos:   p.scanner.Position,
		}
		p.chainToken(res)
	} else if toToken(string(r)) == T_MINUS {
		p.scanner.Next() // consume minus symbol
		fullLit = "-"
		p.scanner.Scan()
		num := p.scanner.TokenText()
		if isFloat, isInt = IsNumber(num); !isFloat && !isInt {
			err = p.unexpected("digit", num)
			return
		} else {
			fullLit += num
		}
		res = &Token{
			Type:  T_NUMBER,
			Raw:   fullLit,
			Value: fullLit,
			Prev:  p.currToken,
			Pos:   p.scanner.Position,
		}
		p.chainToken(res)
	} else {
		err = p.unexpected("- symbol or digit", string(r))
		return
	}
	return
}

// chain token to current token's next pointer
func (p *Parser) chainToken(tok *Token) {
	if p.currToken != nil {
		p.currToken.Next = tok
	}
	p.currToken = tok
	return
}

func (p *Parser) unexpected(found, expected string) error {
	var debug string
	if p.debug {
		_, file, line, _ := runtime.Caller(1)
		debug = fmt.Sprintf(" at %s:%d", file, line)
	}
	return fmt.Errorf("%v: found %q but expected [%s], debug info %s", p.scanner.Position, found, expected, debug)
}
